package nonceapi

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

// APIName :
type APIName string

/* #nosec */
const (
	APINameApplyNonce        = "ApplyNonce"        // 申请一个可用nonce
	APINameNotifyNonceStatus = "NotifyNonceStatus" // 上报一个nonce的状态
)

// APIName2URLMap :
var APIName2URLMap map[string]string

func init() {
	APIName2URLMap = make(map[string]string)

	APIName2URLMap[APINameApplyNonce] = "/api/1/apply-nonce"
	APIName2URLMap[APINameNotifyNonceStatus] = "/api/1/:chain/:account/:nonce"
}

var defaultAPITimeout = 30 * time.Second

// NonceServerAPI :
type NonceServerAPI struct {
	api.BaseAPI
}

// NewNonceServerAPI :
func NewNonceServerAPI(host string) *NonceServerAPI {
	var ns NonceServerAPI
	router, err := rest.MakeRouter(
		/*
			nsAPI
		*/
		rest.Post(APIName2URLMap[APINameApplyNonce], ns.ApplyNonce),
		//rest.Post(APIName2URLMap[APINameNotifyNonceStatus], ns.notifyNonceStatus),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	ns.BaseAPI = api.NewBaseAPI("Nonce-Server", host, router)
	ns.Timeout = defaultAPITimeout
	return &ns
}

// ApplyNonceReq :
type ApplyNonceReq struct {
	api.BaseReq
	api.BaseReqWithResponse
	ChainName string `json:"chain_name"`
	Account   string `json:"account"`
	CancelURL string `json:"cancel_url"`
}

// NewApplyNonceReq :
func NewApplyNonceReq(chainName string, account common.Address, cancelURL string) *ApplyNonceReq {
	return &ApplyNonceReq{
		BaseReq:             api.NewBaseReq(APINameApplyNonce),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		ChainName:           chainName,
		Account:             account.String(),
		CancelURL:           cancelURL,
	}
}

// ApplyNonceResponse :
type ApplyNonceResponse struct {
	Nonce uint64 `json:"nonce"`
}

/*
ApplyNonce :
*/
func (nsAPI *NonceServerAPI) ApplyNonce(w rest.ResponseWriter, r *rest.Request) {
	req := &ApplyNonceReq{
		BaseReq:             api.NewBaseReq(APINameApplyNonce),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		log.Warn("got apply-nonce request but parse err : %s", err.Error())
		api.HTTPReturnJSON(w, api.NewFailResponse(req.GetRequestID(), api.ErrorCodeParamsWrong))
		return
	}
	nsAPI.SendToService(req)
	api.HTTPReturnJSON(w, nsAPI.WaitServiceResponse(req))
}

//func (nsAPI *NonceServerAPI) notifyNonceStatus(w rest.ResponseWriter, r *rest.Request) {
//
//}
