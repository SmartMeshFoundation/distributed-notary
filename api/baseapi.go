package api

import (
	"fmt"
	"time"

	"net/http"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
)

// defaultAPITimeout : 默认api请求超时时间
var defaultAPITimeout = 30 * time.Second

/*
BaseAPI : 提供一些公共方法
*/
type BaseAPI struct {
	serverName  string
	host        string
	router      rest.App
	middleWares []rest.Middleware
	api         *rest.Api
	timeout     time.Duration // 调用service层的超时时间
	requestChan chan Request
}

// NewBaseAPI :
func NewBaseAPI(serverName string, host string, router rest.App, middleWares ...rest.Middleware) BaseAPI {
	return BaseAPI{
		serverName:  serverName,
		host:        host,
		router:      router,
		timeout:     defaultAPITimeout,
		middleWares: middleWares,
		requestChan: make(chan Request, 10),
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
func (ba *BaseAPI) GetRequestChan() <-chan Request {
	return ba.requestChan
}

// SetTimeout :
func (ba *BaseAPI) SetTimeout(timeout time.Duration) {
	ba.timeout = timeout
}

// SendToServiceAndWaitResponse :
func (ba *BaseAPI) SendToServiceAndWaitResponse(req Request, timeout ...time.Duration) *BaseResponse {
	apiRequestLog(req)
	if r, ok := req.(NotaryRequest); ok {
		if !VerifyNotarySignature(r) {
			return NewFailResponse(req.GetRequestID(), ErrorCodePermissionDenied)
		}
	}
	var resp *BaseResponse
	requestTimeout := ba.timeout
	if len(timeout) > 0 && timeout[0] > 0 {
		requestTimeout = timeout[0]
	}
	ba.requestChan <- req
	if requestTimeout > 0 {
		select {
		case resp = <-req.GetResponseChan():
		case <-time.After(requestTimeout):
			resp = NewFailResponse(req.GetRequestID(), ErrorCodeTimeout)
		}
	} else {
		resp = <-req.GetResponseChan()
	}
	apiResponseLog(req, resp)
	return resp
}

/*
tool functions
*/

// Return :
func Return(w rest.ResponseWriter, response *BaseResponse) {
	if w == nil {
		return
	}
	err := w.WriteJson(response)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

func apiRequestLog(req Request) {
	prefix := ""
	body := ""
	if nr, ok := req.(NotaryRequest); ok {
		// NotaryRequest 只打印基本信息,否则太多了
		type requestToLog struct {
			BaseRequest
			BaseNotaryRequest
			BaseCrossChainRequest
		}
		var l requestToLog
		l.BaseRequest.RequestID = req.GetRequestID()
		l.BaseRequest.Name = req.GetRequestName()
		l.BaseNotaryRequest.SessionID = nr.GetSessionID()
		l.BaseNotaryRequest.Sender = nr.GetSender()
		l.BaseNotaryRequest.Signature = nr.getSignature()
		body = utils.ToJSONStringFormat(l)
		prefix = fmt.Sprintf("[SessionID=%s] ", utils.HPex(nr.GetSessionID()))
	} else {
		body = utils.ToJSONStringFormat(req)
	}
	log.Trace("%s===> API Request %s :\n%s", prefix, req.GetRequestID(), body)
}

func apiResponseLog(req Request, resp *BaseResponse) {
	prefix := ""
	if nr, ok := req.(NotaryRequest); ok {
		prefix = fmt.Sprintf("[SessionID=%s] ", utils.HPex(nr.GetSessionID()))
	}
	log.Trace(fmt.Sprintf("%s===> API Response %s :\n%s", prefix, req.GetRequestID(), utils.ToJSONStringFormat(resp)))
}
