package events

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PrepareLockinEvent :
type PrepareLockinEvent struct {
	*chain.BaseEvent
	Account common.Address // lockin的用户地址
	// htlc
	Amount *big.Int
}

// CreatePrepareLockinEvent :
func CreatePrepareLockinEvent(log types.Log) (event PrepareLockinEvent, err error) {
	e := &contracts.EthereumTokenPrepareLockin{}
	err = unpackLog(&ethereumTokenABI, e, EthereumTokenPrepareLockinEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromSpectrumLog(EthereumTokenPrepareLockinEventName, log)
	// params
	event.Account = e.Account
	event.Amount = e.Value
	return
}
