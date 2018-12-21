package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LockinSecretEvent :
type LockinSecretEvent struct {
	BaseEvent
	TokenAddress common.Address // 事件发生的token合约地址
	Secret       common.Hash    // 用户提交到合约的密码
}

// CreateLockinSecretEvent :
func CreateLockinSecretEvent(log types.Log) (event LockinSecretEvent, err error) {
	e := &contracts.EthereumTokenLockinSecret{}
	err = UnpackLog(&EthereumTokenABI, e, EthereumTokenLockinSecretEventName, &log)
	if err != nil {
		return
	}
	createBaseEventFromLog(&event.BaseEvent, EthereumTokenLockinSecretEventName, log)
	// params
	event.TokenAddress = log.Address
	event.Secret = e.Secret
	return
}
