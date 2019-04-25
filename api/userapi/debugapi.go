package userapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// DebugTransferToAccountRequest  :
type DebugTransferToAccountRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
	Account common.Address
}

/*
getPrivateKeyList 私钥列表查询
*/
func (ua *UserAPI) transferToAccount(w rest.ResponseWriter, r *rest.Request) {
	addrStr := r.PathParam("account")
	account := common.HexToAddress(addrStr)
	req := &DebugTransferToAccountRequest{
		BaseReq:             api.NewBaseReq(APIDebugNameTransferToAccount),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		Account:             account,
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// DebugGetAllLockinInfoRequest :
type DebugGetAllLockinInfoRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
}

func (ua *UserAPI) getAllLockinInfo(w rest.ResponseWriter, r *rest.Request) {
	req := &DebugGetAllLockinInfoRequest{
		BaseReq:             api.NewBaseReq(APIDebugNameGetAllLockinInfo),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// DebugGetAllLockoutInfoRequest :
type DebugGetAllLockoutInfoRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
}

func (ua *UserAPI) getAllLockoutInfo(w rest.ResponseWriter, r *rest.Request) {
	req := &DebugGetAllLockoutInfoRequest{
		BaseReq:             api.NewBaseReq(APIDebugNameGetAllLockoutInfo),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}

// DebugGetAllBTCUtxoRequest :
type DebugGetAllBTCUtxoRequest struct {
	api.BaseReq
	api.BaseReqWithResponse
}

func (ua *UserAPI) getAllBTCUtxo(w rest.ResponseWriter, r *rest.Request) {
	req := &DebugGetAllBTCUtxoRequest{
		BaseReq:             api.NewBaseReq(APIDebugNameGetAllUtxo),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
	}
	ua.SendToService(req)
	api.HTTPReturnJSON(w, ua.WaitServiceResponse(req))
}
