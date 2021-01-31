package jettradeapi

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
)

const (
	prefix                            = "jettrade-"
	APIJettradeAdminNameDeplyContract = prefix + "-admin-deploy"
	userPrefix                        = prefix + "user-"
	APIJettradeUserIssuePO            = userPrefix + "issuepo"
	APIJettradeUserSignPO             = userPrefix + "signPO"
	APIJettradeUserSignDOFF           = userPrefix + "signDOFF"
	APIJettradeUserSignDOBuyer        = userPrefix + "signDOBuyer"
	APIJettradeUserIssueINV           = userPrefix + "issueINV"
)

type JettradeAPI struct {
	*api.BaseAPI
}

var APIName2URLMap map[string]string

func init() {
	APIName2URLMap = make(map[string]string)
	//APIName2URLMap[APIJettradeAdminNameDeplyContract] = "/api/1/jettrade/admin/deploy"
	APIName2URLMap[APIJettradeUserIssuePO] = "/api/1/jettrade/user/issuepo"
	APIName2URLMap[APIJettradeUserSignPO] = "/api/1/jettrade/user/signpo"
	APIName2URLMap[APIJettradeUserSignDOFF] = "/api/1/jettrade/user/signdoff"
	APIName2URLMap[APIJettradeUserSignDOBuyer] = "/api/1/jettrade/user/signdobuyer"
	APIName2URLMap[APIJettradeUserIssueINV] = "/api/1/jettrade/user/issueinv"
}
func NewJetTradeAPI(base *api.BaseAPI) *JettradeAPI {
	return &JettradeAPI{base}
}
func (ua *JettradeAPI) GetRoute() []*rest.Route {
	return []*rest.Route{
		//	rest.Get(APIName2URLMap[APIJettradeAdminNameDeplyContract], ua.DeployContract),
		rest.Post(APIName2URLMap[APIJettradeUserIssuePO], ua.IssuePOOnSpectrum),
		rest.Post(APIName2URLMap[APIJettradeUserSignPO], ua.SignPOOnEthereum),
		rest.Post(APIName2URLMap[APIJettradeUserSignDOFF], ua.SignDOFFOnEthereum),
		rest.Post(APIName2URLMap[APIJettradeUserSignDOBuyer], ua.SignDOBuyerOnSpectrum),
		rest.Post(APIName2URLMap[APIJettradeUserIssueINV], ua.IssueINVOnEthereum),
	}
}

type JettradeReq interface {
	IsJettradeReq() bool
}
type ShareData struct {
}

func (s *ShareData) IsJettradeReq() bool {
	return true
}

type IssuePOOnSpectrumRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	ShareData
	TokenId      *big.Int       `json:"token_id"`
	DocumentInfo string         `json:"document_info"`
	PONUm        string         `json:"ponum"`
	Buyer        common.Address `json:"buyer"`
	Farmer       common.Address `json:"farmer"`
}
type SignPOONEthereumRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	ShareData
	TokenId *big.Int `json:"token_id"`
}
type SignDOFFOnEthereumRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	ShareData
	TokenId        *big.Int       `json:"token_id"`
	DocumentInfo   string         `json:"document_info"`
	PONUm          string         `json:"ponum"`
	DONum          string         `json:"donum"`
	Buyer          common.Address `json:"buyer"`
	Farmer         common.Address `json:"farmer"`
	FreightForward common.Address `json:ff`
}
type SignDOBuyerOnSpectrumRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	ShareData
	TokenId *big.Int `json:"token_id"`
}
type IssueINVOnEthereumRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	ShareData
	TokenId      *big.Int       `json:"token_id"`
	DocumentInfo string         `json:"document_info"`
	PONUm        string         `json:"ponum"`
	DONum        string         `json:"donum"`
	INVNUm       string         `json:"invnum"`
	Buyer        common.Address `json:"buyer"`
	Farmer       common.Address `json:"farmer"`
}

/*
在以太坊和spectrum上部署Doc721合约
*/
func (ua *JettradeAPI) DeployContract(w rest.ResponseWriter, r *rest.Request) {
	api.HTTPReturnJSON(w, api.NewSuccessResponse("xxx", "Ok"))
}
func (ua *JettradeAPI) doreq(req api.ReqWithResponse, w rest.ResponseWriter, r *rest.Request) {
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.HTTPReturnJSON(w, api.NewFailResponse(req.GetRequestID(), api.ErrorCodeParamsWrong, fmt.Sprintf("decode json payload err : %s", err.Error())))
		return
	}
	req.NewResponseChan()
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}
func (ua *JettradeAPI) IssuePOOnSpectrum(w rest.ResponseWriter, r *rest.Request) {
	req := &IssuePOOnSpectrumRequest{}
	ua.doreq(req, w, r)
}
func (ua *JettradeAPI) SignPOOnEthereum(w rest.ResponseWriter, r *rest.Request) {
	req := &SignPOONEthereumRequest{}
	ua.doreq(req, w, r)
}

func (ua *JettradeAPI) SignDOFFOnEthereum(w rest.ResponseWriter, r *rest.Request) {
	req := &SignDOFFOnEthereumRequest{}
	ua.doreq(req, w, r)
}
func (ua *JettradeAPI) SignDOBuyerOnSpectrum(w rest.ResponseWriter, r *rest.Request) {
	req := &SignDOBuyerOnSpectrumRequest{}
	ua.doreq(req, w, r)
}
func (ua *JettradeAPI) IssueINVOnEthereum(w rest.ResponseWriter, r *rest.Request) {
	req := &IssueINVOnEthereumRequest{}
	ua.doreq(req, w, r)
}
