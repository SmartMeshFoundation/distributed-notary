package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
)

// GetPrivateKeyListRequest :
type GetPrivateKeyListRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
}

/*
getPrivateKeyList 私钥列表查询
*/
func (ua *UserAPI) getPrivateKeyList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetPrivateKeyListRequest{
		BaseReq:             api.NewBaseReq(APIAdminNameGetPrivateKeyList),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// CreatePrivateKeyRequest :
type CreatePrivateKeyRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
}

/*
CreatePrivateKey 该接口仅仅是发启一次过程,不参与实际协商过程,由SystemService处理
发起一次私钥协商过程,生成一组私钥片
*/
func (ua *UserAPI) createPrivateKey(w rest.ResponseWriter, r *rest.Request) {
	req := &CreatePrivateKeyRequest{
		BaseReq:             api.NewBaseReq(APIAdminNameCreatePrivateKey),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// RegisterSCTokenRequest :
type RegisterSCTokenRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
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
		BaseReq:             api.NewBaseReq(APIAdminNameRegisterNewSCToken),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.MainChainName == "" {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, "main_chain_name can not be null"))
		return
	}
	if req.PrivateKeyID == "" {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, "private_key_id can not be null"))
		return
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}
