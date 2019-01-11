package service

import (
	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

/*
CrossChainService :
负责一个SCToken的所有相关事件及用户请求
*/
type CrossChainService struct {
	selfPrivateKey *ecdsa.PrivateKey
	selfNotaryID   int
	meta           *models.SideChainTokenMetaInfo
	scTokenProxy   chain.ContractProxy
	mcProxy        chain.ContractProxy

	lockinHandler *lockinHandler
}

// NewCrossChainService :
func NewCrossChainService(db *models.DB, dispatchService dispatchServiceBackend, scTokenMetaInfo *models.SideChainTokenMetaInfo) *CrossChainService {
	scChain, err := dispatchService.getChainByName(smcevents.ChainName)
	if err != nil {
		panic("never happen")
	}
	mcChain, err := dispatchService.getChainByName(scTokenMetaInfo.MCName)
	if err != nil {
		panic("never happen")
	}
	return &CrossChainService{
		selfPrivateKey: dispatchService.getSelfPrivateKey(),
		selfNotaryID:   dispatchService.getSelfNotaryInfo().ID,
		meta:           scTokenMetaInfo,
		lockinHandler:  newLockinhandler(db),
		scTokenProxy:   scChain.GetContractProxy(scTokenMetaInfo.SCToken),
		mcProxy:        mcChain.GetContractProxy(scTokenMetaInfo.MCLockedContractAddress),
	}
}

// getMCContractAddress 获取主链合约地址
func (cs *CrossChainService) getMCContractAddress() common.Address {
	return cs.meta.MCLockedContractAddress
}

/*
	contract calls about lockin
*/

// SCPLI 需使用分布式签名
func (cs *CrossChainService) callSCPrepareLockin() (err error) {
	// TODO
	return
}

func (cs *CrossChainService) callMCLockin(userAddressHex string, secret common.Hash) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.mcProxy.Lockin(auth, userAddressHex, secret)
}

func (cs *CrossChainService) callSCCancelLockin(userAddressHex string) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.scTokenProxy.CancelLockin(auth, userAddressHex)
}

/*
	contract calls about lockout
*/

// MCPLO 需使用分布式签名
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
