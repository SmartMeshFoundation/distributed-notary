package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LockoutSecretEvent :
type LockoutSecretEvent struct {
	*chain.BaseEvent
	Secret common.Hash `json:"secret"` // 用户提交到合约的密码
}

// CreateLockoutSecretEvent :
func CreateLockoutSecretEvent(log types.Log) (event LockoutSecretEvent, err error) {
	e := &contracts.LockedSpectrumLockoutSecret{}
	err = unpackLog(&lockedSpectrumABI, e, LockedSpectrumLockoutSecretEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromEthereumLog(LockedSpectrumLockoutSecretEventName, log)
	// params
	event.Secret = e.Secret
	return
}
