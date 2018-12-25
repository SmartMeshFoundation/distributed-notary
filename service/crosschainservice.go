package service

import (
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
CrossChainService :
负责一个SCToken的所有相关事件及用户请求
*/
type CrossChainService struct {
	meta *models.SideChainTokenMetaInfo
}

// GetSCTokenAddress 获取侧链Token合约地址
func (ns *CrossChainService) GetSCTokenAddress() common.Address {
	// TODO
	return utils.EmptyAddress
}

// GetMCContractAddress 获取主链合约地址
func (ns *CrossChainService) GetMCContractAddress() common.Address {
	// TODO
	return utils.EmptyAddress
}
