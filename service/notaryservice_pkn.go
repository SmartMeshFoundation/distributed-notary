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

func (ns *NotaryService) startPKNPhase1(privateKeyID common.Hash) (err error) {
	log.Trace(SessionLogMsg(privateKeyID, "PKNPhase1 start..."))
	// 1. 初始化KeyGenerator
	keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
	// 2. 生成 KeyGenBroadcastMessage1并广播至所有notary
	var msg *models.KeyGenBroadcastMessage1
	msg, err = keyGenerator.GeneratePhase1PubKeyProof()
	if err != nil {
		return
	}
	return ns.BroadcastMsg(privateKeyID, notaryapi.APINamePhase1PubKeyProof, msg, true)
}

func (ns *NotaryService) savePKNPhase1Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage1, senderID int) (finish bool, err error) {
	finish, err = keyGenerator.ReceivePhase1PubKeyProof(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase1 done..."))
	}
	return
}

func (ns *NotaryService) startPKNPhase2(keyGenerator *mecdsa.ThresholdPrivKeyGenerator) (err error) {
	log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase2 start..."))
	var msg *models.KeyGenBroadcastMessage2
	msg, err = keyGenerator.GeneratePhase2PaillierKeyProof()
	if err != nil {
		return
	}
	return ns.BroadcastMsg(keyGenerator.PrivateKeyID, notaryapi.APINAMEPhase2PaillierKeyProof, msg, true)
}

func (ns *NotaryService) savePKNPhase2Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage2, senderID int) (finish bool, err error) {
	finish, err = keyGenerator.ReceivePhase2PaillierPubKeyProof(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase2 done..."))
	}
	return
}

// 定向
func (ns *NotaryService) startPKNPhase3(keyGenerator *mecdsa.ThresholdPrivKeyGenerator) (err error) {
	log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase3 start..."))
	var msgMap map[int]*models.KeyGenBroadcastMessage3
	msgMap, err = keyGenerator.GeneratePhase3SecretShare()
	if err != nil {
		return
	}
	for notaryID, msg := range msgMap {
		// 按ID分别发送phase3消息给其他人
		// 这里虽然是定向发送,但是所有参与者都主动发起SecretShare,所以无需关心返回值,在phase3接口中处理即可 TODO
		err2 := ns.SendMsg(keyGenerator.PrivateKeyID, notaryapi.APINAMEPhase3SecretShare, notaryID, msg)
		if err2 != nil {
			err = err2
			return
		}
	}
	return
}

func (ns *NotaryService) savePKNPhase3Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage3, senderID int) (finish bool, err error) {
	finish, err = keyGenerator.ReceivePhase3SecretShare(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase3 done..."))
	}
	return
}

func (ns *NotaryService) startPKNPhase4(keyGenerator *mecdsa.ThresholdPrivKeyGenerator) (err error) {
	log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase4 start..."))
	var msg *models.KeyGenBroadcastMessage4
	msg, err = keyGenerator.GeneratePhase4PubKeyProof()
	if err != nil {
		return
	}
	return ns.BroadcastMsg(keyGenerator.PrivateKeyID, notaryapi.APINamePhase4PubKeyProof, msg, true)
}

func (ns *NotaryService) savePKNPhase4Msg(keyGenerator *mecdsa.ThresholdPrivKeyGenerator, msg *models.KeyGenBroadcastMessage4, senderID int) (finish bool, err error) {
	finish, err = keyGenerator.ReceivePhase4VerifyTotalPubKey(msg, senderID)
	if finish {
		log.Trace(SessionLogMsg(keyGenerator.PrivateKeyID, "PKNPhase4 done..."))
	}
	return
}
