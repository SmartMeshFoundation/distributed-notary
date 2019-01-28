package service

import (
	"errors"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	ethevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

func (ns *NotaryService) startDSMAsk(notaryNumNeed int) (sessionID common.Hash, notaryIDs []int, err error) {
	sessionID = utils.NewRandomHash()
	log.Trace(SessionLogMsg(sessionID, "DSMAsk start..."))
	notaryIDs = append(notaryIDs, ns.self.ID)
	for _, notary := range ns.otherNotaries {
		// 这里怎么ask会公平一点
		req := notaryapi.NewDSMAskRequest(sessionID, ns.self)
		_, err2 := ns.SendAndWaitResponse(req, notary.ID)
		if err2 == nil {
			notaryIDs = append(notaryIDs, notary.ID)
		} else {
			log.Warn(SessionLogMsg(sessionID, "notary[%d] refuse DSMAsk : %s", err2.Error()))
		}
		if len(notaryIDs) >= notaryNumNeed {
			break
		}
	}
	if len(notaryIDs) < notaryNumNeed {
		err = fmt.Errorf("no enough notary to sign message, need %d but only got %d", notaryNumNeed, len(notaryIDs))
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
	// 2. 通知其他参与的公正人
	req := notaryapi.NewDSMNotifySelectionRequest(sessionID, ns.self, notaryIDs, privateKeyID, msgToSign)
	ns.Broadcast(req, notaryIDs...)
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
	// 0. 校验私钥ID可用性
	var privateKeyInfo *models.PrivateKeyInfo
	privateKeyInfo, err = ns.db.LoadPrivateKeyInfo(req.PrivateKeyID)
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	// 1. 解析msgToSign
	var msgToSign mecdsa.MessageToSign
	msgToSign, err = parseMessageToSign(req.MsgName, req.MsgToSignTransportBytes)
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	// 2. 校验msgToSign
	err = ns.checkMsgToSign(sessionID, privateKeyInfo, msgToSign, req.GetSenderNotaryID())
	if err != nil {
		ns.unlockSession(sessionID)
		return
	}
	// 3. 构造dsm
	_, err = mecdsa.NewDistributedSignMessage(ns.db, ns.self.ID, msgToSign, sessionID, req.PrivateKeyID, req.NotaryIDs)
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
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	req := notaryapi.NewDSMPhase1BroadcastRequest(sessionID, ns.self, privateKeyID, msg)
	ns.Broadcast(req, dsm.L.S...)
	return
}

func (ns *NotaryService) saveDSMPhase1(sessionID, privateKeyID common.Hash, msg *models.SignBroadcastPhase1, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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

// dsmPhase2Response 这里需要解析response,所以定义专用结构体,否则使用BaseResponse的话,会损失结构信息导致解析失败
type dsmPhase2Response struct {
	ErrorCode api.ErrorCode          `json:"error_code"`
	ErrorMsg  string                 `json:"error_msg"`
	RequestID string                 `json:"request_id"`
	Data      *models.MessageBPhase2 `json:"data,omitempty"`
}

// GetErrorCode :
func (r *dsmPhase2Response) GetErrorCode() api.ErrorCode {
	return r.ErrorCode
}

// GetErrorMsg :
func (r *dsmPhase2Response) GetErrorMsg() string {
	return r.ErrorMsg
}

/*
步骤2应该是同步的,start后阻塞至finish再返回
*/
func (ns *NotaryService) startDSMPhase2(sessionID, privateKeyID common.Hash) (finish bool, err error) {
	ns.lockSession(sessionID)
	log.Trace(SessionLogMsg(sessionID, "DSMPhase2 start..."))
	// 1. 获取dsm
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		ns.unlockSession(sessionID)
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
	for _, notaryID := range dsm.L.S {
		// 按ID分发,并接收返回, 暂时同步???
		if notaryID == ns.self.ID {
			continue
		}
		var resp *api.BaseResponse
		req := notaryapi.NewDSMPhase2MessageARequest(sessionID, ns.self, privateKeyID, msg)
		resp, err = ns.SendAndWaitResponse(req, notaryID)
		if err != nil {
			ns.unlockSession(sessionID)
			return
		}
		var respMsg models.MessageBPhase2
		err = resp.ParseData(&respMsg)
		if err != nil {
			log.Error("parse MessageBPhase2 err =%s \n%s", err.Error(), utils.ToJSONStringFormat(resp))
			ns.unlockSession(sessionID)
			return
		}
		finish, err = dsm.ReceivePhase2MessageB(&respMsg, notaryID)
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

/*
仅生成应答消息,不改变本身数据,无需lock
*/
func (ns *NotaryService) saveDSMPhase2(sessionID, privateKeyID common.Hash, msg *models.MessageA, senderID int) (msgResp *models.MessageBPhase2, err error) {
	// 1. 获取dsm
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
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
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	req := notaryapi.NewDSMPhase3DeltaIRequest(sessionID, ns.self, privateKeyID, msg)
	ns.Broadcast(req, dsm.L.S...)
	return
}

func (ns *NotaryService) saveDSMPhase3(sessionID, privateKeyID common.Hash, msg *models.DeltaPhase3, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	req := notaryapi.NewDSMPhase5A5BProofRequest(sessionID, ns.self, privateKeyID, msg)
	ns.Broadcast(req, dsm.L.S...)
	return
}

func (ns *NotaryService) saveDSMPhase5A5B(sessionID, privateKeyID common.Hash, msg *models.Phase5A, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	req := notaryapi.NewDSMPhase5CProofRequest(sessionID, ns.self, privateKeyID, msg)
	ns.Broadcast(req, dsm.L.S...)
	return
}

func (ns *NotaryService) saveDSMPhase5C(sessionID, privateKeyID common.Hash, msg *models.Phase5C, senderID int) (finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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
	req := notaryapi.NewDSMPhase6ReceiveSIRequest(sessionID, ns.self, privateKeyID, msg)
	ns.Broadcast(req, dsm.L.S...)
	return
}

func (ns *NotaryService) saveDSMPhase6(sessionID, privateKeyID common.Hash, msg share.SPrivKey, senderID int) (signature []byte, finish bool, err error) {
	ns.lockSession(sessionID)
	// 1. 获取dsm
	dsm, err := mecdsa.NewDistributedSignMessageFromDB(ns.db, ns.self.ID, sessionID, privateKeyID)
	if err != nil {
		log.Error("NewDistributedSignMessageFromDB err = %s", err.Error())
		ns.unlockSession(sessionID)
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

/*
解析其他公证人传来的待签名消息
*/
func parseMessageToSign(msgName string, buf []byte) (msg mecdsa.MessageToSign, err error) {
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
	default:
		err = fmt.Errorf("got msg to sign which does't support, maybe attack")
	}
	return
}

/*
签名信息校验,根据收到的消息类型,自己生成一份对应的消息体,并与收到的比对
*/
func (ns *NotaryService) checkMsgToSign(sessionID common.Hash, privateKeyInfo *models.PrivateKeyInfo, msg mecdsa.MessageToSign, senderID int) (err error) {
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
		localLockinInfo, err = ns.dispatchService.getLockinInfo(m.UserRequest.SCTokenAddress, m.UserRequest.SecretHash)
		if err != nil {
			return
		}
		// 2. 获取本地scTokenProxy
		var c chain.Chain
		c, err = ns.dispatchService.getChainByName(smcevents.ChainName)
		scTokenProxy := c.GetContractProxy(localLockinInfo.SCTokenAddress)
		// 2. 校验
		err = m.VerifySignData(scTokenProxy, privateKeyInfo, localLockinInfo)
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
	default:
		err = fmt.Errorf("unknow message name=%s", msg.GetName())
	}
	return
}
