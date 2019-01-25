package userapi

import (
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// GetLockoutStatusRequest :
type GetLockoutStatusRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	api.BaseReqWithSCToken
	SecretHash common.Hash `json:"secret_hash"`
}

/*
getLockoutStatus 查询lockout状态
*/
func (ua *UserAPI) getLockoutStatus(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	secretHashStr := r.PathParam("secrethash")
	req := &GetLockoutStatusRequest{
		BaseReq:             api.NewBaseReq(APIUserNameGetLockoutStatus),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
		SecretHash:          common.HexToHash(secretHashStr),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// MCPrepareLockoutRequest :
type MCPrepareLockoutRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	api.BaseReqWithSCToken
	api.BaseReqWithSignature
	SecretHash    common.Hash    `json:"secret_hash"`
	MCUserAddress common.Address `json:"mc_user_address"`     // 主链收款地址,即验证签名
	SCUserAddress common.Address `json:"sc_user_address"`     // 侧链PrepareLockout使用的地址,校验用
	Signature     []byte         `json:"signature,omitempty"` // 用户签名
}

func (ua *UserAPI) mcPrepareLockout(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	req := &MCPrepareLockoutRequest{
		BaseReq:             api.NewBaseReq(APIUserNameMCPrepareLockout),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
	}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.MCUserAddress == utils.EmptyAddress {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.SecretHash == utils.EmptyHash {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	req.BaseReqWithSignature = api.NewBaseReqWithSignature(req.MCUserAddress)
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}
