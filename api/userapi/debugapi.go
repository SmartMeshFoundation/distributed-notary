package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// DebugTransferToAccountRequest  :
type DebugTransferToAccountRequest struct {
	api.BaseRequest
	Account common.Address
}

/*
getPrivateKeyList 私钥列表查询
*/
func (ua *UserAPI) transferToAccount(w rest.ResponseWriter, r *rest.Request) {
	addrStr := r.PathParam("account")
	account := common.HexToAddress(addrStr)
	req := &DebugTransferToAccountRequest{
		BaseRequest: api.NewBaseRequest(APIDebugNameTransferToAccount),
		Account:     account,
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}

// DebugGetAllLockinInfoRequest :
type DebugGetAllLockinInfoRequest struct {
	api.BaseRequest
}

func (ua *UserAPI) getAllLockinInfo(w rest.ResponseWriter, r *rest.Request) {
	req := &DebugGetAllLockinInfoRequest{
		BaseRequest: api.NewBaseRequest(APIDebugNameGetAllLockinInfo),
	}
	api.Return(w, ua.SendToServiceAndWaitResponse(req))
}
