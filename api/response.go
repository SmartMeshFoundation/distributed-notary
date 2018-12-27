package api

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
var errorCode2MsgMap map[ErrorCode]string

func init() {
	errorCode2MsgMap = make(map[ErrorCode]string)
	errorCode2MsgMap[ErrorCodeSuccess] = "success"
	errorCode2MsgMap[ErrorCodeDataNotFound] = "data not found"
	errorCode2MsgMap[ErrorCodePermissionDenied] = "permission denied"
	errorCode2MsgMap[ErrorCodeTimeout] = "request time out"
	errorCode2MsgMap[ErrorCodeParamsWrong] = "params wrong"
	errorCode2MsgMap[ErrorCodeException] = "exception,best call admin"
}

// BaseResponse :
type BaseResponse struct {
	ErrorCode ErrorCode   `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
}

// NewSuccessResponse :
func NewSuccessResponse(requestID string, data interface{}) *BaseResponse {
	r := &BaseResponse{
		ErrorCode: ErrorCodeSuccess,
		ErrorMsg:  errorCode2MsgMap[ErrorCodeSuccess],
		RequestID: requestID,
	}
	if data != nil {
		r.Data = data
	}
	return r
}

// NewFailResponse :
func NewFailResponse(requestID string, errorCode ErrorCode, errorMsg ...string) *BaseResponse {
	r := &BaseResponse{
		ErrorCode: errorCode,
		RequestID: requestID,
	}
	if len(errorMsg) > 0 {
		r.ErrorMsg = errorMsg[0]
	} else {
		var ok bool
		r.ErrorMsg, ok = errorCode2MsgMap[errorCode]
		if !ok {
			r.ErrorMsg = defaultErrorMsg
		}
	}
	return r
}
