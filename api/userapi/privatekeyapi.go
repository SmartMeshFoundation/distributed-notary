package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
)

// CreatePrivateKeyRequest :
type CreatePrivateKeyRequest struct {
	api.BaseRequest
}

/*
CreatePrivateKey 该接口仅仅是发启一次过程,不参与实际协商过程,由SystemService处理
发起一次私钥协商过程,生成一组私钥片
*/
func (ua *UserAPI) createPrivateKey(w rest.ResponseWriter, r *rest.Request) {
	req := &CreatePrivateKeyRequest{
		BaseRequest: api.NewBaseRequest(APIAdminNameCreatePrivateKey),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}
