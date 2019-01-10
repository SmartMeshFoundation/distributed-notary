package events

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PrepareLockoutEvent :
type PrepareLockoutEvent struct {
	*chain.BaseEvent
	Account common.Address `json:"account"` // 提出lockout的用户地址
	Amount  *big.Int       `json:"amount"`  // 金额
}

// CreatePrepareLockoutEvent :
func CreatePrepareLockoutEvent(log types.Log) (event PrepareLockoutEvent, err error) {
	e := &contracts.LockedEthereumPrepareLockout{}
	err = unpackLog(&lockedEthereumABI, e, LockedEthereumPrepareLockoutEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromEthereumLog(LockedEthereumPrepareLockoutEventName, log)
	// params
	event.Account = e.Account
	event.Amount = e.Value
	return
}
