package api

import (
	"fmt"
	"net/http"
	"time"

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
	host        string
	router      rest.App
	middleWares []rest.Middleware
	api         *rest.Api
	timeout     time.Duration // 调用service层的超时时间
	requestChan chan Request
}

// NewBaseAPI :
func NewBaseAPI(host string, router rest.App, middleWares ...rest.Middleware) BaseAPI {
	return BaseAPI{
		host:        host,
		router:      router,
		timeout:     defaultAPITimeout,
		middleWares: middleWares,
		requestChan: make(chan Request, 10),
	}
}

// Start 启动监听线程
func (ba *BaseAPI) Start() {
	ba.api = rest.NewApi()
	ba.api.Use(rest.DefaultDevStack...)
	if len(ba.middleWares) > 0 {
		ba.api.Use(ba.middleWares...)
	}
	ba.api.SetApp(ba.router)
	log.Crit(fmt.Sprintf("http listen and serve :%s", http.ListenAndServe(ba.host, ba.api.MakeHandler())))
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
	log.Trace(fmt.Sprintf("API Request %s :\n%s", req.GetRequestID(), utils.ToJsonStringFormat(req)))
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
			resp = newFailResponse(req.GetRequestID(), ErrorCodeTimeout)
		}
	} else {
		resp = <-req.GetResponseChan()
	}
	log.Trace(fmt.Sprintf("API Response %s :\n%s", req.GetRequestID(), utils.ToJsonStringFormat(resp)))
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
