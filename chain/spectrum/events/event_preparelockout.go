package events

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
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
	e := &contracts.AtmosphereTokenPrepareLockout{}
	err = unpackLog(&atmosphereTokenABI, e, AtmosphereTokenPrepareLockoutEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromSpectrumLog(AtmosphereTokenPrepareLockoutEventName, log)
	// params
	event.Account = e.Account
	event.Amount = e.Value
	return
}
