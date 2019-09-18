package rest

import (
	"encoding/json"
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"

	"github.com/SmartMeshFoundation/distributed-notary/api"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

type scPrepareLockinRequest struct {
	SCToken     common.Address
	UserAddress common.Address //主侧链用户地址
	SecretHash  common.Hash
}
type scPrepareLockinResponse struct {
	TxHash common.Hash
	Req    *userapi.SCPrepareLockinRequest
}

//请求帮助构造scPrepareLockin
func scPrepareLockin(w rest.ResponseWriter, r *rest.Request) {
	var err error
	sr := scPrepareLockinResponse{}
	defer func() {
		if err != nil {
			errError(w, err)
		}
		log.Trace(fmt.Sprintf("Restful Api SendMessage ----> scPrepareLockin ,err=%v", err))
	}()
	var req scPrepareLockinRequest
	err = r.DecodeJsonPayload(&req)
	if err != nil {
		return
	}
	screq := &userapi.SCPrepareLockinRequest{
		BaseReq:              api.NewBaseReq(userapi.APIUserNameSCPrepareLockin),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:   api.NewBaseReqWithSCToken(req.SCToken),
		BaseReqWithSignature: api.NewBaseReqWithSignature(),
		SecretHash:           req.SecretHash,
		MCUserAddress:        req.UserAddress[:],
		//SCUserAddress:        req.UserAddress,
	}
	sr.Req = screq
	data, err := json.Marshal(screq)
	if err != nil {
		return
	}
	sr.TxHash = utils.Sha3(data)
	success(w, sr)
	return
}

type mcPrepareLockoutResponse struct {
	TxHash common.Hash
	Req    *userapi.MCPrepareLockoutRequest
}

//请求帮助构造MCPrepareLockoutRequest
func mcPrepareLockout(w rest.ResponseWriter, r *rest.Request) {
	var err error
	sr := mcPrepareLockoutResponse{}
	defer func() {
		if err != nil {
			errError(w, err)
		}
		log.Trace(fmt.Sprintf("Restful Api SendMessage ----> mcPrepareLockout ,err=%v", err))
	}()
	var req scPrepareLockinRequest
	err = r.DecodeJsonPayload(&req)
	if err != nil {
		return
	}
	screq := &userapi.MCPrepareLockoutRequest{
		BaseReq:              api.NewBaseReq(userapi.APIUserNameMCPrepareLockout),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:   api.NewBaseReqWithSCToken(req.SCToken),
		BaseReqWithSignature: api.NewBaseReqWithSignature(),
		SecretHash:           req.SecretHash,
		//MCUserAddress:        req.UserAddress,
		SCUserAddress: req.UserAddress,
	}
	sr.Req = screq
	data, err := json.Marshal(screq)
	if err != nil {
		return
	}
	sr.TxHash = utils.Sha3(data)
	success(w, sr)
	return
}
