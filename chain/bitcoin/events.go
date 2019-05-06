package bitcoin

import (
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
)

/* #nosec */
const (
	LockinEventName         = "LockinEvent"
	CancelLockinEventName   = "CancelLockinEvent"
	PrepareLockoutEventName = "PrepareLockoutEvent"
	LockoutEventName        = "LockoutSecretEvent"
	CancelLockoutEventName  = "CancelLockoutEvent"
)

// NewBlockEvent :
type NewBlockEvent struct {
	*chain.BaseEvent
}

// createNewBlockEvent :
func createNewBlockEvent(blockNumber uint64) NewBlockEvent {
	e := NewBlockEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = cfg.BTC.Name
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
	e.ChainName = cfg.BTC.Name
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
	e.ChainName = cfg.BTC.Name
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
	SecretHash          common.Hash     `json:"secret_hash"`
	TXHashHex           string          `json:"tx_hash_hex"`
	LockOutpointIndex   int             `json:"lock_outpoint_index"`
	ChangeOutPointIndex int             `json:"change_out_point_index"`
	ChangeAmount        int64           `json:"change_amount"`
	UserAddressHex      string          `json:"user_address_hex"`
	MCExpiration        uint64          `json:"mc_expiration"`
	TxInOutpoint        []wire.OutPoint `json:"-"`
	TxOutLockScriptHex  string          `json:"tx_out_lock_script_hex"`
}

// createPrepareLockoutEvent :
func createPrepareLockoutEvent(blockNumber uint64, tx *wire.MsgTx, outpointRelevantInfo *BTCOutpointRelevantInfo) (event PrepareLockoutEvent) {
	e := PrepareLockoutEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = cfg.BTC.Name
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = PrepareLockoutEventName
	e.SCTokenAddress = utils.EmptyAddress
	e.TXHashHex = tx.TxHash().String()

	e.SecretHash = outpointRelevantInfo.SecretHash
	if len(tx.TxOut) == 1 {
		e.LockOutpointIndex = 1
		e.ChangeOutPointIndex = -1
	}
	if len(tx.TxOut) == 2 {
		// 包含找零vout
		e.LockOutpointIndex = 1
		e.ChangeOutPointIndex = 0
		e.ChangeAmount = tx.TxOut[0].Value
	}
	for _, txIn := range tx.TxIn {
		e.TxInOutpoint = append(e.TxInOutpoint, txIn.PreviousOutPoint)
	}
	e.UserAddressHex = outpointRelevantInfo.Data4PrepareLockout.UserAddressPublicKeyHashHex
	e.MCExpiration = outpointRelevantInfo.Data4PrepareLockout.MCExpiration
	e.TxOutLockScriptHex = outpointRelevantInfo.Data4PrepareLockout.TxOutLockScriptHex
	return e
}

// LockoutSecretEvent :
type LockoutSecretEvent struct {
	*chain.BaseEvent
	Secret common.Hash `json:"secre"`
}

func createLockoutEvent(blockNumber uint64, secret common.Hash) LockoutSecretEvent {
	e := LockoutSecretEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = cfg.BTC.Name
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = LockoutEventName
	e.SCTokenAddress = utils.EmptyAddress
	e.Secret = secret
	return e
}

// CancelLockoutEvent :
type CancelLockoutEvent struct {
	*chain.BaseEvent
	SecretHash common.Hash `json:"secret_hash"`
}

func createCancelLockoutEvent(blockNumber uint64, secretHash common.Hash) CancelLockoutEvent {
	e := CancelLockoutEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = cfg.BTC.Name
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = CancelLockoutEventName
	e.SCTokenAddress = utils.EmptyAddress
	e.SecretHash = secretHash
	return e
}
