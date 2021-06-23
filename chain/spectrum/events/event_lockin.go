package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LockinEvent :
type LockinEvent struct {
	*chain.BaseEvent
	Account    common.Address `json:"account"` // lockout的用户地址
	SecretHash common.Hash    `json:"secret_hash"`
}

// CreateLockinEvent :
func CreateLockinEvent(log types.Log) (event LockinEvent, err error) {
	e := &contracts.LockedSpectrumLockin{}
	err = unpackLog(&lockedSpectrumABI, e, LockedSpectrumLockinEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromEthereumLog(LockedSpectrumLockinEventName, log)
	// params
	event.Account = e.Account
	event.SecretHash = e.SecretHash
	return
}
