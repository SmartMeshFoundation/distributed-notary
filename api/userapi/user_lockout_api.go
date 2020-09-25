package userapi

import (
	"fmt"

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
	SCUserAddress common.Address `json:"sc_user_address"` // 侧链PrepareLockout使用的地址,校验用
}

func (ua *UserAPI) mcPrepareLockout(w rest.ResponseWriter, r *rest.Request) {
	//scTokenStr := r.PathParam("sctoken")
	//req := &MCPrepareLockoutRequest{
	//	BaseReq:             api.NewBaseReq(APIUserNameMCPrepareLockout),
	//	BaseReqWithResponse: api.NewBaseReqWithResponse(),
	//	BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
	//}
	req := &MCPrepareLockoutRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, fmt.Sprintf("decode json payload err : %s", err.Error())))
		return
	}
	if req.SecretHash == utils.EmptyHash {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, "secret hash can not be empty"))
		return
	}
	req.NewResponseChan()
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

func (ua *UserAPI) mcPrepareLockout2(w rest.ResponseWriter, r *rest.Request) {
	//scTokenStr := r.PathParam("sctoken")
	//req := &MCPrepareLockoutRequest{
	//	BaseReq:             api.NewBaseReq(APIUserNameMCPrepareLockout),
	//	BaseReqWithResponse: api.NewBaseReqWithResponse(),
	//	BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
	//}
	req := &MCPrepareLockoutRequest2{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, fmt.Sprintf("decode json payload err : %s", err.Error())))
		return
	}
	if req.SecretHash == utils.EmptyHash {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, "secret hash can not be empty"))
		return
	}
	if !req.VerifySign() {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodePermissionDenied, "signature verified failed"))
		return
	}
	req2 := req.toMCPrepareLockoutRequest()
	req2.NewResponseChan()
	ua.SendToServiceWithoutVerifySign(req2)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req2))
}
