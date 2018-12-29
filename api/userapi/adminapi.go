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
		BaseRequest: api.NewBaseRequest(APINameGetNotaryList),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// GetPrivateKeyListRequest :
type GetPrivateKeyListRequest struct {
	api.BaseRequest
}

/*
getPrivateKeyList 私钥列表查询
*/
func (ua *UserAPI) getPrivateKeyList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetPrivateKeyListRequest{
		BaseRequest: api.NewBaseRequest(APINameGetPrivateKeyList),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))

}
