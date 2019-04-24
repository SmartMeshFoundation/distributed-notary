package bitcoin

import (
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
)

/* #nosec */
const (
	LockinEventName       = "LockinEvent"
	CancelLockinEventName = "CancelLockinEvent"
)

// NewBlockEvent :
type NewBlockEvent struct {
	*chain.BaseEvent
}

// createNewBlockEvent :
func createNewBlockEvent(blockNumber uint64) NewBlockEvent {
	e := NewBlockEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = ChainName
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = chain.NewBlockNumberEventName
	e.SCTokenAddress = utils.EmptyAddress
	return e
}

// LockinEvent :
type LockinEvent struct {
	*chain.BaseEvent
	SecretHash common.Hash   `json:"secret_hash"`
	TxHashStr  string        `json:"tx_hash_str"`
	TxOuts     []*wire.TxOut `json:"tx_outs"`
}

func createLockinEvent(blockNumber uint64, txHashStr string, secretHash common.Hash, txOuts []*wire.TxOut) LockinEvent {
	e := LockinEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = ChainName
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = LockinEventName
	e.SCTokenAddress = utils.EmptyAddress
	e.SecretHash = secretHash
	e.TxHashStr = txHashStr
	e.TxOuts = txOuts
	return e
}

// CancelLockinEvent :
type CancelLockinEvent struct {
	*chain.BaseEvent
	SecretHash common.Hash `json:"secret_hash"`
}

func createCancelLockinEvent(blockNumber uint64, secretHash common.Hash) CancelLockinEvent {
	e := CancelLockinEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = ChainName
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = CancelLockinEventName
	e.SCTokenAddress = utils.EmptyAddress
	e.SecretHash = secretHash
	return e
}
