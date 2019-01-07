package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
)

// GetPrivateKeyListRequest :
type GetPrivateKeyListRequest struct {
	api.BaseRequest
}

/*
getPrivateKeyList 私钥列表查询
*/
func (ua *UserAPI) getPrivateKeyList(w rest.ResponseWriter, r *rest.Request) {
	req := &GetPrivateKeyListRequest{
		BaseRequest: api.NewBaseRequest(APIAdminNameGetPrivateKeyList),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}
