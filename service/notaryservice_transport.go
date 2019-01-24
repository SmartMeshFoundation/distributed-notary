package service

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
)

// Broadcast 广播消息,不关心返回值
func (ns *NotaryService) Broadcast(req api.Req, notaryIDs ...int) {
	if reqWithSignature, ok := req.(api.ReqWithSignature); ok {
		reqWithSignature.Sign(ns.privateKey)
	}
	if len(notaryIDs) > 0 {
		for _, notaryID := range notaryIDs {
			if notaryID == ns.self.ID {
				continue
			}
			ns.notaryClient.SendWSReqToNotary(req, notaryID)
		}
		return
	}
	for _, notary := range ns.notaries {
		if notary.ID == ns.self.ID {
			continue
		}
		ns.notaryClient.SendWSReqToNotary(req, notary.ID)
	}
	return
}

// Send 发送消息,不关心返回值
func (ns *NotaryService) Send(req api.Req, targetNotaryID int) {
	if reqWithSignature, ok := req.(api.ReqWithSignature); ok {
		reqWithSignature.Sign(ns.privateKey)
	}
	ns.notaryClient.SendWSReqToNotary(req, targetNotaryID)
}

// SendAndWaitResponse 发送并阻塞等待返回
func (ns *NotaryService) SendAndWaitResponse(req api.Req, targetNotaryID int) (resp *api.BaseResponse, err error) {
	if reqWithSignature, ok := req.(api.ReqWithSignature); ok {
		reqWithSignature.Sign(ns.privateKey)
	}
	ns.notaryClient.SendWSReqToNotary(req, targetNotaryID)
	return ns.notaryClient.WaitWSResponse(req.(api.ReqWithResponse).GetRequestID())
}
