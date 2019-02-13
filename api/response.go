package api

import (
	"encoding/json"
)

// ErrorCode :
type ErrorCode string

/* #no gosec */
const (
	ErrorCodeSuccess          = "0000"
	ErrorCodeDataNotFound     = "1000"
	ErrorCodePermissionDenied = "2000"
	ErrorCodeTimeout          = "3000"
	ErrorCodeParamsWrong      = "4000"
	ErrorCodeException        = "9999"
)

var defaultErrorMsg = "unknown error"

// ErrorCode2MsgMap :
var ErrorCode2MsgMap map[ErrorCode]string

func init() {
	ErrorCode2MsgMap = make(map[ErrorCode]string)
	ErrorCode2MsgMap[ErrorCodeSuccess] = "success"
	ErrorCode2MsgMap[ErrorCodeDataNotFound] = "data not found"
	ErrorCode2MsgMap[ErrorCodePermissionDenied] = "permission denied"
	ErrorCode2MsgMap[ErrorCodeTimeout] = "request time out"
	ErrorCode2MsgMap[ErrorCodeParamsWrong] = "params wrong"
	ErrorCode2MsgMap[ErrorCodeException] = "exception,best call admin"
}

// Response :
type Response interface {
	GetErrorCode() ErrorCode
	GetErrorMsg() string
}

// BaseResponse :
type BaseResponse struct {
	BaseReq
	BaseReqWithResponse
	ErrorCode ErrorCode `json:"error_code"`
	ErrorMsg  string    `json:"error_msg"`
	//Data      interface{} `json:"data,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

// NewSuccessResponse :
func NewSuccessResponse(requestID string, data interface{}) *BaseResponse {
	r := &BaseResponse{
		BaseReq:   NewBaseReq(APINameResponse),
		ErrorCode: ErrorCodeSuccess,
		ErrorMsg:  ErrorCode2MsgMap[ErrorCodeSuccess],
	}
	r.RequestID = requestID
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		r.Data = buf
	}
	return r
}

// NewFailResponse :
func NewFailResponse(requestID string, errorCode ErrorCode, errorMsg ...string) *BaseResponse {
	r := &BaseResponse{
		BaseReq:   NewBaseReq(APINameResponse),
		ErrorCode: errorCode,
	}
	if len(errorMsg) > 0 {
		r.ErrorMsg = errorMsg[0]
	} else {
		var ok bool
		r.ErrorMsg, ok = ErrorCode2MsgMap[errorCode]
		if !ok {
			r.ErrorMsg = defaultErrorMsg
		}
	}
	r.RequestID = requestID
	return r
}

// ParseData :
func (br *BaseResponse) ParseData(to interface{}) (err error) {
	var buf []byte
	buf, err = json.Marshal(br.Data)
	if err != nil {
		return
	}
	return json.Unmarshal(buf, to)
}

// GetErrorCode : impl Response
func (br *BaseResponse) GetErrorCode() ErrorCode {
	return br.ErrorCode
}

// GetErrorMsg : impl Response
func (br *BaseResponse) GetErrorMsg() string {
	return br.ErrorMsg
}
