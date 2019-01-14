package api

import (
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

// RequestName :
type RequestName string

// Request 普通请求,不带唯一key不带SCTokenAddress
type Request interface {
	GetRequestID() string
	GetRequestName() RequestName
	GetResponseChan() chan *BaseResponse
	WriteSuccessResponse(data interface{})
	WriteErrorResponse(errorCode ErrorCode, errorMsg ...string)
}

// NotaryRequest 公证人之前的请求,带唯一SessionID
type NotaryRequest interface {
	GetSessionID() common.Hash
	GetSender() common.Address
	GetSenderID() int
	getSignature() []byte
	setSignature(sig []byte)
}

// CrossChainRequest 跨链交易相关请求,带SCTokenAddress
type CrossChainRequest interface {
	GetSCTokenAddress() common.Address
}

// BaseRequest :
type BaseRequest struct {
	RequestID    string      `json:"request_id"` // 方便日志查询
	Name         RequestName `json:"name"`
	responseChan chan *BaseResponse
}

// NewBaseRequest :
func NewBaseRequest(name RequestName) BaseRequest {
	var req BaseRequest
	req.Name = name
	req.RequestID = utils.HPex(utils.NewRandomHash())
	return req
}

// GetRequestName :
func (br *BaseRequest) GetRequestName() RequestName {
	return br.Name
}

// GetRequestID :
func (br *BaseRequest) GetRequestID() string {
	return br.RequestID
}

// GetResponseChan :
func (br *BaseRequest) GetResponseChan() chan *BaseResponse {
	if br.responseChan == nil {
		br.responseChan = make(chan *BaseResponse, 1)
	}
	return br.responseChan
}

// WriteSuccessResponse :
func (br *BaseRequest) WriteSuccessResponse(data interface{}) {
	if br.responseChan == nil {
		br.responseChan = make(chan *BaseResponse, 1)
	}
	select {
	case br.responseChan <- NewSuccessResponse(br.RequestID, data):
	default:
		// never block
	}
}

// WriteErrorResponse :
func (br *BaseRequest) WriteErrorResponse(errorCode ErrorCode, errorMsg ...string) {
	if br.responseChan == nil {
		br.responseChan = make(chan *BaseResponse, 1)
	}
	select {
	case br.responseChan <- NewFailResponse(br.RequestID, errorCode, errorMsg...):
	default:
		// never block
	}
}

// BaseNotaryRequest :
type BaseNotaryRequest struct {
	SessionID common.Hash    `json:"session_id,omitempty"`
	Sender    common.Address `json:"sender,omitempty"`
	SenderID  int            `json:"sender_id"`
	Signature []byte         `json:"signature,omitempty"` // 签名内容req全文json序列化后的字符串
}

// NewBaseNotaryRequest :
func NewBaseNotaryRequest(sessionID common.Hash, sender common.Address, senderID int) BaseNotaryRequest {
	return BaseNotaryRequest{
		SessionID: sessionID,
		Sender:    sender,
		SenderID:  senderID,
	}
}

// GetSender :
func (bnr *BaseNotaryRequest) GetSender() common.Address {
	return bnr.Sender
}

// GetSenderID :
func (bnr *BaseNotaryRequest) GetSenderID() int {
	return bnr.SenderID
}

// GetSessionID :
func (bnr *BaseNotaryRequest) GetSessionID() common.Hash {
	return bnr.SessionID
}
func (bnr *BaseNotaryRequest) getSignature() []byte {
	return bnr.Signature
}

func (bnr *BaseNotaryRequest) setSignature(sig []byte) {
	bnr.Signature = sig
}

// BaseCrossChainRequest :
type BaseCrossChainRequest struct {
	SCTokenAddress common.Address `json:"sc_token_address,omitempty"`
}

//NewBaseCrossChainRequest :
func NewBaseCrossChainRequest(scTokenAddress common.Address) BaseCrossChainRequest {
	return BaseCrossChainRequest{
		SCTokenAddress: scTokenAddress,
	}
}

// GetSCTokenAddress :
func (bcr *BaseCrossChainRequest) GetSCTokenAddress() common.Address {
	return bcr.SCTokenAddress
}
