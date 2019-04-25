package service

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	ethevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	ethparams "github.com/ethereum/go-ethereum/params"
)

func (as *AdminService) onDebugTransferToAccountRequest(req *userapi.DebugTransferToAccountRequest) {
	amount := new(big.Int).Mul(big.NewInt(ethparams.Finney), big.NewInt(1000))
	namespace := []string{smcevents.ChainName, ethevents.ChainName}
	for _, name := range namespace {
		c, err := as.dispatchService.getChainByName(name)
		if err != nil {
			req.WriteErrorResponse(api.ErrorCodeException, err.Error())
			return
		}
		err = c.Transfer10ToAccount(as.dispatchService.getSelfPrivateKey(), req.Account, amount)
		if err != nil {
			req.WriteErrorResponse(api.ErrorCodeException, err.Error())
			return
		}
	}
	req.WriteSuccessResponse(nil)
}

func (as *AdminService) onDebugGetAllLockinInfo(req *userapi.DebugGetAllLockinInfoRequest) {
	lockinInfoList, err := as.db.GetAllLockinInfo()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	req.WriteSuccessResponse(lockinInfoList)
}

func (as *AdminService) onDebugGetAllLockoutInfo(req *userapi.DebugGetAllLockoutInfoRequest) {
	lockinInfoList, err := as.db.GetAllLockoutInfo()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	req.WriteSuccessResponse(lockinInfoList)
}

func (as *AdminService) onDebugGetAllBTCUtxo(req *userapi.DebugGetAllBTCUtxoRequest) {
	req.WriteSuccessResponse(as.db.GetBTCOutpointList(-1))
}
