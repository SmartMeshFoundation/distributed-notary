package service

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
)

// OnRequest restful请求处理
func (cs *CrossChainService) OnRequest(req api.Request) {

	switch r := req.(type) {
	case *userapi.GetLockinStatusRequest:
		cs.onGetLockinStatusRequest(r)
	default:
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	return
}

func (cs *CrossChainService) onLockin() {

}

// Lockin状态查询
func (cs *CrossChainService) onGetLockinStatusRequest(req *userapi.GetLockinStatusRequest) {
	lockinInfo, err := cs.lockinHandler.getLockin(req.SecretHash)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
	}
	req.WriteSuccessResponse(lockinInfo)
}
