package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/ethereum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LockinSecretEvent :
type LockinSecretEvent struct {
	BaseEvent
	TokenAddress common.Address // 事件发生的token合约地址
	Secret       common.Hash    // 用户提交到合约的密码
}

// CreateLockoutSecretEvent :
func CreateLockoutSecretEvent(log types.Log) (event LockinSecretEvent, err error) {
	e := &contracts.LockedEthereumLockoutSecret{}
	err = UnpackLog(&LockedEthereumABI, e, LockedEthereumLockoutSecretEventName, &log)
	if err != nil {
		return
	}
	createBaseEventFromLog(&event.BaseEvent, LockedEthereumLockoutSecretEventName, log)
	// params
	event.TokenAddress = log.Address
	event.Secret = e.Secret
	return
}
