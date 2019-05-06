package events

import (
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

// NewBlockEvent :
type NewBlockEvent struct {
	*chain.BaseEvent
}

// CreateNewBlockEvent :
func CreateNewBlockEvent(blockNumber uint64) NewBlockEvent {
	e := NewBlockEvent{}
	e.BaseEvent = &chain.BaseEvent{}
	e.ChainName = cfg.SMC.Name
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = chain.NewBlockNumberEventName
	e.SCTokenAddress = utils.EmptyAddress
	return e
}
