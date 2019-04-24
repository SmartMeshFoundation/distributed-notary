package service

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	ethevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/service/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nkbai/log"
)

// NotaryService :
type NotaryService struct {
	notaryClient       notaryapi.NotaryClient
	privateKey         *ecdsa.PrivateKey
	self               *models.NotaryInfo
	otherNotaries      []*models.NotaryInfo //这里保存除我以外的notary信息
	db                 *models.DB
	dispatchService    dispatchServiceBackend
	sessionLockMap     map[common.Hash]*sync.Mutex
	sessionLockMapLock sync.Mutex

	pknHandlerMap *sync.Map
	dsmHandlerMap *sync.Map
}

// NewNotaryService :
func NewNotaryService(db *models.DB, privateKey *ecdsa.PrivateKey, allNotaries []*models.NotaryInfo, notaryClient notaryapi.NotaryClient, dispatchService dispatchServiceBackend) (ns *NotaryService, err error) {
	ns = &NotaryService{
		db:              db,
		privateKey:      privateKey,
		sessionLockMap:  make(map[common.Hash]*sync.Mutex),
		dispatchService: dispatchService,
		notaryClient:    notaryClient,
		pknHandlerMap:   new(sync.Map),
		dsmHandlerMap:   new(sync.Map),
	}
	// 初始化self, otherNotaries
	selfAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	for _, n := range allNotaries {
		if n.GetAddress() == selfAddress {
			ns.self = n
		} else {
			ns.otherNotaries = append(ns.otherNotaries, n)
		}
	}
	models.SortNotaryInfoSlice(ns.otherNotaries)
	return
}

// OnEvent 链上事件逻辑处理 预留
func (ns *NotaryService) OnEvent(e chain.Event) {
}

// OnRequest restful请求处理
func (ns *NotaryService) OnRequest(req api.Req) {
	switch r := req.(type) {
	case *notaryapi.KeyGenerationPhase1MessageRequest:
		ns.onPKNRequest(r)
	case *notaryapi.KeyGenerationPhase2MessageRequest:
		ns.onPKNRequest(r)
	case *notaryapi.KeyGenerationPhase3MessageRequest:
		ns.onPKNRequest(r)
	case *notaryapi.KeyGenerationPhase4MessageRequest:
		ns.onPKNRequest(r)
	case *notaryapi.DSMAskRequest:
		ns.onDSMAskRequest(r)
	case *notaryapi.DSMNotifySelectionRequest:
		ns.onDSMNotifySelectionRequest(r)
	case *notaryapi.DSMPhase1BroadcastRequest:
		ns.onDSMRequest(r)
	case *notaryapi.DSMPhase2MessageARequest:
		ns.onDSMRequest(r)
	case *notaryapi.DSMPhase3DeltaIRequest:
		ns.onDSMRequest(r)
	case *notaryapi.DSMPhase5A5BProofRequest:
		ns.onDSMRequest(r)
	case *notaryapi.DSMPhase5CProofRequest:
		ns.onDSMRequest(r)
	case *notaryapi.DSMPhase6ReceiveSIRequest:
		ns.onDSMRequest(r)
	default:
		r2, ok := req.(api.ReqWithResponse)
		if ok {
			r2.WriteErrorResponse(api.ErrorCodeParamsWrong)
			return
		}
	}
	return
}

/*
主动开始一次私钥协商
*/
func (ns *NotaryService) startNewPrivateKeyNegotiation() (privateKeyInfo *models.PrivateKeyInfo, err error) {
	sessionID := utils.NewRandomHash() // 初始化一个会话ID
	var otherNotaryIDs []int
	for _, notary := range ns.otherNotaries {
		otherNotaryIDs = append(otherNotaryIDs, notary.ID)
	}
	/*
		创建一个新的PKNHandler,并保存到内存
	*/
	ph := mecdsa.NewPKNHandler(ns.db, ns.self, otherNotaryIDs, sessionID, ns.notaryClient)
	ns.pknHandlerMap.Store(sessionID, ph)
	// 开始pkn并阻塞等待
	return ph.StartPKNAndWaitFinish(nil)
}

func (ns *NotaryService) onPKNRequest(req api.Req) {
	reqWithResponse, needResponse := req.(api.ReqWithResponse)
	// 0. 获取sessionID
	reqWithSessionID, ok := req.(api.ReqWithSessionID)
	if !ok {
		log.Error("unknown msg for PKNHandler :\n%s", utils.ToJSONStringFormat(req))
		if needResponse {
			reqWithResponse.WriteErrorResponse(api.ErrorCodePermissionDenied, "unknown msg for PKNHandler ")
		}
		return
	}
	sessionID := reqWithSessionID.GetSessionID()
	// 1. 获取pkn
	var phSavedInterface interface{}
	var loaded bool
	var phase1Req *notaryapi.KeyGenerationPhase1MessageRequest
	phSavedInterface, loaded = ns.pknHandlerMap.Load(sessionID)
	if !loaded {
		var ok bool
		phase1Req, ok = req.(*notaryapi.KeyGenerationPhase1MessageRequest)
		/*
			不是phase1请求且找不到PKNHandler,拒绝
		*/
		if !ok {
			if needResponse {
				errMsg := SessionLogMsg(sessionID, "can not find PKNHandler with this sessionID")
				reqWithResponse.WriteErrorResponse(api.ErrorCodeException, errMsg)
			}
			return
		}
		/*
			被动开始,新建pkn
		*/
		var otherNotaryIDs []int
		for _, notary := range ns.otherNotaries {
			otherNotaryIDs = append(otherNotaryIDs, notary.ID)
		}
		phSavedInterface, loaded = ns.pknHandlerMap.LoadOrStore(sessionID, mecdsa.NewPKNHandler(ns.db, ns.self, otherNotaryIDs, sessionID, ns.notaryClient))
	}
	/*
		投递到ph,如果是被动开始的,启动协商线程
	*/
	ph := phSavedInterface.(*mecdsa.PKNHandler)
	if loaded {
		// 已经存在,直接投递
		ph.OnRequest(req)
		//ph.receivePhase1PubKeyProof(req.Msg, req.GetSenderNotaryID())
	} else {
		// 不存在,启动协商线程
		go ph.StartPKNAndWaitFinish(phase1Req)
	}
	if needResponse {
		reqWithResponse.WriteSuccessResponse(nil)
	}
	return
}

/*
主动开始一次签名,并等待最终的签名结果
*/
func (ns *NotaryService) startDistributedSignAndWait(msgToSign messagetosign.MessageToSign, privateKeyInfo *models.PrivateKeyInfo) (signature []byte, sessionID common.Hash, err error) {
	// 1. DSMAsk
	notaryNumNeedExpectSelf := params.ThresholdCount
	if msgToSign.GetName() == messagetosign.SpectrumContractDeployTXDataName {
		// 如果需要签名的是部署合约的tx,则要求所有公证人参与
		notaryNumNeedExpectSelf = params.ShareCount - 1
	}
	var otherDSMNotaryIDs []int
	sessionID, otherDSMNotaryIDs, err = ns.startDSMAsk(notaryNumNeedExpectSelf, privateKeyInfo.Key, msgToSign)
	if err != nil {
		log.Error(SessionLogMsg(sessionID, "startDSMAsk err = %s", err.Error()))
		return
	}
	dsmNotaryIDs := append(otherDSMNotaryIDs, ns.self.ID)
	// 2.DSMNotifySelection 通知所有公证人该次签名参与的名单,并等待参与者的返回,否则会出现其他节点尚未处理完NotifySelection请求,就收到发起者发出来的Phase1Msg,造成失败
	log.Trace(SessionLogMsg(sessionID, "DSMNotifySelection start..."))
	wg := sync.WaitGroup{}
	wg.Add(len(otherDSMNotaryIDs))
	for _, notary := range ns.otherNotaries {
		go func(notaryID int) {
			notifySelectionReq := notaryapi.NewDSMNotifySelectionRequest(sessionID, ns.self, dsmNotaryIDs)
			ns.notaryClient.SendWSReqToNotary(notifySelectionReq, notaryID)
			_, err2 := ns.notaryClient.WaitWSResponse(notifySelectionReq.GetRequestID())
			if err2 != nil {
				err = err2
			}
			// 只等待参与者的返回
			for _, id := range otherDSMNotaryIDs {
				if notaryID == id {
					wg.Done()
				}
			}
		}(notary.ID)
	}
	wg.Wait()
	if err != nil {
		log.Error(SessionLogMsg(sessionID, "DSMNotifySelection err = %s", err.Error()))
		return
	}
	log.Trace(SessionLogMsg(sessionID, "DSMNotifySelection done..."))
	// 3. 构造DSMHandler
	dh := mecdsa.NewDSMHandler(ns.db, ns.self, msgToSign, sessionID, privateKeyInfo, ns.notaryClient).RegisterOtherNotaryIDs(otherDSMNotaryIDs)
	// 4. 保存DSMHandler
	ns.dsmHandlerMap.Store(sessionID, dh)
	// 5. 开始签名过程并阻塞等待,协商发起者需要把标志位Running置为true,避免启动2个主线程
	dh.Running = true
	signature, err = dh.StartDSMAndWaitFinish()
	return
}

func (ns *NotaryService) startDSMAsk(notaryNumNeedExpectSelf int, privateKeyID common.Hash, msgToSign messagetosign.MessageToSign) (sessionID common.Hash, otherNotaryIDs []int, err error) {
	sessionID = utils.NewRandomHash()
	log.Trace(SessionLogMsg(sessionID, "DSMAsk start..."))
	m := new(sync.Map)
	wg := sync.WaitGroup{}
	wg.Add(len(ns.otherNotaries))
	for _, notary := range ns.otherNotaries {
		go func(notary *models.NotaryInfo) {
			req := notaryapi.NewDSMAskRequest(sessionID, ns.self, privateKeyID, msgToSign)
			ns.notaryClient.SendWSReqToNotary(req, notary.ID)
			_, err2 := ns.notaryClient.WaitWSResponse(req.GetRequestID())
			if err2 == nil {
				m.Store(notary.ID, true)
			} else {
				log.Warn(SessionLogMsg(sessionID, "notary[%d] refuse DSMAsk : %s", notary.ID, err2.Error()))
			}
			wg.Done()
		}(notary)
	}
	wg.Wait()
	for _, notary := range ns.otherNotaries {
		if _, ok := m.Load(notary.ID); ok {
			otherNotaryIDs = append(otherNotaryIDs, notary.ID)
			if len(otherNotaryIDs) >= notaryNumNeedExpectSelf {
				break
			}
		}
	}
	if len(otherNotaryIDs) < notaryNumNeedExpectSelf {
		err = fmt.Errorf("no enough notary to sign message, need %d but only got %d", notaryNumNeedExpectSelf+1, len(otherNotaryIDs)+1)
		return
	}
	log.Trace(SessionLogMsg(sessionID, "DSMAsk done..."))
	return
}

/*
收到 DSMAskRequest
*/
func (ns *NotaryService) onDSMAskRequest(req *notaryapi.DSMAskRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	_, ok := ns.getNotaryInfoByAddress(req.GetSignerETHAddress())
	if !ok {
		log.Warn(SessionLogMsg(req.GetSessionID(), "unknown notary %s, maybe attack", req.GetSignerETHAddress().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 校验私钥ID可用性
	var privateKeyInfo *models.PrivateKeyInfo
	privateKeyInfo, err := ns.db.LoadPrivateKeyInfo(req.PrivateKeyID)
	if err != nil {
		log.Error(SessionLogMsg(sessionID, "unknown PrivateKeyID %s", req.PrivateKeyID.String()))
		req.WriteErrorResponse(api.ErrorCodeException, "unknown PrivateKeyID")
		return
	}
	// 3. 解析msgToSign
	var msgToSign messagetosign.MessageToSign
	msgToSign, err = parseMessageToSign(req.MsgName, req.MsgToSignTransportBytes)
	if err != nil {
		log.Error(SessionLogMsg(sessionID, "parseMessageToSign err : %s", err.Error()))
		req.WriteErrorResponse(api.ErrorCodeException, "parseMessageToSign error")
		return
	}
	// 4. 校验msgToSign
	err = ns.checkMsgToSign(sessionID, privateKeyInfo, msgToSign, req.GetSenderNotaryID())
	if err != nil {
		log.Error(SessionLogMsg(sessionID, "checkMessageToSign err : %s", err.Error()))
		req.WriteErrorResponse(api.ErrorCodeException, "checkMessageToSign error")
		return
	}
	// 5. 构造DSMHandler
	dh := mecdsa.NewDSMHandler(ns.db, ns.self, msgToSign, sessionID, privateKeyInfo, ns.notaryClient)
	// 6. 保存dh,这里不会重复
	ns.dsmHandlerMap.Store(sessionID, dh)
	log.Trace(SessionLogMsg(req.GetSessionID(), "accept DSMAsk from %s", utils.APex(req.GetSignerETHAddress())))
	req.WriteSuccessResponse(nil)
	return
}

/*
收到 DSMNotifySelectionRequest
*/
func (ns *NotaryService) onDSMNotifySelectionRequest(req *notaryapi.DSMNotifySelectionRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	_, ok := ns.getNotaryInfoByAddress(req.GetSignerETHAddress())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSignerETHAddress().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 判断自己是否需要参与
	amIIn := false
	var otherNotaryIDs []int
	for _, notaryID := range req.NotaryIDs {
		if notaryID == ns.self.ID {
			amIIn = true
		} else {
			otherNotaryIDs = append(otherNotaryIDs, notaryID)
		}
	}
	// 3. 如果无需参与,即没被选中,删除dsmHandler并返回
	if !amIIn {
		ns.dsmHandlerMap.Delete(sessionID)
		req.WriteSuccessResponse(nil)
		return
	}
	// 4. 如果需要参与,获取DSMHandler
	var dsmHandlerInterface interface{}
	var loaded bool
	// 这里找不到直接报错
	dsmHandlerInterface, loaded = ns.dsmHandlerMap.Load(sessionID)
	if !loaded {
		errMsg := SessionLogMsg(sessionID, "can not find DSMHandler with this sessionID")
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	dh := dsmHandlerInterface.(*mecdsa.DSMHandler)
	// 注册参与者到dsmHandler中
	dh.RegisterOtherNotaryIDs(otherNotaryIDs)
	log.Trace(SessionLogMsg(sessionID, "DSMNotifySelection done..."))
	// 5. 返回
	req.WriteSuccessResponse(nil)
	return
}

func (ns *NotaryService) onDSMRequest(req api.Req) {
	reqWithResponse, needResponse := req.(api.ReqWithResponse)
	// 0. 获取sessionID
	reqWithSessionID, ok := req.(api.ReqWithSessionID)
	if !ok {
		log.Error("unknown msg for DSNHandler :\n%s", utils.ToJSONStringFormat(req))
		if needResponse {
			reqWithResponse.WriteErrorResponse(api.ErrorCodePermissionDenied, "unknown msg for DSNHandler ")
		}
		return
	}
	sessionID := reqWithSessionID.GetSessionID()
	// 1. 获取DSMHandler
	var dsmHandlerInterface interface{}
	var loaded bool
	// 这里找不到直接报错
	dsmHandlerInterface, loaded = ns.dsmHandlerMap.Load(sessionID)
	if !loaded {
		if needResponse {
			errMsg := SessionLogMsg(sessionID, "can not find DSMHandler with this sessionID")
			reqWithResponse.WriteErrorResponse(api.ErrorCodeException, errMsg)
		}
		return
	}
	// 由于发起者会在所有人NotifySelection处理完成后才开始发送消息,所以被动参与者可以考虑在收到发起者的phase1Msg消息之后才启动dsmHandler的主线程发送自己的phase1Msg消息
	// 如果这样,那么被动参与方收到的第一条消息一定是发起者发来的phase1Msg,并且收到的时候,dsmHandler一定已经初始化了,所以直接启动主线程
	dh := dsmHandlerInterface.(*mecdsa.DSMHandler)
	if !dh.Running {
		dh.Running = true
		go dh.StartDSMAndWaitFinish()
	}
	// 2. 投递消息
	dh.OnRequest(req)
	return
}

func (ns *NotaryService) getNotaryInfoByAddress(addr common.Address) (notaryInfo *models.NotaryInfo, ok bool) {
	for _, v := range ns.otherNotaries {
		if v.GetAddress() == addr {
			notaryInfo = v
			ok = true
			return
		}
	}
	ok = false
	return
}

/*
解析其他公证人传来的待签名消息
*/
func parseMessageToSign(msgName string, buf []byte) (msg messagetosign.MessageToSign, err error) {
	switch msgName {
	case messagetosign.SpectrumContractDeployTXDataName:
		msg = new(messagetosign.SpectrumContractDeployTXData)
		err = msg.Parse(buf)
	case messagetosign.SpectrumPrepareLockinTxDataName:
		msg = new(messagetosign.SpectrumPrepareLockinTxData)
		err = msg.Parse(buf)
	case messagetosign.EthereumPrepareLockoutTxDataName:
		msg = new(messagetosign.EthereumPrepareLockoutTxData)
		err = msg.Parse(buf)
	case messagetosign.EthereumCancelNonceTxDataName:
		msg = new(messagetosign.EthereumCancelNonceTxData)
		err = msg.Parse(buf)
	case messagetosign.BitcoinLockinTXDataName:
		msg = new(messagetosign.BitcoinLockinTXData)
		err = msg.Parse(buf)
	default:
		err = fmt.Errorf("got msg to sign which does't support, maybe attack")
	}
	return
}

/*
签名信息校验,根据收到的消息类型,自己生成一份对应的消息体,并与收到的比对
*/
func (ns *NotaryService) checkMsgToSign(sessionID common.Hash, privateKeyInfo *models.PrivateKeyInfo, msg messagetosign.MessageToSign, senderID int) (err error) {
	switch m := msg.(type) {
	// 1. 合约部署消息
	case *messagetosign.SpectrumContractDeployTXData:
		log.Trace(SessionLogMsg(sessionID, "Got %s-%s MsgToSign,run checkMsgToSign...", m.GetName(), m.DeployChainName))
		var c chain.Chain
		c, err = ns.dispatchService.getChainByName(m.DeployChainName)
		if err != nil {
			return
		}
		err = m.VerifySignBytes(c, privateKeyInfo.ToAddress())
	// 2. 侧链PrepareLockin合约调用消息
	case *messagetosign.SpectrumPrepareLockinTxData:
		log.Trace(SessionLogMsg(sessionID, "Got %s MsgToSign,run checkMsgToSign...", m.GetName()))
		// 1. 获取本地lockinInfo
		var localLockinInfo *models.LockinInfo
		localLockinInfo, err = ns.dispatchService.getLockInInfoBySCPrepareLockInRequest(m.UserRequest)
		if err != nil {
			return
		}
		// 2. 获取本地scTokenProxy
		var c chain.Chain
		c, err = ns.dispatchService.getChainByName(smcevents.ChainName)
		scTokenProxy := c.GetContractProxy(localLockinInfo.SCTokenAddress)
		// 2. 校验
		err = m.VerifySignData(scTokenProxy, privateKeyInfo, localLockinInfo, ns.dispatchService.getBtcNetworkParam())
		if err != nil {
			return
		}
		// 2.5. 更新本地locinInfo的NotaryIDInChargeID,记录该lockinInfo的负责人
		err = ns.dispatchService.updateLockinInfoNotaryIDInChargeID(localLockinInfo.SCTokenAddress, localLockinInfo.SecretHash, senderID)
	// 3. 主链PrepareLockout合约调用消息
	case *messagetosign.EthereumPrepareLockoutTxData:
		log.Trace(SessionLogMsg(sessionID, "Got %s MsgToSign,run checkMsgToSign...", m.GetName()))
		// 1. 获取本地lockoutInfo
		var localLockoutInfo *models.LockoutInfo
		localLockoutInfo, err = ns.dispatchService.getLockoutInfo(m.UserRequest.SCTokenAddress, m.UserRequest.SecretHash)
		if err != nil {
			return
		}
		// 2. 获取scTokenInfo
		scToken := ns.dispatchService.getSCTokenMetaInfoBySCTokenAddress(localLockoutInfo.SCTokenAddress)
		// 3. 获取本地mcProxy
		var c chain.Chain
		c, err = ns.dispatchService.getChainByName(ethevents.ChainName)
		mcProxy := c.GetContractProxy(scToken.MCLockedContractAddress)
		// 4. 校验
		err = m.VerifySignData(mcProxy, privateKeyInfo, localLockoutInfo)
		if err != nil {
			return
		}
		// 5. 更新本地lockoutInfo的NotaryIDInChargeID,记录该lockinInfo的负责人
		err = ns.dispatchService.updateLockoutInfoNotaryIDInChargeID(localLockoutInfo.SCTokenAddress, localLockoutInfo.SecretHash, senderID)
	// 4. nonce销毁消息
	case *messagetosign.EthereumCancelNonceTxData:
		log.Trace(SessionLogMsg(sessionID, "Got %s MsgToSign,run checkMsgToSign...", m.GetName()))
		// 1. 获取chain
		var c chain.Chain
		c, err = ns.dispatchService.getChainByName(m.ChainName)
		if err != nil {
			return
		}
		// 2. 校验account
		account := common.HexToAddress(m.Account)
		_, err = ns.db.LoadPrivateKeyInfoByAccountAddress(account)
		if err != nil {
			return
		}
		// 3. 校验
		err = m.VerifySignData(c, account)
		if err != nil {
			return
		}
		// 5. BitcointLockin消息
	case *messagetosign.BitcoinLockinTXData:
		log.Trace(SessionLogMsg(sessionID, "Got %s MsgToSign,run checkMsgToSign...", m.GetName()))
		// 0. 获取bs
		var c chain.Chain
		c, err = ns.dispatchService.getChainByName(bitcoin.ChainName)
		if err != nil {
			return
		}
		bs := c.(*bitcoin.BTCService)
		// 1. 获取本地lockinInfo
		var localLockinInfo *models.LockinInfo
		localLockinInfo, err = ns.dispatchService.getLockinInfo(m.SCTokenAddress, m.SecretHash)
		if err != nil {
			return
		}
		// 2. 校验
		err = m.VerifySignData(bs, localLockinInfo, privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()))
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unknow message chainName=%s", msg.GetName())
	}
	return
}
