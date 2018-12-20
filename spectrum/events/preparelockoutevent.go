package events

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PrepareLockoutEvent :
type PrepareLockoutEvent struct {
	BaseEvent
	TokenAddress common.Address // 事件发生的token合约地址
	Account      common.Address // 提出lockout的用户地址
	Amount       *big.Int       // 金额
}

// CreatePrepareLockoutEvent :
func CreatePrepareLockoutEvent(log types.Log) (event PrepareLockoutEvent, err error) {
	e := &contracts.EthereumTokenPrePareLockedOut{}
	err = UnpackLog(&EthereumTokenABI, e, EthereumTokenPrePareLockedOutEventName, &log)
	if err != nil {
		return
	}
	// BaseEvent
	event.Name = EthereumTokenPrePareLockedOutEventName
	event.BlockNumber = log.BlockNumber
	// params
	event.TokenAddress = log.Address
	event.Account = e.From
	event.Amount = e.Value
	return
}
