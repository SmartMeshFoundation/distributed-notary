package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/heco/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LockoutEvent :
type LockoutEvent struct {
	*chain.BaseEvent
	Account    common.Address `json:"account"` // lockout的用户地址
	SecretHash common.Hash    `json:"secret_hash"`
}

// CreateLockoutEvent :
func CreateLockoutEvent(log types.Log) (event LockoutEvent, err error) {
	e := &contracts.HecoTokenLockout{}
	err = unpackLog(&hecoTokenABI, e, HecoTokenLockoutEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromHecoLog(HecoTokenLockoutEventName, log)
	// params
	event.Account = e.Account
	event.SecretHash = e.SecretHash
	return
}
