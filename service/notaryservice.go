package service

import (
	"crypto/ecdsa"
	"time"

	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	mecdsa2 "github.com/SmartMeshFoundation/distributed-notary/service/mecdsa"
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
		ns.onKeyGenerationPhase1MessageRequest(r)
	case *notaryapi.KeyGenerationPhase2MessageRequest:
		ns.onKeyGenerationPhase2MessageRequest(r)
	case *notaryapi.KeyGenerationPhase3MessageRequest:
		ns.onKeyGenerationPhase3MessageRequest(r)
	case *notaryapi.KeyGenerationPhase4MessageRequest:
		ns.onKeyGenerationPhase4MessageRequest(r)
	case *notaryapi.DSMAskRequest:
		ns.onDSMAskRequest(r)
	case *notaryapi.DSMNotifySelectionRequest:
		ns.onDSMNotifySelectionRequest(r)
	case *notaryapi.DSMPhase1BroadcastRequest:
		ns.onDSMPhase1BroadcastRequest(r)
	case *notaryapi.DSMPhase2MessageARequest:
		ns.onDSMPhase2MessageARequest(r)
	case *notaryapi.DSMPhase3DeltaIRequest:
		ns.onDSMPhase3DeltaIRequest(r)
	case *notaryapi.DSMPhase5A5BProofRequest:
		ns.onDSMPhase5A5BProofRequest(r)
	case *notaryapi.DSMPhase5CProofRequest:
		ns.onDSMPhase5CProofRequest(r)
	case *notaryapi.DSMPhase6ReceiveSIRequest:
		ns.onDSMPhase6ReceiveSIRequest(r)
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
	ph := mecdsa2.NewPKNHandler(ns.db, ns.self, otherNotaryIDs, sessionID, ns.notaryClient)
	ns.pknHandlerMap.Store(sessionID, ph)
	// 开始pkn并阻塞等待
	return ph.StartPKNAndWaitFinish(nil)
}

/*
收到KeyGenerationPhase1MessageRequest, 被动开始一次私钥协商
*/
func (ns *NotaryService) onKeyGenerationPhase1MessageRequest(req *notaryapi.KeyGenerationPhase1MessageRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok || senderNotaryInfo.ID != req.GetSenderNotaryID() {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 获取pknHandler
	var otherNotaryIDs []int
	for _, notary := range ns.otherNotaries {
		otherNotaryIDs = append(otherNotaryIDs, notary.ID)
	}
	phInterface, loaded := ns.pknHandlerMap.LoadOrStore(sessionID, mecdsa2.NewPKNHandler(ns.db, ns.self, otherNotaryIDs, sessionID, ns.notaryClient))
	ph := phInterface.(*mecdsa2.PKNHandler)
	if loaded {
		// 已经存在,直接保存
		ph.ReceivePhase1PubKeyProof(req.Msg, req.GetSenderNotaryID())
	} else {
		// 不存在,启动协商线程
		go ph.StartPKNAndWaitFinish(req)
	}
	req.WriteSuccessResponse(nil)
	return
}

/*
收到KeyGenerationPhase2MessageRequest
*/
func (ns *NotaryService) onKeyGenerationPhase2MessageRequest(req *notaryapi.KeyGenerationPhase2MessageRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok || senderNotaryInfo.ID != req.GetSenderNotaryID() {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 获取pkn
	phInterface, ok := ns.pknHandlerMap.Load(sessionID)
	if !ok {
		errMsg := SessionLogMsg(sessionID, "can not find PKNHandler with this sessionID")
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3. 保存信息并返回
	phInterface.(*mecdsa2.PKNHandler).ReceivePhase2PaillierPubKeyProof(req.Msg, req.GetSenderNotaryID())
	req.WriteSuccessResponse(nil)
	return
}

/*
收到KeyGenerationPhase3MessageRequest
*/
func (ns *NotaryService) onKeyGenerationPhase3MessageRequest(req *notaryapi.KeyGenerationPhase3MessageRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok || senderNotaryInfo.ID != req.GetSenderNotaryID() {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 获取pkn
	phInterface, ok := ns.pknHandlerMap.Load(sessionID)
	if !ok {
		errMsg := SessionLogMsg(sessionID, "can not find PKNHandler with this sessionID")
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3. 保存信息并返回
	phInterface.(*mecdsa2.PKNHandler).ReceivePhase3SecretShare(req.Msg, req.GetSenderNotaryID())
	req.WriteSuccessResponse(nil)
	return
}

/*
收到KeyGenerationPhase4MessageRequest
*/
func (ns *NotaryService) onKeyGenerationPhase4MessageRequest(req *notaryapi.KeyGenerationPhase4MessageRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok || senderNotaryInfo.ID != req.GetSenderNotaryID() {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 获取pkn
	phInterface, ok := ns.pknHandlerMap.Load(sessionID)
	if !ok {
		errMsg := SessionLogMsg(sessionID, "can not find PKNHandler with this sessionID")
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3. 保存信息并返回
	phInterface.(*mecdsa2.PKNHandler).ReceivePhase4VerifyTotalPubKey(req.Msg, req.GetSenderNotaryID())
	req.WriteSuccessResponse(nil)
	return
}

/*
主动开始一次签名,并等待最终的签名结果
*/
func (ns *NotaryService) startDistributedSignAndWait(msgToSign mecdsa.MessageToSign, privateKeyInfo *models.PrivateKeyInfo) (signature []byte, sessionID common.Hash, err error) {
	// 1. DSMAsk
	start := time.Now()
	notaryNumNeed := params.ThresholdCount + 1
	if msgToSign.GetName() == messagetosign.SpectrumContractDeployTXDataName {
		// 如果需要签名的是部署合约的tx,则要求所有公证人参与
		notaryNumNeed = params.ShareCount
	}
	var notaryIDs []int
	sessionID, notaryIDs, err = ns.startDSMAsk(notaryNumNeed)
	if err != nil {
		log.Error("startDSMAsk err = %s", err.Error())
		return
	}
	// 2.DSMNotifySelection
	err = ns.startDSMNotifySelection(sessionID, notaryIDs, privateKeyInfo.Key, msgToSign)
	if err != nil {
		log.Error("startDSMNotifySelection err = %s", err.Error())
		return
	}
	// 3. 通知完毕之后,直接开始签名过程
	err = ns.startDSMPhase1(sessionID, privateKeyInfo.Key)
	if err != nil {
		log.Error("startDSMPhase1 err = %s", err.Error())
		return
	}
	// 4. 轮询数据库,等待签名完成
	times := 0
	for {
		time.Sleep(time.Second) // TODO 这里轮询周期设置为多少合适,是否需要设置超时
		dsm2, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyInfo.Key)
		if err2 != nil {
			log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
			return
		}
		var finish bool
		signature, finish, err2 = dsm2.GetFinalSignature()
		if err2 != nil {
			log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
			return
		}
		if finish {
			timeUsed := time.Since(start)
			log.Trace("distributedSignMessage end ,total use %f seconds", timeUsed.Seconds())
			return
		}
		if times%10 == 0 {
			log.Trace(SessionLogMsg(sessionID, "waiting for distributedSignMessage..."))
		}
		times++
	}
}

/*
收到 DSMAskRequest
*/
func (ns *NotaryService) onDSMAskRequest(req *notaryapi.DSMAskRequest) {
	// 1. 校验sender
	_, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(req.GetSessionID(), "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. TODO 暂时所有公证人都默认愿意参与所有签名工作
	log.Trace(SessionLogMsg(req.GetSessionID(), "accept DSMAsk from %s", utils.APex(req.GetSigner())))
	req.WriteSuccessResponse(nil)
	return
}

/*
收到 DSMNotifySelectionRequest
*/
func (ns *NotaryService) onDSMNotifySelectionRequest(req *notaryapi.DSMNotifySelectionRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	_, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 构造dsm
	err := ns.saveDSMNotifySelection(req)
	if err != nil {
		errMsg := SessionLogMsg(sessionID, "saveDSMNotifySelection err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3. 这里先行返回,以免阻塞发起人
	req.WriteSuccessResponse(nil)
	// 4.开始phase1
	err = ns.startDSMPhase1(req.GetSessionID(), req.PrivateKeyID)
	if err != nil {
		log.Error(SessionLogMsg(sessionID, "startDSMPhase1 err = %s", err.Error()))
	}
}

/*
收到 DSMPhase1BroadcastRequest
*/
func (ns *NotaryService) onDSMPhase1BroadcastRequest(req *notaryapi.DSMPhase1BroadcastRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. save
	finish, err := ns.saveDSMPhase1(sessionID, req.PrivateKeyID, req.Msg, notaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(sessionID, "saveDSMPhase1 err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3.先行返回,避免阻塞调用方
	req.WriteSuccessResponse(nil)
	// 4.开始下一步
	if finish {
		finish, err = ns.startDSMPhase2(sessionID, req.PrivateKeyID)
		if err != nil {
			log.Error(SessionLogMsg(sessionID, "startDSMPhase2 err = %s", err.Error()))
			return
		}
		// 5.phase2是同步的,如果完成,直接开始pahse3
		if finish {
			err = ns.startDSMPhase3(sessionID, req.PrivateKeyID)
			if err != nil {
				log.Error(SessionLogMsg(sessionID, "startDSMPhase3 err = %s", err.Error()))
			}
		}
	}
}

/*
收到 DSMPhase2MessageARequest
*/
func (ns *NotaryService) onDSMPhase2MessageARequest(req *notaryapi.DSMPhase2MessageARequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2.
	msgResp, err := ns.saveDSMPhase2(sessionID, req.PrivateKeyID, req.Msg, notaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(sessionID, "saveDSMPhase2 err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	req.WriteSuccessResponse(msgResp)
	return
}

/*
收到 DSMPhase3DeltaIRequest
*/
func (ns *NotaryService) onDSMPhase3DeltaIRequest(req *notaryapi.DSMPhase3DeltaIRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. save
	finish, err := ns.saveDSMPhase3(sessionID, req.PrivateKeyID, req.Msg, notaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(sessionID, "saveDSMPhase3 err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3.先行返回,避免阻塞调用方
	req.WriteSuccessResponse(nil)
	// 4.开始下一步
	if finish {
		err = ns.startDSMPhase5A5B(sessionID, req.PrivateKeyID)
		if err != nil {
			log.Error(SessionLogMsg(sessionID, "startDSMPhase5A5B err = %s", err.Error()))
		}
	}
}

/*
收到 DSMPhase5A5BProofRequest
*/
func (ns *NotaryService) onDSMPhase5A5BProofRequest(req *notaryapi.DSMPhase5A5BProofRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. save
	finish, err := ns.saveDSMPhase5A5B(sessionID, req.PrivateKeyID, req.Msg, notaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(sessionID, "saveDSMPhase5A5B err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3.先行返回,避免阻塞调用方
	req.WriteSuccessResponse(nil)
	// 4.开始下一步
	if finish {
		err = ns.startDSMPhase5C(sessionID, req.PrivateKeyID)
		if err != nil {
			log.Error(SessionLogMsg(sessionID, "startDSMPhase5C err = %s", err.Error()))
		}
	}
}

/*
收到 DSMPhase5CProofRequest
*/
func (ns *NotaryService) onDSMPhase5CProofRequest(req *notaryapi.DSMPhase5CProofRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. save
	finish, err := ns.saveDSMPhase5C(sessionID, req.PrivateKeyID, req.Msg, notaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(sessionID, "saveDSMPhase5C err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3.先行返回,避免阻塞调用方
	req.WriteSuccessResponse(nil)
	// 4.开始下一步
	if finish {
		err = ns.startDSMPhase6(sessionID, req.PrivateKeyID)
		if err != nil {
			log.Error(SessionLogMsg(sessionID, "startDSMPhase6 err = %s", err.Error()))
		}
	}
}

/*
收到 DSMPhase6ReceiveSIRequest
*/
func (ns *NotaryService) onDSMPhase6ReceiveSIRequest(req *notaryapi.DSMPhase6ReceiveSIRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSigner())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSigner().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. save
	signature, finish, err := ns.saveDSMPhase6(sessionID, req.PrivateKeyID, req.Msg, notaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(sessionID, "saveDSMPhase6 err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 4. 如果结束,签名生成成功,打印日志
	if finish {
		log.Info(SessionLogMsg(sessionID, "Distributed sign message END, signature=%s", common.BytesToHash(signature).String()))
	}
	req.WriteSuccessResponse(nil)
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

func (ns *NotaryService) lockSession(sessionID common.Hash) {
	var lock *sync.Mutex
	var ok bool
	ns.sessionLockMapLock.Lock()
	if lock, ok = ns.sessionLockMap[sessionID]; !ok {
		lock = &sync.Mutex{}
		ns.sessionLockMap[sessionID] = lock
	}
	ns.sessionLockMapLock.Unlock()
	lock.Lock()
}
func (ns *NotaryService) unlockSession(sessionID common.Hash) {
	var lock *sync.Mutex
	ns.sessionLockMapLock.Lock()
	lock = ns.sessionLockMap[sessionID]
	ns.sessionLockMapLock.Unlock()
	lock.Unlock()
}
func (ns *NotaryService) removeSessionLock(sessionID common.Hash) {
	ns.sessionLockMapLock.Lock()
	delete(ns.sessionLockMap, sessionID)
	ns.sessionLockMapLock.Unlock()
}
