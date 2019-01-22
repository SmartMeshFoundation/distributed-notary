package service

import (
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
PKN = PrivateKeyNegotiation
*/

func (ns *NotaryService) startPKNPhase1(privateKeyID common.Hash, receivedMsg *models.KeyGenBroadcastMessage1, senderID int) (finish bool, err error) {
	// 这里主动开始的时候需要锁,被动开始不需要锁,上层已经锁了但是后续解锁流程不变
	if receivedMsg == nil {
		ns.lockSession(privateKeyID)
	}
	log.Trace(SessionLogMsg(privateKeyID, "PKNPhase1 start..."))
	// 1. 初始化KeyGenerator
	keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
	// 2. 生成 KeyGenBroadcastMessage1并广播至所有notary
	var msg *models.KeyGenBroadcastMessage1
	msg, err = keyGenerator.GeneratePhase1PubKeyProof()
	if err != nil {
		ns.unlockSession(privateKeyID)
		return
	}
	// 3. 如果是被动开始,即receivedMsg != nil,保存
	if receivedMsg != nil {
		finish, err = keyGenerator.ReceivePhase1PubKeyProof(receivedMsg, senderID)
		if err != nil {
			ns.unlockSession(privateKeyID)
			return
		}
	}
	ns.unlockSession(privateKeyID)
	err = ns.BroadcastMsg(privateKeyID, notaryapi.APINamePKNPhase1PubKeyProof, msg, true)
	return
}

func (ns *NotaryService) savePKNPhase1Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage1, senderID int) (finish bool, err error) {
	ns.lockSession(keyGenerator.PrivateKeyID)
	finish, err = keyGenerator.ReceivePhase1PubKeyProof(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase1 done..."))
	}
	ns.unlockSession(keyGenerator.PrivateKeyID)
	return
}

func (ns *NotaryService) startPKNPhase2(keyGenerator *mecdsa.ThresholdPrivKeyGenerator) (err error) {
	ns.lockSession(keyGenerator.PrivateKeyID)
	log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase2 start..."))
	var msg *models.KeyGenBroadcastMessage2
	msg, err = keyGenerator.GeneratePhase2PaillierKeyProof()
	if err != nil {
		ns.unlockSession(keyGenerator.PrivateKeyID)
		return
	}
	ns.unlockSession(keyGenerator.PrivateKeyID)
	return ns.BroadcastMsg(keyGenerator.PrivateKeyID, notaryapi.APINamePKNPhase2PaillierKeyProof, msg, true)
}

func (ns *NotaryService) savePKNPhase2Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage2, senderID int) (finish bool, err error) {
	/*
		这里存在先收到某个节点的phase2后收到phase1消息的可能,而phase2消息的处理依赖于同节点发来的phase1消息,如果出现这种情况,就会出现NPE
		目前仅在这一步中发现这种情况,所以先暂时这么处理,测试数据中看来该段等待对总体key生成时间来说影响较小,加不加都在3秒作用.后续再优化
	*/
	//times := 0
	//for {
	//	var p *models.PrivateKeyInfo
	//	p, err = ns.db.LoadPrivateKeyInfo(keyGenerator.PrivateKeyID)
	//	if err != nil {
	//		return
	//	}
	//	if _, ok := p.PubKeysProof1[senderID]; ok || times > 10 {
	//		break
	//	}
	//	times++
	//	time.Sleep(100 * time.Millisecond)
	//	log.Warn("wait for phase1 message of notary %d...", senderID)
	//}
	ns.lockSession(keyGenerator.PrivateKeyID)
	finish, err = keyGenerator.ReceivePhase2PaillierPubKeyProof(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase2 done..."))
	}
	ns.unlockSession(keyGenerator.PrivateKeyID)
	return
}

// 定向
func (ns *NotaryService) startPKNPhase3(keyGenerator *mecdsa.ThresholdPrivKeyGenerator) (err error) {
	ns.lockSession(keyGenerator.PrivateKeyID)
	log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase3 start..."))
	var msgMap map[int]*models.KeyGenBroadcastMessage3
	msgMap, err = keyGenerator.GeneratePhase3SecretShare()
	if err != nil {
		ns.unlockSession(keyGenerator.PrivateKeyID)
		return
	}
	ns.unlockSession(keyGenerator.PrivateKeyID)
	for notaryID, msg := range msgMap {
		// 按ID分别发送phase3消息给其他人
		// 这里虽然是定向发送,但是所有参与者都主动发起SecretShare,所以无需关心返回值,在phase3接口中处理即可
		err2 := ns.SendMsg(keyGenerator.PrivateKeyID, notaryapi.APINamePKNPhase3SecretShare, notaryID, msg, nil)
		if err2 != nil {
			err = err2
			return
		}
	}
	return
}

func (ns *NotaryService) savePKNPhase3Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage3, senderID int) (finish bool, err error) {
	ns.lockSession(keyGenerator.PrivateKeyID)
	finish, err = keyGenerator.ReceivePhase3SecretShare(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase3 done..."))
	}
	ns.unlockSession(keyGenerator.PrivateKeyID)
	return
}

func (ns *NotaryService) startPKNPhase4(keyGenerator *mecdsa.ThresholdPrivKeyGenerator) (err error) {
	ns.lockSession(keyGenerator.PrivateKeyID)
	log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase4 start..."))
	var msg *models.KeyGenBroadcastMessage4
	msg, err = keyGenerator.GeneratePhase4PubKeyProof()
	if err != nil {
		ns.unlockSession(keyGenerator.PrivateKeyID)
		return
	}
	ns.unlockSession(keyGenerator.PrivateKeyID)
	return ns.BroadcastMsg(keyGenerator.PrivateKeyID, notaryapi.APINamePKNPhase4PubKeyProof, msg, true)
}

func (ns *NotaryService) savePKNPhase4Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage4, senderID int) (finish bool, err error) {
	ns.lockSession(keyGenerator.PrivateKeyID)
	finish, err = keyGenerator.ReceivePhase4VerifyTotalPubKey(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase4 done..."))
	}
	ns.unlockSession(keyGenerator.PrivateKeyID)
	return
}
