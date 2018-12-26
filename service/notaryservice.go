package service

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

// NotaryService TODO
type NotaryService struct {
	self     models.NotaryInfo
	notaries []models.NotaryInfo //这里保存除我以外的notary信息
	db       *models.DB
}

// NewNotaryService :
func NewNotaryService(db *models.DB) (ns *NotaryService, err error) {
	ns = &NotaryService{
		db: db,
	}
	// TODO 初始化self, notaries
	return
}

// OnRequest restful请求处理
func (ns *NotaryService) OnRequest(req api.Request) {
	//TODO
	//switch r := req.(type) {
	//case *userapi.CreatePrivateKeyRequest:
	//	ns.onCreatePrivateKeyRequest(r)
	//}
	return
}

/*
主动开始一次私钥协商
*/
func (ns *NotaryService) startNewPrivateKeyNegotiation() (privateKeyID common.Hash, err error) {
	sessionID := utils.NewRandomHash() // 初始化会话ID
	privateKeyID = sessionID           // 将会话ID作为私钥Key
	// 1.初始化KeyGenerator
	keyGenerator := mecdsa.NewThresholdPrivKeyGenerator(ns.self.ID, ns.db, privateKeyID)
	// 2.生成 KeyGenBroadcastMessage1并广播至所有notary
	var msg *models.KeyGenBroadcastMessage1
	msg, err = keyGenerator.GeneratePhase1PubKeyProof()
	if err != nil {
		return
	}
	err = ns.BroadcastMsg(notaryapi.APINamePhase1PubKeyProof, msg, true)
	return
}
