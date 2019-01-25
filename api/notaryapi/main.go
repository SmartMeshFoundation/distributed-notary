package notaryapi

import (
	"fmt"
	"net/http"
	"sync"

	"encoding/json"

	"errors"

	"bytes"

	"encoding/binary"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ant0ine/go-json-rest/rest"
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
)

/*
NotaryAPI :
提供给其他公证人节点的API
*/
type NotaryAPI struct {
	api.BaseAPI
	notaries             []models.NotaryInfo
	wsReqSendingQueueMap *sync.Map // 存储消息发送队列,每个与我通信过的公证人节点都会在里面有一对ID-chan的key-value,且对应存在一个常驻发送线程
	dealingWSReqMap      *sync.Map // 暂存我发出去的,正在等待返回的请求
	notaryWSConnMap      *sync.Map
	notaryWSConnMapLock  sync.Mutex
}

// NewNotaryAPI :
func NewNotaryAPI(host string, notaries []models.NotaryInfo) *NotaryAPI {
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
	notaryAPI.wsReqSendingQueueMap = new(sync.Map)
	notaryAPI.dealingWSReqMap = new(sync.Map)
	notaryAPI.notaryWSConnMap = new(sync.Map)
	notaryAPI.notaries = notaries
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
	reqWithSessionID, ok := req.(api.ReqWithSessionID)
	if !ok {
		log.Error("exception req with out sessionID and notaryID")
		return
	}
	senderID := reqWithSessionID.GetSenderNotaryID()
	/*
		2. 保存连接,如过这里已经存在连接,有两种可能:
				a. 内存中保存的是之前对方找我创建的连接,如果对方找我再次创建,说明之前的连接已经断了,
				   而如果老的连接断掉的话,我这边会删除掉,所以除非恶意攻击,不可能出现这种情况,如果出现了,直接使用老的连接即可
				b. 内存中保存的是我给对方发的时候创建的连接,直接使用
	*/
	na.notaryWSConnMapLock.Lock()
	if old, ok := na.notaryWSConnMap.Load(senderID); ok {
		oldWS := old.(*websocket.Conn)
		// 4. 使用老的连接处理消息,然后直接返回,关闭这次多余的连接
		na.dealReq(oldWS, req)
		na.notaryWSConnMapLock.Unlock()
		log.Warn("got new websocket connection with notary[ID=%d], but already have one,use old and ignore new", senderID)
		return
	}
	na.notaryWSConnMap.Store(senderID, ws)
	// 4. 处理消息
	na.dealReq(ws, req)
	na.notaryWSConnMapLock.Unlock()
	// 5. 启动消息接收线程,这里不能返回,返回这个ws连接就被框架回收了
	na.notaryMsgReceiveLoop(ws, senderID)
}

func (na *NotaryAPI) notaryMsgReceiveLoop(ws *websocket.Conn, senderID int) {
	log.Trace("notaryMsgReceiveLoop with notary[ID=%d] start...", senderID)
	for {
		buf, err := readBytesFromWSConn(ws)
		if err != nil {
			// 这里直接删除内存中的连接并返回就行,不用close连接,交由上层回收
			na.notaryWSConnMap.Delete(senderID)
			log.Warn("notaryMsgReceiveLoop with notary[ID=%d] end because readBytesFromWSConn err : %s, maybe reconnect", senderID, err.Error())
			return
		}
		req, err := na.parseNotaryRequest(buf)
		if err != nil {
			// 解析失败也断掉连接,交由上层回收
			na.notaryWSConnMap.Delete(senderID)
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
		oldReqInterface, ok := na.dealingWSReqMap.Load(resp.GetRequestID())
		if !ok {
			//log.Warn("get response of req[RequestID=%s], but can not found req in dealingWSReqMap,\n%s\nignore", resp.GetRequestID(), utils.ToJSONStringFormat(resp))
			return
		}
		oldReq := oldReqInterface.(api.ReqWithResponse)
		oldReq.WriteResponse(resp)
		na.dealingWSReqMap.Delete(resp.GetRequestID())
		return
	}
	// 2. 普通消息处理
	na.SendToService(req)
	// 如果是需要返回的请求,另起一个线程等待返回,并回写至ws
	if reqWithResponse, needResponse := req.(api.ReqWithResponse); needResponse {
		go func(ws *websocket.Conn, reqWithResponse api.ReqWithResponse) {
			resp := na.WaitServiceResponse(reqWithResponse)
			api.WSWriteJSON(ws, resp)
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
		// TODO 新增dsmPhase2Response解析
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
	}
	if err != nil {
		return
	}
	//if reqWithSignature, ok := req.(api.ReqWithSignature); ok {
	//	if !reqWithSignature.VerifySign(reqWithSignature) {
	//		err = errors.New(api.ErrorCode2MsgMap[api.ErrorCodePermissionDenied])
	//	}
	//}
	return
}

func readBytesFromWSConn(ws *websocket.Conn) (buf []byte, err error) {
	if ws == nil {
		panic("readBytesFromWSConn ws can not be nil")
	}
	lengthBytes := make([]byte, 4)
	if _, err = ws.Read(lengthBytes); err != nil {
		return
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
	var n int32
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
