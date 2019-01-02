package api

import (
	"crypto/ecdsa"
	"encoding/json"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

// NotarySign 公证人通信签名
func NotarySign(req NotaryRequest, privateKey *ecdsa.PrivateKey) []byte {
	buf, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	sig, err := utils.SignData(privateKey, buf)
	if err != nil {
		panic(err)
	}
	req.setSignature(sig)
	return sig
}

// VerifyNotarySignature 验证公证人通信签名
func VerifyNotarySignature(v NotaryRequest) bool {
	if v == nil {
		return false
	}
	sig := v.getSignature()
	if sig == nil || len(sig) == 0 {
		return false
	}
	v.setSignature(nil)
	dataWithoutSig, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	dataHash := utils.Sha3(dataWithoutSig)
	signer, err := utils.Ecrecover(dataHash, sig)
	if err != nil {
		panic(err)
	}
	v.setSignature(sig)
	return signer == v.GetSender()
}
