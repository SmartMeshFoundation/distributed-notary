package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
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
