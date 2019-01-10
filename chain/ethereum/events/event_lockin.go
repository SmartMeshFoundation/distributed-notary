package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
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
	e := &contracts.LockedEthereumLockin{}
	err = unpackLog(&lockedEthereumABI, e, LockedEthereumLockinEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromEthereumLog(LockedEthereumLockinEventName, log)
	// params
	event.Account = e.Account
	event.SecretHash = e.SecretHash
	return
}
