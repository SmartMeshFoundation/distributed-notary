package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
)

// GetNotaryListRequest :
type GetNotaryListRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
}

/*
getNotaryList 公证人列表查询
*/
func (ua *UserAPI) getNotaryList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetNotaryListRequest{
		BaseReq:             api.NewBaseReq(APIUserNameGetNotaryList),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// GetSCTokenListRequest :
type GetSCTokenListRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
}

/*
getSCTokenList 当前支持的侧链Token列表查询
*/
func (ua *UserAPI) getSCTokenList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetSCTokenListRequest{
		BaseReq:             api.NewBaseReq(APIUserNameGetSCTokenList),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}
