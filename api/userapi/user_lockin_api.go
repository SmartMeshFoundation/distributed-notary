package userapi

import (
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// GetLockinStatusRequest :
type GetLockinStatusRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	api.BaseReqWithSCToken
	SecretHash common.Hash `json:"secret_hash"`
}

/*
getLockinStatus 查询lockin状态
*/
func (ua *UserAPI) getLockinStatus(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	secretHashStr := r.PathParam("secrethash")
	req := &GetLockinStatusRequest{
		BaseReq:             api.NewBaseReq(APIUserNameGetLockinStatus),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
		SecretHash:          common.HexToHash(secretHashStr),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// SCPrepareLockinRequest :
type SCPrepareLockinRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	api.BaseReqWithSCToken
	api.BaseReqWithSignature
	SecretHash    common.Hash    `json:"secret_hash"`
	MCUserAddress common.Address `json:"mc_user_address"` // 主链PrepareLockin使用的地址,校验用
	SCUserAddress common.Address `json:"sc_user_address"` // 侧链收款地址,即验证签名
}

func (ua *UserAPI) scPrepareLockin(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	req := &SCPrepareLockinRequest{
		BaseReq:             api.NewBaseReq(APIUserNameSCPrepareLockin),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
	}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.SCUserAddress == utils.EmptyAddress {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.SecretHash == utils.EmptyHash {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}
