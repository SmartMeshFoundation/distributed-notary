package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// GetNotaryListRequest :
type GetNotaryListRequest struct {
	api.BaseRequest
}

/*
getNotaryList 公证人列表查询
*/
func (ua *UserAPI) getNotaryList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetNotaryListRequest{
		BaseRequest: api.NewBaseRequest(APIUserNameGetNotaryList),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// GetSCTokenListRequest :
type GetSCTokenListRequest struct {
	api.BaseRequest
}

/*
getSCTokenList 当前支持的侧链Token列表查询
*/
func (ua *UserAPI) getSCTokenList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetSCTokenListRequest{
		BaseRequest: api.NewBaseRequest(APIUserNameGetSCTokenList),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// GetLockinStatusRequest :
type GetLockinStatusRequest struct {
	api.BaseRequest
	api.BaseCrossChainRequest
	SecretHash common.Hash
}

/*
getLockinStatus 查询lockin状态
*/
func (ua *UserAPI) getLockinStatus(w rest.ResponseWriter, r *rest.Request) {
	scTokenStr := r.PathParam("sctoken")
	secretHashStr := r.PathParam("secrethash")
	req := &GetLockinStatusRequest{
		BaseRequest:           api.NewBaseRequest(APIUserNameGetLockinStatus),
		BaseCrossChainRequest: api.NewBaseCrossChainRequest(common.HexToAddress(scTokenStr)),
		SecretHash:            common.HexToHash(secretHashStr),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}
