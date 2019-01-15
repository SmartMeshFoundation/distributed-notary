package userapi

import (
	"crypto/ecdsa"

	"encoding/json"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// GetNotaryListRequest :
type GetNotaryListRequest struct {
	api.BaseRequest
}

/*
getNotaryList 公证人列表查询
*/
func (ua *UserAPI) getNotaryList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetNotaryListRequest{
		BaseRequest: api.NewBaseRequest(APIUserNameGetNotaryList),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// GetSCTokenListRequest :
type GetSCTokenListRequest struct {
	api.BaseRequest
}

/*
getSCTokenList 当前支持的侧链Token列表查询
*/
func (ua *UserAPI) getSCTokenList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetSCTokenListRequest{
		BaseRequest: api.NewBaseRequest(APIUserNameGetSCTokenList),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// GetLockinStatusRequest :
type GetLockinStatusRequest struct {
	api.BaseRequest
	api.BaseCrossChainRequest
	SecretHash common.Hash `json:"secret_hash"`
}

/*
getLockinStatus 查询lockin状态
*/
func (ua *UserAPI) getLockinStatus(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	secretHashStr := r.PathParam("secrethash")
	req := &GetLockinStatusRequest{
		BaseRequest:           api.NewBaseRequest(APIUserNameGetLockinStatus),
		BaseCrossChainRequest: api.NewBaseCrossChainRequest(common.HexToAddress(scTokenStr)),
		SecretHash:            common.HexToHash(secretHashStr),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// SCPrepareLockinRequest :
type SCPrepareLockinRequest struct {
	api.BaseRequest
	api.BaseCrossChainRequest
	SecretHash    common.Hash    `json:"secret_hash"`
	MCUserAddress common.Address `json:"mc_user_address"`     // 主链PrepareLockin使用的地址,校验用
	SCUserAddress common.Address `json:"sc_user_address"`     // 侧链收款地址,即验证签名
	Signature     []byte         `json:"signature,omitempty"` // 用户签名
}

// Sign :
func (r *SCPrepareLockinRequest) Sign(privateKey *ecdsa.PrivateKey) []byte {
	r.Signature = nil
	buf, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	sig, err := utils.SignData(privateKey, buf)
	if err != nil {
		panic(err)
	}
	r.Signature = sig
	return sig
}

// VerifySign :
func (r *SCPrepareLockinRequest) VerifySign() bool {
	sig := r.Signature
	r.Signature = nil
	if sig == nil || len(sig) == 0 {
		return false
	}
	dataWithoutSig, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	dataHash := utils.Sha3(dataWithoutSig)
	signer, err := utils.Ecrecover(dataHash, sig)
	if err != nil {
		panic(err)
	}
	r.Signature = sig
	return signer == r.SCUserAddress
}

func (ua *UserAPI) scPrepareLockin(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	req := &SCPrepareLockinRequest{
		BaseRequest:           api.NewBaseRequest(APIUserNameSCPrepareLockin),
		BaseCrossChainRequest: api.NewBaseCrossChainRequest(common.HexToAddress(scTokenStr)),
	}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.SCUserAddress == utils.EmptyAddress {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.SecretHash == utils.EmptyHash {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if !req.VerifySign() {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodePermissionDenied))
		return
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}
