package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LockinSecretEvent :
type LockinSecretEvent struct {
	*chain.BaseEvent
	Secret common.Hash // 用户提交到合约的密码
}

// CreateLockoutSecretEvent :
func CreateLockoutSecretEvent(log types.Log) (event LockinSecretEvent, err error) {
	e := &contracts.LockedEthereumLockoutSecret{}
	err = UnpackLog(&lockedEthereumABI, e, LockedEthereumLockoutSecretEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromEthereumLog(LockedEthereumLockoutSecretEventName, log)
	// params
	event.Secret = e.Secret
	return
}