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
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jinzhu/gorm"
	"github.com/nkbai/log"
)

// NotaryService :
type NotaryService struct {
	privateKey     *ecdsa.PrivateKey
	self           models.NotaryInfo
	notaries       []models.NotaryInfo //这里保存除我以外的notary信息
	db             *models.DB
	sessionLockMap map[common.Hash]*sync.Mutex
}

// NewNotaryService :
func NewNotaryService(db *models.DB, privateKey *ecdsa.PrivateKey, allNotaries []models.NotaryInfo) (ns *NotaryService, err error) {
	ns = &NotaryService{
		db:             db,
		privateKey:     privateKey,
		sessionLockMap: make(map[common.Hash]*sync.Mutex),
	}
	// 初始化self, notaries
	selfAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	for _, n := range allNotaries {
		if n.GetAddress() == selfAddress {
			ns.self = n
		} else {
			ns.notaries = append(ns.notaries, n)
		}
	}
	return
}

// OnEvent 链上事件逻辑处理 预留
func (ns *NotaryService) OnEvent(e chain.Event) {

}

// OnRequest restful请求处理
func (ns *NotaryService) OnRequest(req api.Request) {
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
	}
	return
}

/*
主动开始一次私钥协商
*/
func (ns *NotaryService) startNewPrivateKeyNegotiation() (privateKeyID common.Hash, err error) {
	sessionID := utils.NewRandomHash() // 初始化一个会话ID
	privateKeyID = sessionID           // 将会话ID作为私钥Key
	log.Info(SessionLogMsg(privateKeyID, "Private key negotiation BEGIN"))
	_, err = ns.startPKNPhase1(privateKeyID, nil, 0)
	return
}

/*
收到KeyGenerationPhase1MessageRequest, 被动开始一次私钥协商
*/
func (ns *NotaryService) onKeyGenerationPhase1MessageRequest(req *notaryapi.KeyGenerationPhase1MessageRequest) {
	privateKeyID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(privateKeyID, "unknown notary %s, maybe attack", req.GetSender().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 从db获取privateKey
	_, err := ns.db.LoadPrivateKeyInfo(privateKeyID)
	if err == gorm.ErrRecordNotFound {
		// 3. 如果不存在,开始一次私钥协商
		log.Info(SessionLogMsg(privateKeyID, "Private key negotiation BEGIN"))
		var finish bool
		finish, err = ns.startPKNPhase1(privateKeyID, req.Msg, senderNotaryInfo.ID)
		if err != nil {
			errMsg := SessionLogMsg(privateKeyID, "startPKNPhase1 err = %s", err.Error())
			req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		} else {
			req.WriteSuccessResponse(nil)
		}
		if finish {
			// 3.5 这种情况理论上只会在只有两个公证人的时候出现
			keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
			err = ns.startPKNPhase2(keyGenerator)
			if err != nil {
				log.Error(SessionLogMsg(privateKeyID, "startPKNPhase2 err = %s", err.Error()))
			}
		}
		return
	}
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "LoadPrivateKeyInfo err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 4. 如果存在,保存phase1消息
	keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
	finish, err := ns.savePKNPhase1Msg(keyGenerator, req.Msg, senderNotaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "savePKNPhase1Msg err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 保存完毕直接返回成功,防止调用方api阻塞
	req.WriteSuccessResponse(nil)
	// 5. 如果phase1完成,开始phase2
	if finish {
		err = ns.startPKNPhase2(keyGenerator)
		if err != nil {
			log.Error(SessionLogMsg(privateKeyID, "startPKNPhase2 err = %s", err.Error()))
		}
	}
}

/*
收到KeyGenerationPhase2MessageRequest
*/
func (ns *NotaryService) onKeyGenerationPhase2MessageRequest(req *notaryapi.KeyGenerationPhase2MessageRequest) {
	privateKeyID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(privateKeyID, "unknown notary %s, maybe attack", req.GetSender().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 从db获取privateKey
	_, err := ns.db.LoadPrivateKeyInfo(privateKeyID)
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "LoadPrivateKeyInfo err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3. 保存phase2消息
	keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
	finish, err := ns.savePKNPhase2Msg(keyGenerator, req.Msg, senderNotaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "savePKNPhase2Msg err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 保存完毕直接返回成功,防止调用方api阻塞
	req.WriteSuccessResponse(nil)
	// 4. 如果结果,开始phase3
	if finish {
		err = ns.startPKNPhase3(keyGenerator)
		if err != nil {
			log.Error(SessionLogMsg(privateKeyID, "startPKNPhase3 err = %s", err.Error()))
		}
	}
}

/*
收到KeyGenerationPhase3MessageRequest
*/
func (ns *NotaryService) onKeyGenerationPhase3MessageRequest(req *notaryapi.KeyGenerationPhase3MessageRequest) {
	privateKeyID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(privateKeyID, "unknown notary %s, maybe attack", req.GetSender().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 从db获取privateKey
	_, err := ns.db.LoadPrivateKeyInfo(privateKeyID)
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "LoadPrivateKeyInfo err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3. 保存phase3消息
	keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
	finish, err := ns.savePKNPhase3Msg(keyGenerator, req.Msg, senderNotaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "savePKNPhase3Msg err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 保存完毕直接返回成功,防止调用方api阻塞
	req.WriteSuccessResponse(nil)
	// 4. 如果结果,开始phase4
	if finish {
		err = ns.startPKNPhase4(keyGenerator)
		if err != nil {
			log.Error(SessionLogMsg(privateKeyID, "startPKNPhase4 err = %s", err.Error()))
		}
	}
}

/*
收到KeyGenerationPhase4MessageRequest
*/
func (ns *NotaryService) onKeyGenerationPhase4MessageRequest(req *notaryapi.KeyGenerationPhase4MessageRequest) {
	privateKeyID := req.GetSessionID()
	// 1. 校验sender
	senderNotaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(privateKeyID, "unknown notary %s, maybe attack", req.GetSender().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. 从db获取privateKey
	_, err := ns.db.LoadPrivateKeyInfo(privateKeyID)
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "LoadPrivateKeyInfo err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 3. 保存phase4消息
	keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
	finish, err := ns.savePKNPhase4Msg(keyGenerator, req.Msg, senderNotaryInfo.ID)
	if err != nil {
		errMsg := SessionLogMsg(privateKeyID, "savePKNPhase4Msg err = %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		return
	}
	// 4. 如果结束,私钥生成成功,打印日志
	if finish {
		log.Info(SessionLogMsg(privateKeyID, "Private key negotiation END"))
	}
	req.WriteSuccessResponse(nil)
}

/*
主动开始一次签名,并等待最终的签名结果
*/
func (ns *NotaryService) startDistributedSignAndWait(msgToSign mecdsa.MessageToSign, privateKeyInfo *models.PrivateKeyInfo) (signature []byte, err error) {
	// 1. DSMAsk
	var sessionID common.Hash
	var notaryIDs []int
	sessionID, notaryIDs, err = ns.startDSMAsk()
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
		dsm2, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyInfo.Key)
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
			return
		}
		if times%10 == 0 {
			log.Trace(SessionLogMsg(sessionID, "waiting for istributedSignMessage..."))
		}
		times++
	}
}

/*
收到 DSMAskRequest
*/
func (ns *NotaryService) onDSMAskRequest(req *notaryapi.DSMAskRequest) {
	// 1. 校验sender
	_, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(req.GetSessionID(), "unknown notary %s, maybe attack", req.GetSender().String()))
		req.WriteErrorResponse(api.ErrorCodePermissionDenied)
		return
	}
	// 2. TODO 暂时所有公证人都默认愿意参与所有签名工作
	log.Trace(SessionLogMsg(req.GetSessionID(), "accept DSMAsk from %s", utils.APex(req.GetSender())))
	req.WriteSuccessResponse(nil)
	return
}

/*
收到 DSMNotifySelectionRequest
*/
func (ns *NotaryService) onDSMNotifySelectionRequest(req *notaryapi.DSMNotifySelectionRequest) {
	sessionID := req.GetSessionID()
	// 1. 校验sender
	_, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSender().String()))
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
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSender().String()))
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
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSender().String()))
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
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSender().String()))
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
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSender().String()))
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
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSender().String()))
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
	notaryInfo, ok := ns.getNotaryInfoByAddress(req.GetSender())
	if !ok {
		log.Warn(SessionLogMsg(sessionID, "unknown notary %s, maybe attack", req.GetSender().String()))
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
	for _, v := range ns.notaries {
		if v.GetAddress() == addr {
			notaryInfo = &v
			ok = true
			return
		}
	}
	ok = false
	return
}

func (ns *NotaryService) lockSession(sessionID common.Hash) {
	if _, ok := ns.sessionLockMap[sessionID]; !ok {
		ns.sessionLockMap[sessionID] = &sync.Mutex{}
	}
	ns.sessionLockMap[sessionID].Lock()
}
func (ns *NotaryService) unlockSession(sessionID common.Hash) {
	ns.sessionLockMap[sessionID].Unlock()
}
func (ns *NotaryService) removeSessionLock(sessionID common.Hash) {
	delete(ns.sessionLockMap, sessionID)
}
