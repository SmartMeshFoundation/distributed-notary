package service

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
)

/*
CrossChainService :
负责一个SCToken的所有相关事件及用户请求
*/
type CrossChainService struct {
	self         models.NotaryInfo
	db           *models.DB
	meta         *models.SideChainTokenMetaInfo
	scTokenProxy chain.ContractProxy
	mcProxy      chain.ContractProxy
}

// NewCrossChainService :
func NewCrossChainService(db *models.DB, self models.NotaryInfo, scTokenMetaInfo *models.SideChainTokenMetaInfo) *CrossChainService {
	// TODO init proxy
	return &CrossChainService{
		self: self,
		db:   db,
		meta: scTokenMetaInfo,
	}
}

// getMCContractAddress 获取主链合约地址
func (cs *CrossChainService) getMCContractAddress() common.Address {
	return cs.meta.MCLockedContractAddress
}

/*
	contract calls about lockin
*/

func (cs *CrossChainService) callSCPrepareLockin() (err error) {
	// TODO
	return
}

func (cs *CrossChainService) callMCLockin() (err error) {
	// TODO
	return
}

func (cs *CrossChainService) callSCCancelLockin() (err error) {
	// TODO
	return
}

/*
	contract calls about lockout
*/

func (cs *CrossChainService) callMCPrepareLockout() (err error) {
	// TODO
	return
}

func (cs *CrossChainService) callSCLockout() (err error) {
	// TODO
	return
}

func (cs *CrossChainService) callMCCancelLockout() (err error) {
	// TODO
	return
}
