package service

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
)

// OnRequest restful请求处理
func (cs *CrossChainService) OnRequest(req api.Request) {

	switch r := req.(type) {
	// lockin
	case *userapi.GetLockinStatusRequest:
		cs.onGetLockinStatusRequest(r)
	case *userapi.SCPrepareLockinRequest:
		cs.onSCPrepareLockinRequest(r)
	// lockout
	case *userapi.GetLockoutStatusRequest:
		cs.onGetLockoutStatusRequest(r)
	case *userapi.MCPrepareLockoutRequest:
		cs.onMCPrepareLockoutRequest(r)
	default:
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	return
}

// Lockin状态查询
func (cs *CrossChainService) onGetLockinStatusRequest(req *userapi.GetLockinStatusRequest) {
	lockinInfo, err := cs.lockinHandler.getLockin(req.SecretHash)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	req.WriteSuccessResponse(lockinInfo)
}

// Lockout状态查询
func (cs *CrossChainService) onGetLockoutStatusRequest(req *userapi.GetLockoutStatusRequest) {
	lockinInfo, err := cs.lockoutHandler.getLockout(req.SecretHash)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	req.WriteSuccessResponse(lockinInfo)
}

/*
用户随意选择一个公证人调用该接口,让公证人发起侧链PrepareLockin
注意:该接口不修改本地信息,在事件监听中才修改
*/
func (cs *CrossChainService) onSCPrepareLockinRequest(req *userapi.SCPrepareLockinRequest) {
	// 1. privateKey状态校验
	privateKeyInfo, err := cs.lockinHandler.db.LoadPrivateKeyInfo(cs.meta.SCTokenOwnerKey)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	if privateKeyInfo.Status != models.PrivateKeyNegotiateStatusFinished {
		panic("never happen")
	}
	// 2. lockinInfo状态校验
	lockinInfo, err := cs.lockinHandler.getLockin(req.SecretHash)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	if lockinInfo.MCUserAddressHex != req.MCUserAddress.String() {
		req.WriteErrorResponse(api.ErrorCodeException, "MCUserAddress wrong")
		return
	}
	if lockinInfo.MCLockStatus != models.LockStatusLock {
		req.WriteErrorResponse(api.ErrorCodeException, "MCLockStatus wrong")
		return
	}
	if lockinInfo.SCLockStatus != models.LockStatusNone {
		req.WriteErrorResponse(api.ErrorCodeException, "SCLockStatus wrong")
		return
	}
	// 3. 发起合约调用
	err = cs.callSCPrepareLockin(req, privateKeyInfo, lockinInfo)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, fmt.Sprintf("callSCPrepareLockin err = %s", err.Error()))
		return
	}
	// 4. 更新NotaryInCharge
	lockinInfo.NotaryIDInCharge = cs.selfNotaryID
	err = cs.lockinHandler.updateLockin(lockinInfo)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, fmt.Sprintf("lockinHandler.updateLockin err = %s", err.Error()))
		return
	}
	req.WriteSuccessResponse(lockinInfo)
}

/*
用户随意选择一个公证人调用该接口,让公证人发起主链PrepareLockout
注意:该接口不修改本地信息,在事件监听中才修改
*/
func (cs *CrossChainService) onMCPrepareLockoutRequest(req *userapi.MCPrepareLockoutRequest) {
	// 1. privateKey状态校验
	privateKeyInfo, err := cs.lockoutHandler.db.LoadPrivateKeyInfo(cs.meta.SCTokenOwnerKey)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	if privateKeyInfo.Status != models.PrivateKeyNegotiateStatusFinished {
		panic("never happen")
	}
	// 2. lockoutInfo状态校验
	lockoutInfo, err := cs.lockoutHandler.getLockout(req.SecretHash)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	if lockoutInfo.SCUserAddress != req.SCUserAddress {
		req.WriteErrorResponse(api.ErrorCodeException, "SCUserAddress wrong")
		return
	}
	if lockoutInfo.SCLockStatus != models.LockStatusLock {
		req.WriteErrorResponse(api.ErrorCodeException, "SCLockStatus wrong")
		return
	}
	if lockoutInfo.MCLockStatus != models.LockStatusNone {
		req.WriteErrorResponse(api.ErrorCodeException, "MCLockStatus wrong")
		return
	}
	if lockoutInfo.MCExpiration-cs.mcLastedBlockNumber <= params.MinLockoutMCExpiration {
		req.WriteErrorResponse(api.ErrorCodeException, "too late to MCPrepareLockout")
		return
	}
	// 3. 发起合约调用
	err = cs.callMCPrepareLockout(req, privateKeyInfo, lockoutInfo)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, fmt.Sprintf("callMCPrepareLockout err = %s", err.Error()))
		return
	}
	// 4. 更新NotaryInCharge
	lockoutInfo.NotaryIDInCharge = cs.selfNotaryID
	err = cs.lockoutHandler.updateLockout(lockoutInfo)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, fmt.Sprintf("lockoutHandler.updateLockout err = %s", err.Error()))
		return
	}
	req.WriteSuccessResponse(lockoutInfo)
}
