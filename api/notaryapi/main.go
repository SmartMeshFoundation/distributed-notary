package notaryapi

import (
	"fmt"
	"net/http"
	"sync"

	"encoding/json"

	"errors"

	"bytes"

	"encoding/binary"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nkbai/log"
	"golang.org/x/net/websocket"
)

// APIName :
type APIName string

/* #nosec */
const (
	APINameNotifySCTokenDeployed = "NotifySCTokenDeployed" // 该接口在公证人参与签名部署合约时,由合约部署操作发起人在合约部署成功后,将合约信息广播给所有公证人

	APINamePKNPrefix                 = "PKN-"
	APINamePKNPhase1PubKeyProof      = APINamePKNPrefix + "Phase1PubKeyProof"
	APINamePKNPhase2PaillierKeyProof = APINamePKNPrefix + "Phase2PaillierKeyProof"
	APINamePKNPhase3SecretShare      = APINamePKNPrefix + "Phase3SecretShare"
	APINamePKNPhase4PubKeyProof      = APINamePKNPrefix + "Phase4PubKeyProof"

	APINameDSMAsk             = "DSM-Ask"
	APINameDSMNotifySelection = "DSM-NotifySelection"
	APINameDSMPhase1Broadcast = "DSM-Phase1Broadcast"
	APINameDSMPhase2MessageA  = "DSM-Phase2MessageA" // response中带Phase2MessageB
	APINameDSMPhase3DeltaI    = "DSM-Phase3DeltaI"
	APINameDSMPhase5A5BProof  = "DSM-Phase5A5BProof"
	APINameDSMPhase5CProof    = "DSM-Phase5CProof"
	APINameDSMPhase6ReceiveSI = "DSM-Phase6ReceiveSI"
	APINamePBFTMessage        = "PBFT-Message"
)

/*
NotaryAPI :
提供给其他公证人节点的API
*/
type NotaryAPI struct {
	api.BaseAPI
	notaries       []*models.NotaryInfo
	selfPrivateKey *ecdsa.PrivateKey

	/*
		发送相关
	*/
	sendingWSConnMap   *sync.Map // 保存websocket连接,key为NotaryID,都为自己建立的
	sendingChanMap     *sync.Map // 保存发送队列,key为NotaryID
	waitingResponseMap *sync.Map // 保存已经发送出去并在等待连接的请求,key为requestID
	/*
		接收相关
	*/
	receivingWSConnMap *sync.Map // 保存对方找我建立的websocket连接
}

// NewNotaryAPI :
func NewNotaryAPI(host string, selfPrivateKey *ecdsa.PrivateKey, notaries []*models.NotaryInfo) *NotaryAPI {
	var notaryAPI NotaryAPI
	wsHandler := websocket.Handler(notaryAPI.wsHandlerFunc)
	router, err := rest.MakeRouter(
		rest.Get("/", func(w rest.ResponseWriter, r *rest.Request) {
			wsHandler.ServeHTTP(w.(http.ResponseWriter), r.Request)
		}),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	notaryAPI.BaseAPI = api.NewBaseAPI("NotaryAPI-WS-Server", host, router)
	notaryAPI.notaries = notaries
	notaryAPI.selfPrivateKey = selfPrivateKey
	notaryAPI.sendingWSConnMap = new(sync.Map)
	notaryAPI.sendingChanMap = new(sync.Map)
	notaryAPI.waitingResponseMap = new(sync.Map)
	notaryAPI.receivingWSConnMap = new(sync.Map)
	// 直接在启动的时候初始化好队列以及发送线程,但不初始化连接,在给某个公证人发送/接收第一条消息的时候再初始化连接
	for _, notary := range notaries {
		if common.HexToAddress(notary.AddressStr) == crypto.PubkeyToAddress(selfPrivateKey.PublicKey) {
			continue
		}
		sendingChan := make(chan api.Req, 10*params.ShareCount)
		notaryAPI.sendingChanMap.Store(notary.ID, sendingChan)
		go notaryAPI.notaryMsgSendLoop(notary.ID, sendingChan)
	}
	return &notaryAPI
}

/*
	统一处理所有公证人的连接请求
	当第一次接到某公证人请求时,会创建连接,此时需要保存下来,并启动连接对应的监听线程
*/
func (na *NotaryAPI) wsHandlerFunc(ws *websocket.Conn) {
	// 1. 解析请求
	buf, err := readBytesFromWSConn(ws)
	if err != nil {
		log.Error("readBytesFromWSConn err : %s", err.Error())
		return
	}
	req, err := na.parseNotaryRequest(buf)
	if err != nil {
		log.Error("parse notary request err : %s", err.Error())
		return
	}
	//if req.GetRequestName() == api.APINameResponse {
	//	// 如果是response,不保存连接
	//	na.dealReq(nil, req)
	//	return
	//}
	//reqWithSessionID, ok := req.(api.ReqWithSessionID)
	//if !ok {
	//	log.Error("exception req with out sessionID and notaryID")
	//	return
	//}
	//senderID := reqWithSessionID.GetSenderNotaryID()
	//
	//wsSavedInterface, loaded := na.receivingWSConnMap.LoadOrStore(senderID, ws)
	//// 2. 保存连接,并使用存储下来的连接处理请求
	//wsSaved := wsSavedInterface.(*websocket.Conn)
	na.dealReq(ws, req)
	//if loaded {
	//	// 重复创建,出现这个情况说明我在创建连接的过程中,另外一个线程也创建了连接并保存进去了,概率低,且影响不大
	//	// 如果出现了,使用已经存进去的,抛弃现有连接
	//	log.Warn("websocket connection with notary[ID=%d] repeat create,use old one and close new", senderID)
	//	return
	//}
	//// 第一次连接,此时连接已经保存到内存,保留当前线程不关闭作为该连接的消息接收线程
	//log.Info("new websocket connection with notary[ID=%d]", senderID)
	// 启动消息接收线程,这里不能返回,返回这个ws连接就被框架回收了
	na.notaryMsgReceiveLoop(ws, 0)
}

func (na *NotaryAPI) notaryMsgReceiveLoop(ws *websocket.Conn, senderID int) {
	//log.Trace("notaryMsgReceiveLoop with notary[ID=%d] start...", senderID)
	for {
		buf, err := readBytesFromWSConn(ws)
		if err != nil {
			// 这里直接删除内存中的连接并返回就行,不用close连接,交由上层回收
			//na.receivingWSConnMap.Delete(senderID)
			log.Warn("notaryMsgReceiveLoop with notary[ID=%d] end because readBytesFromWSConn err : %s, maybe reconnect", senderID, err.Error())
			return
		}
		req, err := na.parseNotaryRequest(buf)
		if err != nil {
			// 解析失败也断掉连接,交由上层回收
			//na.receivingWSConnMap.Delete(senderID)
			log.Error("parseNotaryRequest from notary[ID=%d] err : %s", senderID, err.Error())
			log.Error("request body : \n%s", string(buf))
			api.WSWriteJSON(ws, api.NewFailResponse("", api.ErrorCodeParamsWrong, err.Error()))
			return
		}
		na.dealReq(ws, req)
	}
}

func (na *NotaryAPI) dealReq(ws *websocket.Conn, req api.Req) {
	// 1. ACK消息处理
	if req.GetRequestName() == api.APINameResponse {
		resp := req.(*api.BaseResponse)
		oldReqInterface, ok := na.waitingResponseMap.Load(resp.GetRequestID())
		if !ok {
			//log.Warn("get response of req[RequestID=%s], but can not found req in dealingWSReqMap,\n%s\nignore", resp.GetRequestID(), utils.ToJSONStringFormat(resp))
			return
		}
		oldReq := oldReqInterface.(api.ReqWithResponse)
		oldReq.WriteResponse(resp)
		na.waitingResponseMap.Delete(resp.GetRequestID())
		return
	}
	// 2. 普通消息处理
	na.SendToService(req)
	// 3. 如果是需要等待service层返回的请求,另起一个线程等待返回,并回写至ws
	if reqWithResponse, needResponse := req.(api.ReqWithResponse); needResponse {
		go func(ws *websocket.Conn, reqWithResponse api.ReqWithResponse) {
			resp := na.WaitServiceResponse(reqWithResponse)
			if reqWithSessionID, ok := reqWithResponse.(api.ReqWithSessionID); ok {
				// 公证人请求是单向连接,返回结果也使用发送线程来发送
				na.SendWSReqToNotary(resp, reqWithSessionID.GetSenderNotaryID())
			} else {
				api.WSWriteJSON(ws, resp)
			}
		}(ws, reqWithResponse)
		return
	}
}

func (na *NotaryAPI) parseNotaryRequest(content []byte) (req api.Req, err error) {
	// 1. 基础解析,获取api name
	baseReq := &api.BaseReq{}
	err = json.Unmarshal(content, &baseReq)
	if err != nil {
		return
	}
	switch baseReq.GetRequestName() {
	/*
		admin
	*/
	case APINameNotifySCTokenDeployed:
		req = &NotifySCTokenDeployedRequest{}
		err = json.Unmarshal(content, &req)
	case api.APINameResponse:
		req = &api.BaseResponse{}
		err = json.Unmarshal(content, &req)
	/*
		pkn
	*/
	case APINamePKNPhase1PubKeyProof:
		req = &KeyGenerationPhase1MessageRequest{}
		err = json.Unmarshal(content, &req)
	case APINamePKNPhase2PaillierKeyProof:
		req = &KeyGenerationPhase2MessageRequest{}
		err = json.Unmarshal(content, &req)
	case APINamePKNPhase3SecretShare:
		req = &KeyGenerationPhase3MessageRequest{}
		err = json.Unmarshal(content, &req)
	case APINamePKNPhase4PubKeyProof:
		req = &KeyGenerationPhase4MessageRequest{}
		err = json.Unmarshal(content, &req)
	/*
		dsm
	*/
	case APINameDSMAsk:
		req = &DSMAskRequest{}
		err = json.Unmarshal(content, &req)
	case APINameDSMNotifySelection:
		req = &DSMNotifySelectionRequest{}
		err = json.Unmarshal(content, &req)
	case APINameDSMPhase1Broadcast:
		req = &DSMPhase1BroadcastRequest{}
		err = json.Unmarshal(content, &req)
	case APINameDSMPhase2MessageA:
		req = &DSMPhase2MessageARequest{}
		err = json.Unmarshal(content, &req)
	case APINameDSMPhase3DeltaI:
		req = &DSMPhase3DeltaIRequest{}
		err = json.Unmarshal(content, &req)
	case APINameDSMPhase5A5BProof:
		req = &DSMPhase5A5BProofRequest{}
		err = json.Unmarshal(content, &req)
	case APINameDSMPhase5CProof:
		req = &DSMPhase5CProofRequest{}
		err = json.Unmarshal(content, &req)
	case APINameDSMPhase6ReceiveSI:
		req = &DSMPhase6ReceiveSIRequest{}
		err = json.Unmarshal(content, &req)
	case APINamePBFTMessage:
		req = &PBFTMessage{}
		err = json.Unmarshal(content, &req)
	}
	if err != nil {
		return
	}
	// 在接收到需要返回的请求时,第一时间初始化responseChan
	if r2, ok := req.(api.ReqWithResponse); ok {
		r2.NewResponseChan()
	}
	r1, ok1 := req.(api.ReqWithSignature)
	r2, ok2 := req.(api.ReqWithSessionID)
	if ok1 && ok2 {
		sender, err2 := na.getNotaryInfoByID(r2.GetSenderNotaryID())
		if err2 != nil {
			err = err2
			return
		}
		if sender.GetAddress() != r1.GetSignerETHAddress() {
			err = fmt.Errorf("signer address not euqal with msg.sender : signer=%s sender=%s", sender.AddressStr, r1.GetSignerETHAddress().String())
			return
		}
	}
	return
}

func readBytesFromWSConn(ws *websocket.Conn) (buf []byte, err error) {
	if ws == nil {
		panic("readBytesFromWSConn ws can not be nil")
	}
	var n int32
	lengthBytes := make([]byte, 4)
	for n < 4 {
		var n1 int
		if n1, err = ws.Read(lengthBytes[n:]); err != nil {
			return
		}
		n += int32(n1)
	}
	var length int32
	err = binary.Read(bytes.NewBuffer(lengthBytes), binary.BigEndian, &length)
	if err != nil {
		return
	}
	// TODO 暂定最大消息长度10k
	if length > 1024*10 {
		err = errors.New("msg too long")
		return
	}
	buf2 := make([]byte, length)
	n = 0
	for n < length {
		var n1 int
		if n1, err = ws.Read(buf2[n:]); err != nil {
			return
		}
		n += int32(n1)
	}
	if int32(n) != length {
		log.Error("read data from web socket err : data len=%d, read len=%d", length, n)
	}
	buf = buf2[:n]
	return
}
