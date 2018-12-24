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
	Account common.Address // 提出lockout的用户地址
	Amount  *big.Int       // 金额
}

// CreatePrepareLockoutEvent :
func CreatePrepareLockoutEvent(log types.Log) (event PrepareLockoutEvent, err error) {
	e := &contracts.EthereumTokenPrePareLockedOut{}
	err = unpackLog(&ethereumTokenABI, e, EthereumTokenPrePareLockedOutEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromSpectrumLog(EthereumTokenPrePareLockedOutEventName, log)
	// params
	event.Account = e.Account
	event.Amount = e.Value
	return
}
