package events

import (
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/spectrum/contracts"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

// CreateLockoutEvent :
func CreateIssueDocumentPOEvent(log types.Log) (event chainjettrade.IssueDocumentPOEvent, err error) {
	e := &contracts.Doc721IssueDocumentPO{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameIssueDocumentPO, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentPO, log)
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}

// CreateLockoutEvent :
func CreateSignDocumentPOEvent(log types.Log) (event chainjettrade.SignDocumentPOEvent, err error) {
	e := &contracts.Doc721SignDocumentPO{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameSignDocumentPO, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameSignDocumentPO, log)
	// params
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}

// CreateLockoutEvent :
func CreateIssueDocumentDOEvent(log types.Log) (event chainjettrade.IssueDocumentDOEvent, err error) {
	e := &contracts.Doc721IssueDocumentDO{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameIssueDocumentDO, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentDO, log)
	// params
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}

// CreateLockoutEvent :
func CreateSignDocumentDOFFEvent(log types.Log) (event chainjettrade.SignDocumentDOFFEvent, err error) {
	e := &contracts.Doc721SignDocumentDOFF{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameSignDocumentDOFF, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameSignDocumentDOFF, log)
	// params
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}

// CreateLockoutEvent :
func CreateSignDocumentDOBuyerEvent(log types.Log) (event chainjettrade.SignDocumentDOBuyerEvent, err error) {
	e := &contracts.Doc721SignDocumentDOBuyer{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameSignDocumentDOBuyer, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameSignDocumentDOBuyer, log)
	// params
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}

// CreateLockoutEvent :
func CreateIssueDocumentINVEvent(log types.Log) (event chainjettrade.IssueDocumentINVEvent, err error) {
	e := &contracts.Doc721IssueDocumentINV{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameIssueDocumentINV, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentINV, log)
	// params
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}

// CreateNewBlockEvent :
func CreateNewBlockEvent(blockNumber uint64) chainjettrade.NewBlockEvent {
	e := chainjettrade.NewBlockEvent{}
	e.BaseEvent = &chainjettrade.BaseEvent{}
	e.ChainName = cfg.ETH.Name
	e.FromAddress = utils.EmptyAddress
	e.BlockNumber = blockNumber
	e.Time = time.Now()
	e.EventName = chain.NewBlockNumberEventName
	e.SCTokenAddress = utils.EmptyAddress
	return e
}
