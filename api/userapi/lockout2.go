package userapi

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
封装mcPrepareLockout,为了js方便的调用
*/
// MCPrepareLockoutRequest :
type MCPrepareLockoutRequest2 struct {
	api.BaseReq
	api.BaseReqWithResponse
	api.BaseReqWithSCToken
	Signer        string         `json:"signer,omitempty"`
	Signature     string         `json:"signature,omitempty"`
	SecretHash    common.Hash    `json:"secret_hash"`
	SCUserAddress common.Address `json:"sc_user_address"` // 侧链PrepareLockout使用的地址,校验用
}

// VerifySign impl ReqWithSignature
func (r *MCPrepareLockoutRequest2) VerifySign() bool {
	// 清空不参与签名的数据
	sigString := r.Signature
	sig, _ := hex.DecodeString(r.Signature)
	r.Signature = ""
	// 验证签名
	data, err := json.Marshal(r)
	if err != nil {
		return false
	}
	dataHash := utils.Sha3(data)
	log.Trace("data : %s", string(data))
	log.Trace("data hash : %s", dataHash.String())
	publicKey, err := utils.Ecrecover(dataHash, sig)
	if err != nil {
		log.Error(fmt.Sprintf("ecrecover err : %s", err.Error()))
		return false
	}
	r.Signature = sigString
	signerEthAddress := crypto.PubkeyToAddress(*publicKey)
	if signerEthAddress == r.GetSignerSMCAddress() {
		return true
	}
	//todo 为了兼容来自浏览器的请求,go相关代码不会走到这里
	sig[64] = 1
	publicKey, err = utils.Ecrecover(dataHash, sig)
	if err != nil {
		log.Error(fmt.Sprintf("ecrecover err : %s", err.Error()))
		return false
	}
	r.Signature = sigString
	signerEthAddress = crypto.PubkeyToAddress(*publicKey)
	return signerEthAddress == r.GetSignerSMCAddress()
}

// GetSignerSMCAddress impl ReqWithSignature
func (r *MCPrepareLockoutRequest2) GetSignerSMCAddress() common.Address {
	if len(r.Signer) == 0 {
		return utils.EmptyAddress
	}
	return crypto.PubkeyToAddress(*r.GetSigner())
}

// GetSigner impl ReqWithSignature
func (r *MCPrepareLockoutRequest2) GetSigner() *ecdsa.PublicKey {
	key, err := hex.DecodeString(r.Signer)
	if err != nil {
		panic(err)
	}
	pubKey, err := crypto.DecompressPubkey(key)
	if err != nil {
		log.Error(" crypto.DecompressPubkey err : %s", err.Error())
		panic(err)
	}
	return pubKey
}

func (r *MCPrepareLockoutRequest2) toMCPrepareLockoutRequest() *MCPrepareLockoutRequest {
	signer, _ := hex.DecodeString(r.Signer)
	signature, _ := hex.DecodeString(r.Signature)

	var r2 = &MCPrepareLockoutRequest{
		BaseReq:             r.BaseReq,
		BaseReqWithResponse: r.BaseReqWithResponse,
		BaseReqWithSCToken:  r.BaseReqWithSCToken,
		BaseReqWithSignature: api.BaseReqWithSignature{
			Signer:    signer,
			Signature: signature,
		},
		SecretHash:    r.SecretHash,
		SCUserAddress: r.SCUserAddress,
	}
	return r2
}
