package api

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// RequestName :
type RequestName string

// Request 普通请求,不带唯一key不带SCTokenAddress
type Request interface {
	GetRequestName() RequestName
	GetRequestID() string
	GetResponseChan() chan *BaseResponse
	WriteSuccessResponse(data interface{})
	WriteErrorResponse(errorCode ErrorCode, errorMsg ...string)
}

// NotaryRequest 公证人之前的请求,带唯一SessionID
type NotaryRequest interface {
	GetSessionID() common.Hash
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
	req.RequestID = fmt.Sprintf("%d", time.Now().Nanosecond())
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
	case br.responseChan <- newSuccessResponse(br.RequestID, data):
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
	case br.responseChan <- newFailResponse(br.RequestID, errorCode, errorMsg...):
	default:
		// never block
	}
}

// BaseNotaryRequest :
type BaseNotaryRequest struct {
	SessionID common.Hash `json:"session_id"`
}

// GetSessionID :
func (bnr *BaseNotaryRequest) GetSessionID() common.Hash {
	return bnr.SessionID
}

// BaseCrossChainRequest :
type BaseCrossChainRequest struct {
	SCTokenAddress common.Address `json:"sc_token_address"`
}

// GetSCTokenAddress :
func (bcr *BaseCrossChainRequest) GetSCTokenAddress() common.Address {
	return bcr.SCTokenAddress
}
