package service

import (
	"crypto/ecdsa"

	"fmt"

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
	privateKey *ecdsa.PrivateKey
	self       models.NotaryInfo
	notaries   []models.NotaryInfo //这里保存除我以外的notary信息
	db         *models.DB
}

// NewNotaryService :
func NewNotaryService(db *models.DB, privateKey *ecdsa.PrivateKey, allNotaries []models.NotaryInfo) (ns *NotaryService, err error) {
	ns = &NotaryService{
		db:         db,
		privateKey: privateKey,
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
	err = ns.startPKNPhase1(privateKeyID)
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
		err = ns.startPKNPhase1(privateKeyID)
		if err != nil {
			errMsg := SessionLogMsg(privateKeyID, "startPKNPhase1 err = %s", err.Error())
			req.WriteErrorResponse(api.ErrorCodeException, errMsg)
		} else {
			req.WriteSuccessResponse(nil)
		}
		// 这里继续处理,保存已经收到的phase1消息
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
	fmt.Println("==========================", finish)
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
