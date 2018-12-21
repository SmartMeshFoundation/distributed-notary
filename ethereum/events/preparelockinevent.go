package events

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/ethereum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PrepareLockinEvent :
type PrepareLockinEvent struct {
	BaseEvent
	TokenAddress common.Address // prepareLockin发生的token合约地址
	Account      common.Address // lockin的用户地址
	// htlc
	Amount *big.Int
}

// CreatePrepareLockinEvent :
func CreatePrepareLockinEvent(log types.Log) (event PrepareLockinEvent, err error) {
	e := &contracts.LockedEthereumPrepareLockin{}
	err = UnpackLog(&LockedEthereumABI, e, LockedEthereumPrepareLockinEventName, &log)
	if err != nil {
		return
	}
	createBaseEventFromLog(&event.BaseEvent, LockedEthereumPrepareLockinEventName, log)
	// params
	event.TokenAddress = log.Address
	event.Account = e.Account
	event.Amount = e.Value
	return
}
