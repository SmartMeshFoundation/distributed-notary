package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
)

// RegisterSCTokenRequest :
type RegisterSCTokenRequest struct {
	api.BaseRequest
	MainChainName string `json:"main_chain_name,omitempty"` // 主链名,目前仅支持以太坊
	PrivateKeyID  string `json:"private_key_id,omitempty"`  // 部署合约使用的私钥ID
}

/*
registerNewSCToken :
	注册一个新的侧链token地址,dnotary将完成以下工作:
	1.部署一个主链合约,一个侧链合约
*/
func (ua *UserAPI) registerNewSCToken(w rest.ResponseWriter, r *rest.Request) {
	req := &RegisterSCTokenRequest{
		BaseRequest: api.NewBaseRequest(""),
	}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}
