package service

import (
	"errors"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

func (ns *NotaryService) startDSMAsk() (sessionID common.Hash, notaryIDs []int, err error) {
	sessionID = utils.NewRandomHash()
	log.Trace(SessionLogMsg(sessionID, "DSMAsk start..."))
	notaryIDs = append(notaryIDs, ns.self.ID)
	for _, notary := range ns.notaries {
		// 这里怎么ask会公平一点
		_, err = ns.SendMsg(sessionID, notaryapi.APINameDSMAsk, notary.ID, notaryapi.NewDSMAskRequest(sessionID, ns.self.GetAddress()))
		if err == nil {
			notaryIDs = append(notaryIDs, notary.ID)
		}
		if len(notaryIDs) > params.ThresholdCount {
			break
		}
	}
	if len(notaryIDs) <= params.ThresholdCount {
		err = errors.New("no enough notary to sign message")
		return
	}
	log.Trace(SessionLogMsg(sessionID, "DSMAsk done..."))
	return
}

/*
仅签名牵头人会调用,生成DistributedSignMessage
*/
func (ns *NotaryService) startDSMNotifySelection(sessionID common.Hash, notaryIDs []int, privateKeyID common.Hash, msgToSign mecdsa.MessageToSign) (err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMNotifySelection start..."))
	// 1. 构造dsm
	_, err = mecdsa.NewDistributedSignMessage(ns.db, ns.self.ID, msgToSign, sessionID, privateKeyID, notaryIDs)
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	ns.unlockSession(sessionID)
	// 2. 通知其他参与的公正人,这里同步通知
	req := notaryapi.NewDSMNotifySelectionRequest(sessionID, ns.self.GetAddress(), notaryIDs, privateKeyID, msgToSign)
	err = ns.BroadcastMsg(sessionID, notaryapi.APINameDSMNotifySelection, req, true, notaryIDs...)
	if err != nil {
		return
	}
	log.Trace(SessionLogMsg(sessionID, "DSMNotifySelection done..."))
	return
}

/*
仅被邀请签名的公证人会调用,生成DistributedSignMessage
*/
func (ns *NotaryService) saveDSMNotifySelection(req *notaryapi.DSMNotifySelectionRequest) (err error) {
	sessionID := req.GetSessionID()
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMNotifySelection start..."))
	// 1. 构造dsm
	_, err = mecdsa.NewDistributedSignMessage(ns.db, ns.self.ID, req.MsgToSign, sessionID, req.PrivateKeyID, req.NotaryIDs)
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	ns.unlockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMNotifySelection done..."))
	return
}

func (ns *NotaryService) startDSMPhase1(sessionID, privateKeyID common.Hash) (err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase1 start..."))
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2. 生成 SignBroadcastPhase1并广播至参与者
	var msg *models.SignBroadcastPhase1
	msg, err = dsm.GeneratePhase1Broadcast()
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	ns.unlockSession(sessionID)
	// 3. 广播
	req := notaryapi.NewDSMPhase1BroadcastRequest(sessionID, ns.self.GetAddress(), privateKeyID, msg)
	err = ns.BroadcastMsg(dsm.Key, notaryapi.APINameDSMPhase1Broadcast, req, true, dsm.S...)
	return
}

func (ns *NotaryService) saveDSMPhase1(sessionID, privateKeyID common.Hash, msg *models.SignBroadcastPhase1, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	finish, err = dsm.ReceivePhase1Broadcast(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(dsm.Key, "DSMPhase1 done..."))
	}
	ns.unlockSession(sessionID)
	return
}

/*
步骤2应该是同步的,start后阻塞至finish再返回
*/
func (ns *NotaryService) startDSMPhase2(sessionID, privateKeyID common.Hash) (finish bool, err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase2 start..."))
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	var msg *models.MessageA
	msg, err = dsm.GeneratePhase2MessageA()
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	for _, notaryID := range dsm.S {
		// 按ID分发,并接收返回, 暂时同步???
		if notaryID == ns.self.ID {
			continue
		}
		var resp api.BaseResponse
		req := notaryapi.NewDSMPhase2MessageARequest(sessionID, ns.self.GetAddress(), privateKeyID, msg)
		resp, err = ns.SendMsg(dsm.Key, notaryapi.APINameDSMPhase2MessageA, notaryID, req)
		if err != nil {
			ns.unlockSession(sessionID)
			return
		}
		var msgResp models.MessageBPhase2
		err = resp.ParseData(&msgResp)
		if err != nil {
			ns.unlockSession(sessionID)
			return
		}
		finish, err = dsm.ReceivePhase2MessageB(&msgResp, notaryID)
		if err != nil {
			ns.unlockSession(sessionID)
			return
		}
	}
	if !finish {
		err = errors.New("expect finish but not ")
	}
	log.Trace(SessionLogMsg(dsm.Key, "DSMPhase2 done..."))
	ns.unlockSession(sessionID)
	return
}

func (ns *NotaryService) saveDSMPhase2(sessionID, privateKeyID common.Hash, msg *models.MessageA, senderID int) (msgResp *models.MessageBPhase2, err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase2 receive..."))
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2. save
	msgResp, err = dsm.ReceivePhase2MessageA(msg, senderID)
	return
}

func (ns *NotaryService) startDSMPhase3(sessionID, privateKeyID common.Hash) (err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase2 start..."))
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	var msg *models.DeltaPhase3
	msg, err = dsm.GeneratePhase3DeltaI()
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	ns.unlockSession(sessionID)
	// 3. 广播
	req := notaryapi.NewDSMPhase3DeltaIRequest(sessionID, ns.self.GetAddress(), privateKeyID, msg)
	return ns.BroadcastMsg(sessionID, notaryapi.APINameDSMPhase3DeltaI, req, true, dsm.S...)
}

func (ns *NotaryService) saveDSMPhase3(sessionID, privateKeyID common.Hash, msg *models.DeltaPhase3, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	finish, err = dsm.ReceivePhase3DeltaI(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(sessionID, "DSMPhase3 done..."))
	}
	ns.unlockSession(sessionID)
	return
}

func (ns *NotaryService) startDSMPhase5A5B(sessionID, privateKeyID common.Hash) (err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase5A5B start..."))
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	_, err = dsm.GeneratePhase4R()
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	var msg *models.Phase5A
	msg, err = dsm.GeneratePhase5a5bZkProof()
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	ns.unlockSession(sessionID)
	// 3. 广播
	req := notaryapi.NewDSMPhase5A5BProofRequest(sessionID, ns.self.GetAddress(), privateKeyID, msg)
	return ns.BroadcastMsg(sessionID, notaryapi.APINameDSMPhase5A5BProof, req, true, dsm.S...)
}

func (ns *NotaryService) saveDSMPhase5A5B(sessionID, privateKeyID common.Hash, msg *models.Phase5A, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	finish, err = dsm.ReceivePhase5A5BProof(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(sessionID, "DSMPhase5A5B done..."))
	}
	ns.unlockSession(sessionID)
	return
}

func (ns *NotaryService) startDSMPhase5C(sessionID, privateKeyID common.Hash) (err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase5C start..."))
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	var msg *models.Phase5C
	msg, err = dsm.GeneratePhase5CProof()
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	ns.unlockSession(sessionID)
	// 3. 广播
	req := notaryapi.NewDSMPhase5CProofRequest(sessionID, ns.self.GetAddress(), privateKeyID, msg)
	return ns.BroadcastMsg(sessionID, notaryapi.APINameDSMPhase5CProof, req, true, dsm.S...)
}

func (ns *NotaryService) saveDSMPhase5C(sessionID, privateKeyID common.Hash, msg *models.Phase5C, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	finish, err = dsm.ReceivePhase5cProof(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(sessionID, "DSMPhase5C done..."))
	}
	ns.unlockSession(sessionID)
	return
}

func (ns *NotaryService) startDSMPhase6(sessionID, privateKeyID common.Hash) (err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase6 start..."))
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	var msg share.SPrivKey
	msg, err = dsm.Generate5dProof()
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	ns.unlockSession(sessionID)
	// 3. 广播
	req := notaryapi.NewDSMPhase6ReceiveSIRequest(sessionID, ns.self.GetAddress(), privateKeyID, msg)
	return ns.BroadcastMsg(sessionID, notaryapi.APINameDSMPhase6ReceiveSI, req, true, dsm.S...)
}

func (ns *NotaryService) saveDSMPhase6(sessionID, privateKeyID common.Hash, msg share.SPrivKey, senderID int) (signature []byte, finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err2 := mecdsa.NewDistributedSignMessageFromDB(ns.db, sessionID, privateKeyID)
	if err2 != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		return
	}
	// 2.
	signature, finish, err = dsm.RecevieSI(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(sessionID, "DSMPhase6 done..."))
	}
	ns.unlockSession(sessionID)
	return
}
