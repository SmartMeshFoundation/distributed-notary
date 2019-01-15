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
	amount := new(big.Int).Mul(big.NewInt(ethparams.Ether), big.NewInt(10))
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
