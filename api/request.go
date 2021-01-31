package api

import (
	"crypto/ecdsa"
	"encoding/json"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nkbai/log"
)

// RequestName :
type RequestName string

// #nosec
const (
	APINameResponse = "Response" // 其他节点对我发过去的请求的应答消息,根据RequestID查询内存中暂存的req,然后处理
)

/*
Req **********************************************
基础请求
*/
type Req interface {
	GetRequestName() RequestName
}

// BaseReq 基类
type BaseReq struct {
	Name RequestName `json:"name"`
}

// NewBaseReq constructor
func NewBaseReq(name RequestName) BaseReq {
	var req BaseReq
	req.Name = name
	return req
}

// GetRequestName impl BaseReq
func (r *BaseReq) GetRequestName() RequestName {
	return r.Name
}

/*
ReqWithSCToken **********************************************
该类请求需要CrossChainService处理
*/
type ReqWithSCToken interface {
	GetSCTokenAddress() common.Address
}

// BaseReqWithSCToken 基类
type BaseReqWithSCToken struct {
	SCTokenAddress common.Address `json:"sc_token_address,omitempty"`
}

// NewBaseReqWithSCToken constructor
func NewBaseReqWithSCToken(scTokenAddress common.Address) BaseReqWithSCToken {
	return BaseReqWithSCToken{
		SCTokenAddress: scTokenAddress,
	}
}

// GetSCTokenAddress impl BaseReqWithSCToken
func (r *BaseReqWithSCToken) GetSCTokenAddress() common.Address {
	return r.SCTokenAddress
}

/*
ReqWithSignature **********************************************
该类请求在接收时需要校验Sender签名
*/
type ReqWithSignature interface {
	GetSigner() *ecdsa.PublicKey
	SetSigner(publicKey *ecdsa.PublicKey)
	GetSignature() []byte
	SetSignature(sig []byte)
	Sign(req ReqWithSignature, key *ecdsa.PrivateKey)
	VerifySign(req ReqWithSignature) bool

	GetSignerETHAddress() common.Address
	GetSignerBTCPublicKey(net *chaincfg.Params) *btcutil.AddressPubKey
}

// BaseReqWithSignature 基类
type BaseReqWithSignature struct {
	Signer    []byte `json:"signer,omitempty"`
	Signature []byte `json:"signature,omitempty"`
}

// NewBaseReqWithSignature constructor
func NewBaseReqWithSignature() BaseReqWithSignature {
	return BaseReqWithSignature{
		Signer:    nil,
		Signature: nil,
	}
}

// GetSigner impl ReqWithSignature
func (r *BaseReqWithSignature) GetSigner() *ecdsa.PublicKey {
	pubKey, err := crypto.DecompressPubkey(r.Signer)
	if err != nil {
		log.Error(" crypto.DecompressPubkey err : %s", err.Error())
	}
	return pubKey
}

// SetSigner impl ReqWithSignature
func (r *BaseReqWithSignature) SetSigner(publicKey *ecdsa.PublicKey) {
	r.Signer = crypto.CompressPubkey(publicKey)

}

// GetSignature impl ReqWithSignature
func (r *BaseReqWithSignature) GetSignature() []byte {
	return r.Signature
}

// SetSignature impl ReqWithSignature
func (r *BaseReqWithSignature) SetSignature(sig []byte) {
	r.Signature = sig
}

// Sign impl ReqWithSignature
func (r *BaseReqWithSignature) Sign(req ReqWithSignature, key *ecdsa.PrivateKey) {
	// 清空不参与签名的数据
	req.SetSignature(nil)
	// 填入数据
	req.SetSigner(&key.PublicKey)
	// 签名
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	//fmt.Println("data to sign :", string(data))
	sig, err := utils.SignData(key, data)
	//fmt.Println("sig:", common.Bytes2Hex(sig))
	if err != nil {
		panic(err)
	}
	req.SetSignature(sig)
	return
}

// VerifySign impl ReqWithSignature
func (r *BaseReqWithSignature) VerifySign(req ReqWithSignature) bool {
	// 清空不参与签名的数据
	sig := req.GetSignature()
	req.SetSignature(nil)
	// 验证签名
	data, err := json.Marshal(req)
	if err != nil {
		return false
	}
	dataHash := utils.Sha3(data)
	log.Trace("data : %s", string(data))
	log.Trace("data hash : %s", dataHash.String())
	//fmt.Println("data to verify :", string(data))
	//fmt.Println("data hash :", dataHash.String())
	//fmt.Println("sig:", common.Bytes2Hex(sig))
	publicKey, err := utils.Ecrecover(dataHash, sig)
	if err != nil {
		log.Error("ecrecover err : %s", err.Error())
		return false
	}
	req.SetSignature(sig)
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
	req.SetSignature(sig)
	signerEthAddress = crypto.PubkeyToAddress(*publicKey)
	return signerEthAddress == r.GetSignerETHAddress()
}

// GetSignerETHAddress impl ReqWithSignature
func (r *BaseReqWithSignature) GetSignerETHAddress() common.Address {
	if r.Signer == nil {
		return utils.EmptyAddress
	}
	return crypto.PubkeyToAddress(*r.GetSigner())
}

// GetSignerBTCPublicKey impl ReqWithSignature
func (r *BaseReqWithSignature) GetSignerBTCPublicKey(net *chaincfg.Params) *btcutil.AddressPubKey {
	if r.Signer == nil {
		return nil
	}
	pubKey := (*btcec.PublicKey)(r.GetSigner())
	addressPubKey, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), net)
	if err != nil {
		panic(err)
	}
	return addressPubKey
}

/*
ReqWithResponse **********************************************
该类请求提供返回值相关处理方法
*/
type ReqWithResponse interface {
	Req
	GetRequestID() string
	GetResponseChan() chan *BaseResponse
	NewResponseChan()
	WriteResponse(resp *BaseResponse)
	WriteSuccessResponse(data interface{})
	WriteErrorResponse(errorCode ErrorCode, errorMsg ...string)
}

// BaseReqWithResponse 基类
type BaseReqWithResponse struct {
	RequestID    string `json:"request_id"`
	responseChan chan *BaseResponse
}

// NewBaseReqWithResponse constructor
func NewBaseReqWithResponse() BaseReqWithResponse {
	return BaseReqWithResponse{
		RequestID:    utils.NewRandomHash().String(),
		responseChan: make(chan *BaseResponse, 1),
	}
}

// GetRequestID impl ReqWithResponse
func (r *BaseReqWithResponse) GetRequestID() string {
	return r.RequestID
}

// GetResponseChan impl ReqWithResponse
func (r *BaseReqWithResponse) GetResponseChan() chan *BaseResponse {
	return r.responseChan
}

// NewResponseChan :
func (r *BaseReqWithResponse) NewResponseChan() {
	r.responseChan = make(chan *BaseResponse, 1)
}

// WriteResponse impl ReqWithResponse
func (r *BaseReqWithResponse) WriteResponse(resp *BaseResponse) {
	select {
	case r.responseChan <- resp:
	default:
		log.Error(fmt.Sprintf("response of requestID = %s double write", r.RequestID))
		panic("never block")
		// never block
	}
}

// WriteSuccessResponse impl ReqWithResponse
func (r *BaseReqWithResponse) WriteSuccessResponse(data interface{}) {
	writeFail := false
	select {
	case r.responseChan <- NewSuccessResponse(r.RequestID, data):
	default:
		writeFail = true
	}
	if writeFail {
		log.Error("responseChan full with requestID=%s", r.RequestID)
	}
}

// WriteErrorResponse impl ReqWithResponse
func (r *BaseReqWithResponse) WriteErrorResponse(errorCode ErrorCode, errorMsg ...string) {
	//if r.responseChan == nil {
	//	r.responseChan = make(chan *BaseResponse, 1)
	//}
	select {
	case r.responseChan <- NewFailResponse(r.RequestID, errorCode, errorMsg...):
	default:
		log.Error("WriteErrorResponse lost,errcode=%d,errorMsg=%s", errorCode, errorMsg)
		// never block
	}
}

/*
ReqWithSessionID **********************************************
该类请求需要NotaryService处理
*/
type ReqWithSessionID interface {
	GetSessionID() common.Hash
	GetSenderNotaryID() int
}

// BaseReqWithSessionID 基类
type BaseReqWithSessionID struct {
	SessionID      common.Hash `json:"session_id,omitempty"`
	SenderNotaryID int         `json:"sender_notary_id"`
}

// NewBaseReqWithSessionID constructor
func NewBaseReqWithSessionID(sessionID common.Hash, senderNotaryID int) BaseReqWithSessionID {
	return BaseReqWithSessionID{
		SessionID:      sessionID,
		SenderNotaryID: senderNotaryID,
	}
}

// GetSenderNotaryID impl ReqWithSessionID
func (r *BaseReqWithSessionID) GetSenderNotaryID() int {
	return r.SenderNotaryID
}

// GetSessionID impl ReqWithSessionID
func (r *BaseReqWithSessionID) GetSessionID() common.Hash {
	return r.SessionID
}
