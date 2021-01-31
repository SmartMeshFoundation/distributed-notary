package chain

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// EventName :
type EventName string

const (
	// NewBlockNumberEventName 新块事件,所有链公用
	NewBlockNumberEventName = "NewBlockNumber"
)

// Event :
type Event interface {
	GetSCTokenAddress() common.Address
	GetFromAddress() common.Address
	GetEventName() EventName
	GetChainName() string
	GetBlockNumber() uint64
}

/*
BaseEvent :
*/
type BaseEvent struct {
	ChainName   string         `json:"chain_name"`   // 事件所属链名
	FromAddress common.Address `json:"from_address"` // 产生该事件的合约地址
	BlockNumber uint64         `json:"block_number"` // 区块高度
	Time        time.Time      `json:"time"`         // 事件接收时间

	EventName      EventName      `json:"event_name"`
	SCTokenAddress common.Address `json:"sc_token_address"` // 该事件对应的侧链Token地址,主链事件该值为utils.EmptyHash
	TxHash         common.Hash
}

// GetSCTokenAddress :
func (be *BaseEvent) GetSCTokenAddress() common.Address {
	return be.SCTokenAddress
}

// GetFromAddress :
func (be *BaseEvent) GetFromAddress() common.Address {
	return be.FromAddress
}

// GetEventName :
func (be *BaseEvent) GetEventName() EventName {
	return be.EventName
}

// GetChainName :
func (be *BaseEvent) GetChainName() string {
	return be.ChainName
}

// GetBlockNumber :
func (be *BaseEvent) GetBlockNumber() uint64 {
	return be.BlockNumber
}
