package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
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
	e := &contracts.AtmosphereTokenLockout{}
	err = unpackLog(&atmosphereTokenABI, e, AtmosphereTokenLockoutEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromSpectrumLog(AtmosphereTokenLockoutEventName, log)
	// params
	event.Account = e.Account
	event.SecretHash = e.SecretHash
	return
}
