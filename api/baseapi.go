package api

import (
	"fmt"
	"time"

	"net/http"

	"os"

	"encoding/json"

	"encoding/binary"

	"bytes"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
	"golang.org/x/net/websocket"
)

// defaultAPITimeout : 默认api请求超时时间
var defaultAPITimeout = 1000 * time.Second

/*
BaseAPI : 提供一些公共方法
*/
type BaseAPI struct {
	serverName  string
	host        string
	router      rest.App
	middleWares []rest.Middleware
	api         *rest.Api
	requestChan chan Req
	Timeout     time.Duration // 调用api的超时时间
}

// NewBaseAPI :
func NewBaseAPI(serverName string, host string, router rest.App, middleWares ...rest.Middleware) BaseAPI {
	return BaseAPI{
		serverName:  serverName,
		host:        host,
		router:      router,
		Timeout:     defaultAPITimeout,
		middleWares: middleWares,
		requestChan: make(chan Req, 10),
	}
}

// Start 启动监听线程
func (ba *BaseAPI) Start(sync bool) {
	ba.api = rest.NewApi()
	ba.api.Use(rest.DefaultCommonStack...)
	if len(ba.middleWares) > 0 {
		ba.api.Use(ba.middleWares...)
	}
	ba.api.SetApp(ba.router)
	log.Info("%s listen at %s", ba.serverName, ba.host)
	if sync {
		err := http.ListenAndServe(ba.host, ba.api.MakeHandler())
		if err != nil {
			log.Error("http server start err : %s", err.Error())
			os.Exit(-1)
		}
	} else {
		go func() {
			err := http.ListenAndServe(ba.host, ba.api.MakeHandler())
			if err != nil {
				log.Error("http server start err : %s", err.Error())
				os.Exit(-1)
			}
		}()
	}
}

// GetRequestChan :
func (ba *BaseAPI) GetRequestChan() <-chan Req {
	return ba.requestChan
}

// SetTimeout :
func (ba *BaseAPI) SetTimeout(timeout time.Duration) {
	ba.Timeout = timeout
}

// SendToService 该方法只负责发送到Service,不考虑返回值
func (ba *BaseAPI) SendToService(req Req) {
	var resp *BaseResponse
	reqWithResponse, needResponse := req.(ReqWithResponse)
	if r, ok := req.(ReqWithSignature); ok {
		if !r.VerifySign(r) {
			resp = NewFailResponse(reqWithResponse.GetRequestID(), ErrorCodePermissionDenied)
			if needResponse {
				reqWithResponse.WriteResponse(resp)
			}
			return
		}
	}
	ba.requestChan <- req
	if _, ok := req.(ReqWithResponse); !ok {
		LogAPI(req, nil)
	}
	return
}

// WaitServiceResponse 阻塞等待Service层对请求的返回
func (ba *BaseAPI) WaitServiceResponse(req ReqWithResponse, timeout ...time.Duration) *BaseResponse {
	var resp *BaseResponse
	requestTimeout := ba.Timeout
	if len(timeout) > 0 && timeout[0] > 0 {
		requestTimeout = timeout[0]
	}
	if requestTimeout > 0 {
		select {
		case resp = <-req.GetResponseChan():
		case <-time.After(requestTimeout):
			resp = NewFailResponse(req.(ReqWithResponse).GetRequestID(), ErrorCodeTimeout)
		}
	} else {
		resp = <-req.GetResponseChan()
	}
	LogAPI(req.(Req), resp)
	return resp
}

/*
tool functions
*/

// HTTPReturnJSON :
func HTTPReturnJSON(w rest.ResponseWriter, response *BaseResponse) {
	if w == nil {
		return
	}
	err := w.WriteJson(response)
	if err != nil {
		log.Warn(fmt.Sprintf("HTTPReturnJSON writejson err %s", err))
	}
}

// WSWriteJSON websocket回写
func WSWriteJSON(ws *websocket.Conn, data interface{}) {
	if ws == nil {
		panic("writeToWSConn ws can not be nil")
	}
	if data == nil {
		return
	}
	buf, err := json.Marshal(data)
	if err != nil {
		log.Error(fmt.Sprintf("WSWriteJSON marshal json err %s", err))
		return
	}
	var length int32
	length = int32(len(buf))
	totalBuf := new(bytes.Buffer)
	err = binary.Write(totalBuf, binary.BigEndian, length)
	if err != nil {
		panic(err)
	}
	err = binary.Write(totalBuf, binary.BigEndian, buf)
	if err != nil {
		panic(err)
	}
	n, err := ws.Write(totalBuf.Bytes())
	if err != nil {
		//log.Error("send data to web socket err : %s, data:\n%s", err.Error(), utils.ToJSONStringFormat(data))
		return
	}
	if n != totalBuf.Len() {
		log.Error("send data to web socket err : data len=%d, send len=%d", totalBuf.Len(), n)
	}
	//fmt.Printf("================== send data length=%d totalBuf.Len()=%d n=%d\n", length, totalBuf.Len(), n)
}

// LogAPI 统一打印接收到的API调用日志
func LogAPI(req Req, resp *BaseResponse) {
	var sessionIDStr, requestIDStr, nameStr, senderIDStr string
	nameStr = fmt.Sprintf("Name=%s ", req.GetRequestName())
	if r2, ok := req.(ReqWithSessionID); ok {
		sessionIDStr = fmt.Sprintf("[SessionID=%s] ", utils.HPex(r2.GetSessionID()))
		senderIDStr = fmt.Sprintf("SenderID=%d ", r2.GetSenderNotaryID())
	}
	if r2, ok := req.(ReqWithResponse); ok {
		requestIDStr = fmt.Sprintf("RequestID=%s ", r2.GetRequestID())
	}
	msg := fmt.Sprintf("%s===> API [ %s%s%s]", sessionIDStr, requestIDStr, nameStr, senderIDStr)
	if resp == nil {
		log.Trace(msg)
	} else if resp.GetErrorCode() == ErrorCodeSuccess {
		log.Trace(fmt.Sprintf("%s deal SUCCESS", msg))
	} else {
		body := utils.ToJSONStringFormat(req)
		log.Error(fmt.Sprintf("%s deal FAIL: \nRequest :\n%s\nResponse :\n%s", msg, body, utils.ToJSONStringFormat(resp)))
	}
}

// LogAPICall 统一打印发送出去的API调用日志
func LogAPICall(req Req, resp *BaseResponse, targetNotaryID ...int) {
	var sessionIDStr, requestIDStr, nameStr, targetIDStr string
	if targetNotaryID != nil && len(targetNotaryID) > 0 {
		targetIDStr = fmt.Sprintf("TargetID=%d ", targetNotaryID[0])
	}
	nameStr = fmt.Sprintf("Name=%s ", req.GetRequestName())
	if r2, ok := req.(ReqWithSessionID); ok {
		sessionIDStr = fmt.Sprintf("[SessionID=%s] ", utils.HPex(r2.GetSessionID()))
	}
	if r2, ok := req.(ReqWithResponse); ok {
		requestIDStr = fmt.Sprintf("RequestID=%s ", r2.GetRequestID())
	}
	msg := fmt.Sprintf("%s===> CALL [ %s%s%s]", sessionIDStr, requestIDStr, nameStr, targetIDStr)
	if resp == nil {
		log.Trace(msg)
	} else if resp.GetErrorCode() == ErrorCodeSuccess {
		log.Trace(fmt.Sprintf("%s deal SUCCESS", msg))
	} else {
		body := utils.ToJSONStringFormat(req)
		log.Error(fmt.Sprintf("%s deal FAIL: \nRequest :\n%s\nResponse :\n%s", msg, body, utils.ToJSONStringFormat(resp)))
	}
}
