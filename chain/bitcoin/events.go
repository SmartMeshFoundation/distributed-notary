package bitcoin

import (
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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

// PrepareLockoutEvent :
type PrepareLockoutEvent struct {
	*chain.BaseEvent
	SecretHash     common.Hash   `json:"secret_hash"`
	LockOutPoint   wire.OutPoint `json:"lock_out_point"` // 锁定到脚本的vout
	UserAddressHex string        `json:"user_address_hex"`
	MCExpiration   uint64        `json:"mc_expiration"`
}

// createPrepareLockoutEvent :
func createPrepareLockoutEvent(blockNumber uint64, txHash chainhash.Hash, outpointRelevantInfo *BTCOutpointRelevantInfo) (event PrepareLockoutEvent) {
	e := PrepareLockoutEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = ChainName
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = CancelLockinEventName
	e.SCTokenAddress = utils.EmptyAddress

	e.SecretHash = outpointRelevantInfo.SecretHash
	e.LockOutPoint = wire.OutPoint{
		Hash:  txHash,
		Index: 0,
	}
	e.UserAddressHex = outpointRelevantInfo.Data4PrepareLockout.UserAddressPublicKeyHashHex
	e.MCExpiration = outpointRelevantInfo.Data4PrepareLockout.MCExpiration
	return e
}
