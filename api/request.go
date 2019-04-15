package api

import (
	"crypto/ecdsa"
	"encoding/json"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
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
	GetSigner() common.Address
	GetSignature() []byte
	SetSignature(sig []byte)
	Sign(req ReqWithSignature, key *ecdsa.PrivateKey)
	VerifySign(req ReqWithSignature) bool
}

// BaseReqWithSignature 基类
type BaseReqWithSignature struct {
	Signer    common.Address `json:"signer,omitempty"`
	Signature []byte         `json:"signature,omitempty"`
}

// NewBaseReqWithSignature constructor
func NewBaseReqWithSignature(signer common.Address) BaseReqWithSignature {
	return BaseReqWithSignature{
		Signer:    signer,
		Signature: nil,
	}
}

// GetSigner impl ReqWithSignature
func (r *BaseReqWithSignature) GetSigner() common.Address {
	return r.Signer
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
	sig := req.GetSignature()
	req.SetSignature(nil)
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	//fmt.Println("data to sign :", string(data))
	sig, err = utils.SignData(key, data)
	//fmt.Println("sig:", common.Bytes2Hex(sig))
	if err != nil {
		panic(err)
	}
	req.SetSignature(sig)
	return
}

// VerifySign impl ReqWithSignature
func (r *BaseReqWithSignature) VerifySign(req ReqWithSignature) bool {
	if r.GetSigner() == utils.EmptyAddress {
		return false
	}
	sig := req.GetSignature()
	req.SetSignature(nil)
	data, err := json.Marshal(req)
	if err != nil {
		return false
	}
	//fmt.Println("data to verify :", string(data))
	//fmt.Println("sig:", common.Bytes2Hex(sig))
	dataHash := utils.Sha3(data)
	signer, err := utils.Ecrecover(dataHash, sig)
	if err != nil {
		panic(err)
	}
	req.SetSignature(sig)
	if signer == r.GetSigner() {
		return true
	}
	//todo 为了兼容来自浏览器的请求,go相关代码不会走到这里
	sig[64] = 1
	signer, err = utils.Ecrecover(dataHash, sig)
	if err != nil {
		panic(err)
	}
	req.SetSignature(sig)
	return signer == r.GetSigner()
}

/*
ReqWithResponse **********************************************
该类请求提供返回值相关处理方法
*/
type ReqWithResponse interface {
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
