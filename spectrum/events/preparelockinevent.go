package events

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PrepareLockinEvent :
type PrepareLockinEvent struct {
	BaseEvent
	TokenAddress common.Address // prepareLockin发生的token合约地址
	Account      common.Address // lockin的用户地址
	// htlc
	SecretHash common.Hash
	Expiration uint64
	Amount     *big.Int
}

// CreatePrepareLockinEvent :
func CreatePrepareLockinEvent(log types.Log) (event PrepareLockinEvent, err error) {
	e := &contracts.EthereumTokenPrepareLockin{}
	err = UnpackLog(&EthereumTokenABI, e, EthereumTokenPrepareLockinEventName, &log)
	if err != nil {
		return
	}
	// BaseEvent
	event.Name = EthereumTokenPrepareLockinEventName
	event.BlockNumber = log.BlockNumber
	// params
	event.TokenAddress = log.Address
	event.Account = e.Account
	event.SecretHash = e.SecretHash
	event.Expiration = e.Expiration.Uint64()
	event.Amount = e.Value
	return
}
