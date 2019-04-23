package service

import (
	"errors"
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/Spectrum/log"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
)

// OnRequest restful请求处理
func (cs *CrossChainService) OnRequest(req api.Req) {

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
		r2, ok := req.(api.ReqWithResponse)
		if ok {
			r2.WriteErrorResponse(api.ErrorCodeParamsWrong)
			return
		}
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
	lockinInfo, err := cs.getLockInInfoBySCPrepareLockInRequest(req)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
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

func (cs *CrossChainService) getLockInInfoBySCPrepareLockInRequest(req *userapi.SCPrepareLockinRequest) (lockinInfo *models.LockinInfo, err error) {
	if cs.meta.MCName == events.ChainName {
		// 以太坊,收到用户请求时,lockinInfo必须已经存在,仅做校验工作
		mcUserAddress := common.BytesToAddress(req.MCUserAddress)
		lockinInfo, err = cs.lockinHandler.getLockin(req.SecretHash)
		if err != nil {
			req.WriteErrorResponse(api.ErrorCodeException, err.Error())
			return
		}
		if lockinInfo.MCUserAddressHex != mcUserAddress.String() {
			err = errors.New("MCUserAddress wrong")
			return
		}
		if lockinInfo.MCLockStatus != models.LockStatusLock {
			err = fmt.Errorf("MCLockStatus %d wrong", lockinInfo.MCLockStatus)
			return
		}
		if lockinInfo.SCLockStatus != models.LockStatusNone {
			err = fmt.Errorf("SCLockStatus %d wrong", lockinInfo.SCLockStatus)
			return
		}
		return
	}
	if cs.meta.MCName == bitcoin.ChainName {
		c, err2 := cs.dispatchService.getChainByName(cs.meta.MCName)
		if err2 != nil {
			panic(err2)
		}
		btcService := c.(*bitcoin.BTCService)
		net := cs.dispatchService.getBtcNetworkParam()
		mcUserAddress, err2 := btcutil.NewAddressPubKeyHash(req.MCUserAddress, net)
		if err2 != nil {
			log.Error(err2.Error())
			return nil, err2
		}
		notaryPublicKeyHash, err2 := btcutil.DecodeAddress(cs.meta.MCLockedPublicKeyHashStr, net)
		if err2 != nil {
			panic("never happen")
		}
		notaryAddress := notaryPublicKeyHash.(*btcutil.AddressPubKeyHash)
		// 比特币,收到用户请求之前,公证人无从得知任何信息,所以这里必须根据用户的请求验证链上数据,并构造lockinInfo
		// 1. 根据用户参数本地生成脚本hash
		builder := btcService.GetPrepareLockInScriptBuilder(mcUserAddress, notaryAddress, req.MCLockedAmount, req.SecretHash[:], req.MCExpiration)
		_, lockAddr, _ := builder.GetPKScript()
		// 2. 从链上查询tx并在其中查找锁定地址为上一步生成的hash值的outpoint及交易发生的blockNumber
		btcPrepareLockinInfo, err2 := btcService.GetPrepareLockinInfo(req.MCTXHash, lockAddr.String(), req.MCLockedAmount)
		if err2 != nil {
			return nil, err2
		}
		// 3. 校验mcExpiration
		mcExpiration := req.MCExpiration.Uint64()
		if mcExpiration-btcPrepareLockinInfo.BlockNumber <= params.MinLockinMCExpiration {
			err2 = fmt.Errorf("mcExpiration must bigger than %d", params.MinLockinMCExpiration)
			return nil, err2
		}
		// 4 计算scExpiration
		scExpiration := cs.scLastedBlockNumber + (mcExpiration - btcPrepareLockinInfo.BlockNumber - 5*params.ForkConfirmNumber - 1)
		// 5. 构造LockinInfo
		txHash := btcPrepareLockinInfo.TxHash
		lockinInfo = &models.LockinInfo{
			MCChainName:      cs.meta.MCName,
			SecretHash:       req.SecretHash,
			Secret:           utils.EmptyHash,
			MCUserAddressHex: mcUserAddress.String(),
			SCUserAddress:    req.SCUserAddress,
			SCTokenAddress:   cs.meta.SCToken,
			Amount:           big.NewInt(int64(req.MCLockedAmount)),
			MCExpiration:     mcExpiration,
			SCExpiration:     scExpiration,
			MCLockStatus:     models.LockStatusLock,
			SCLockStatus:     models.LockStatusNone,
			//Data:               data,
			NotaryIDInCharge:       models.UnknownNotaryIDInCharge,
			StartTime:              btcPrepareLockinInfo.BlockNumberTime,
			StartMCBlockNumber:     btcPrepareLockinInfo.BlockNumber,
			BTCPrepareLockinTXHash: &txHash,
			BTCPrepareLockinVout:   uint32(btcPrepareLockinInfo.Index),
		}
		// 6. 调用handler处理
		err2 = cs.lockinHandler.registerLockin(lockinInfo)
		if err2 != nil {
			err2 = fmt.Errorf("lockinHandler.registerLockin err2 = %s", err2.Error())
			return nil, err2
		}
		return
	}
	err = fmt.Errorf("unknown chain : %s", cs.meta.MCName)
	return
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
