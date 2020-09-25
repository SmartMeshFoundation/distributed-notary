package userapi

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
{
	"name": "User-SCPrepareLockin",
	"request_id": "0x03bc1c3ddeb1e428b61959b2eefe04e19cc679d1ec86bdb81a4c3d848bcdaa00",
	"sc_token_address": "0x326dee230e67e5c124e9c36eae2126c2158bf361",
	"signer": "AvGCZuCe4L74U5nymitF4QQZHrtbutUkucq6kubGLJ2G",
	"secret_hash": "0x1305497a65c3fdb9f3676548feea03e068662da2ad7f6e2cc23152082e365f4f",
	"mc_user_address": "k9r6wbXp/D3vcpnd/xkIX2z86hg=",
	"mc_expiration": 8644816,
	"signature": "fKQdLOuTSb9cOuNI++Zyg+SAhT4dEvO2u0Nh7WTMcWgdmFmjefdjR9dhqLjm0xyqwCCGO4INvdidghzmCbbKsQA="
}
*/
type SCPrepareLockinRequest2 struct {
	api.BaseReq
	api.BaseReqWithResponse
	api.BaseReqWithSCToken
	Signer        string         `json:"signer,omitempty"`    //hex编码的公钥x
	Signature     string         `json:"signature,omitempty"` //hex编码的签名
	SecretHash    common.Hash    `json:"secret_hash"`
	MCUserAddress common.Address `json:"mc_user_address"` // 主链PrepareLockin使用的地址,校验用
	//MCTXHash       chainhash.Hash `json:"mc_tx_hash,omitempty"`       // 当主链为BTC的时候使用
	MCTXHash       string         `json:"mc_tx_hash,omitempty"` // 当主链为BTC的时候使用
	MCExpiration   *big.Int       `json:"mc_expiration,omitempty"`
	MCLockedAmount btcutil.Amount `json:"mc_locked_amount,omitempty"` // 当主链为BTC的时候使用
}

func (r *SCPrepareLockinRequest2) toSCPrepareLockinRequest() *SCPrepareLockinRequest {
	signer, _ := hex.DecodeString(r.Signer)
	signature, _ := hex.DecodeString(r.Signature)
	mctxhash, _ := hex.DecodeString(r.MCTXHash)
	var r2 = &SCPrepareLockinRequest{
		BaseReq:             r.BaseReq,
		BaseReqWithResponse: r.BaseReqWithResponse,
		BaseReqWithSCToken:  r.BaseReqWithSCToken,
		BaseReqWithSignature: api.BaseReqWithSignature{
			Signer:    signer,
			Signature: signature,
		},
		SecretHash:     r.SecretHash,
		MCUserAddress:  r.MCUserAddress.Bytes(),
		MCTXHash:       mctxhash,
		MCExpiration:   r.MCExpiration,
		MCLockedAmount: r.MCLockedAmount,
	}
	return r2
}

// VerifySign impl ReqWithSignature
func (r *SCPrepareLockinRequest2) VerifySign() bool {
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
	if signerEthAddress == r.GetSignerETHAddress() {
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
	return signerEthAddress == r.GetSignerETHAddress()
}

// GetSignerETHAddress impl ReqWithSignature
func (r *SCPrepareLockinRequest2) GetSignerETHAddress() common.Address {
	if len(r.Signer) == 0 {
		return utils.EmptyAddress
	}
	return crypto.PubkeyToAddress(*r.GetSigner())
}

// GetSigner impl ReqWithSignature
func (r *SCPrepareLockinRequest2) GetSigner() *ecdsa.PublicKey {
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
