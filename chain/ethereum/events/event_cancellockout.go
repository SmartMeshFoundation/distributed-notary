package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// CancelLockoutEvent :
type CancelLockoutEvent struct {
	*chain.BaseEvent
	Account    common.Address `json:"account"` // lockout的用户地址
	SecretHash common.Hash    `json:"secret_hash"`
}

// CreateCancelLockoutEvent :
func CreateCancelLockoutEvent(log types.Log) (event CancelLockoutEvent, err error) {
	e := &contracts.LockedEthereumCancelLockout{}
	err = unpackLog(&lockedEthereumABI, e, LockedEthereumCancelLockoutEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromEthereumLog(LockedEthereumCancelLockoutEventName, log)
	// params
	event.Account = e.Account
	event.SecretHash = e.SecretHash
	return
}
