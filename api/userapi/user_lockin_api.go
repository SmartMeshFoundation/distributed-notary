package userapi

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/btcsuite/btcutil"
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

// SCPrepareLockinRequest : 用户在主链prepareLockin完成一段时间后,通知相关公证人,需要侧链PrepareLockin
type SCPrepareLockinRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	api.BaseReqWithSCToken
	api.BaseReqWithSignature
	SecretHash    common.Hash `json:"secret_hash"`
	MCUserAddress []byte      `json:"mc_user_address"` // 主链PrepareLockin使用的地址,校验用
	//MCTXHash       chainhash.Hash `json:"mc_tx_hash,omitempty"`       // 当主链为BTC的时候使用
	MCTXHash       []byte         `json:"mc_tx_hash,omitempty"`       // 当主链为BTC的时候使用
	MCExpiration   *big.Int       `json:"mc_expiration,omitempty"`    // 当主链为BTC的时候使用
	MCLockedAmount btcutil.Amount `json:"mc_locked_amount,omitempty"` // 当主链为BTC的时候使用
}

func (ua *UserAPI) scPrepareLockin(w rest.ResponseWriter, r *rest.Request) {
	//scTokenStr := r.PathParam("sctoken")
	//req := &SCPrepareLockinRequest{
	//	BaseReq:             api.NewBaseReq(APIUserNameSCPrepareLockin),
	//	BaseReqWithResponse: api.NewBaseReqWithResponse(),
	//	BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
	//}
	req := &SCPrepareLockinRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, fmt.Sprintf("decode json payload err : %s", err.Error())))
		return
	}
	//if req.SCUserAddress == utils.EmptyAddress {
	//	api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
	//	return
	//}
	if req.SecretHash == utils.EmptyHash {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, "secret hash cat"))
		return
	}
	req.NewResponseChan()
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

/*
对scPrepareLockin封装,方便js处理,现在这种编码方式对js很不优化.
*/
func (ua *UserAPI) scPrepareLockin2(w rest.ResponseWriter, r *rest.Request) {
	//scTokenStr := r.PathParam("sctoken")
	//req := &SCPrepareLockinRequest{
	//	BaseReq:             api.NewBaseReq(APIUserNameSCPrepareLockin),
	//	BaseReqWithResponse: api.NewBaseReqWithResponse(),
	//	BaseReqWithSCToken:  api.NewBaseReqWithSCToken(common.HexToAddress(scTokenStr)),
	//}
	req := &SCPrepareLockinRequest2{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, fmt.Sprintf("decode json payload err : %s", err.Error())))
		return
	}
	if req.SecretHash == utils.EmptyHash {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, "secret hash cat"))
		return
	}
	req.NewResponseChan()
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}
