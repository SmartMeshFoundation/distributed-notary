package api

// ErrorCode :
type ErrorCode string

/* #no gosec */
const (
	ErrorCodeSuccess          = "0000"
	ErrorCodeDataNotFound     = "1000"
	ErrorCodePermissionDenied = "2000"
	ErrorCodeException        = "9999"
)

var defaultErrorMsg = "unknown error"
var errorCode2MsgMap map[ErrorCode]string

func init() {
	errorCode2MsgMap = make(map[ErrorCode]string)
	errorCode2MsgMap[ErrorCodeSuccess] = "success"
	errorCode2MsgMap[ErrorCodeDataNotFound] = "data not found"
	errorCode2MsgMap[ErrorCodePermissionDenied] = "Permission denied"
	errorCode2MsgMap[ErrorCodeException] = "exception,best call admin"
}

// BaseResponse :
type BaseResponse struct {
	ErrorCode ErrorCode   `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data"`
}

// NewSuccessResponse :
func NewSuccessResponse(data interface{}) *BaseResponse {
	r := &BaseResponse{
		ErrorCode: ErrorCodeSuccess,
		ErrorMsg:  errorCode2MsgMap[ErrorCodeSuccess],
	}
	if data != nil {
		r.Data = data
	}
	return r
}

// NewFailResponse :
func NewFailResponse(errorCode ErrorCode, errorMsg ...string) *BaseResponse {
	r := &BaseResponse{
		ErrorCode: errorCode,
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
