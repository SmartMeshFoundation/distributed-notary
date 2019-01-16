package userapi

import (
	"crypto/ecdsa"

	"encoding/json"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// GetLockoutStatusRequest :
type GetLockoutStatusRequest struct {
	api.BaseRequest
	api.BaseCrossChainRequest
	SecretHash common.Hash `json:"secret_hash"`
}

/*
getLockoutStatus 查询lockout状态
*/
func (ua *UserAPI) getLockoutStatus(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	secretHashStr := r.PathParam("secrethash")
	req := &GetLockoutStatusRequest{
		BaseRequest:           api.NewBaseRequest(APIUserNameGetLockoutStatus),
		BaseCrossChainRequest: api.NewBaseCrossChainRequest(common.HexToAddress(scTokenStr)),
		SecretHash:            common.HexToHash(secretHashStr),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// MCPrepareLockoutRequest :
type MCPrepareLockoutRequest struct {
	api.BaseRequest
	api.BaseCrossChainRequest
	SecretHash    common.Hash    `json:"secret_hash"`
	MCUserAddress common.Address `json:"mc_user_address"`     // 主链收款地址,即验证签名
	SCUserAddress common.Address `json:"sc_user_address"`     // 侧链PrepareLockout使用的地址,校验用
	Signature     []byte         `json:"signature,omitempty"` // 用户签名
}

// Sign :
func (r *MCPrepareLockoutRequest) Sign(privateKey *ecdsa.PrivateKey) []byte {
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
func (r *MCPrepareLockoutRequest) VerifySign() bool {
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
	return signer == r.MCUserAddress
}

func (ua *UserAPI) mcPrepareLockout(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	req := &MCPrepareLockoutRequest{
		BaseRequest:           api.NewBaseRequest(APIUserNameMCPrepareLockout),
		BaseCrossChainRequest: api.NewBaseCrossChainRequest(common.HexToAddress(scTokenStr)),
	}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if req.MCUserAddress == utils.EmptyAddress {
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
