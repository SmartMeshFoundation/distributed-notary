package service

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
NotaryService :
负责一组公证人之间的消息通讯及业务处理,需保证线程安全
*/
type NotaryService struct {
	meta *models.SideChainTokenMetaInfo
}

// OnChainEvent 链上事件逻辑处理
func (ns *NotaryService) OnChainEvent(e chain.Event) error {
	//TODO
	return nil
}

// GetMCContractAddress 获取主链合约地址
func (ns *NotaryService) GetMCContractAddress() common.Address {
	// TODO
	return utils.EmptyAddress
}
