package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/ethereum/contracts"
	"github.com/ethereum/go-ethereum/core/types"
)

// CreateLockoutEvent :
func CreateIssueDocumentPOEvent(log types.Log) (event chainjettrade.IssueDocumentPOEvent, err error) {
	e := &contracts.Doc721IssueDocumentPO{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameIssueDocumentPO, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentPO, cfg.ETH.Name, log)
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
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameSignDocumentPO, cfg.ETH.Name, log)
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
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentDO, cfg.ETH.Name, log)
	// params
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}

// CreateLockoutEven t :
func CreateSignDocumentDOFFEvent(log types.Log) (event chainjettrade.SignDocumentDOFFEvent, err error) {
	e := &contracts.Doc721SignDocumentDOFF{}
	err = chainjettrade.UnpackLog(&docABI, e, EventNameSignDocumentDOFF, &log)
	if err != nil {
		return
	}
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameSignDocumentDOFF, cfg.ETH.Name, log)
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
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameSignDocumentDOBuyer, cfg.ETH.Name, log)
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
	event.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentINV, cfg.ETH.Name, log)
	// params
	// params
	event.From = e.From
	event.To = e.To
	event.TokenID = e.TokenId
	return
}
