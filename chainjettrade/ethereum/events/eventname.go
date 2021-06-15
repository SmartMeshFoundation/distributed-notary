package events

import (
	"fmt"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/ethereum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// docABI :
var docABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	docABI, err = abi.JSON(strings.NewReader(contracts.Doc721ABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err  %s", err))
	}
	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[docABI.Events[EventNameIssueDocumentPO].Id()] = EventNameIssueDocumentPO
	TopicToEventName[docABI.Events[EventNameSignDocumentPO].Id()] = EventNameSignDocumentPO
	TopicToEventName[docABI.Events[EventNameIssueDocumentDO].Id()] = EventNameIssueDocumentDO
	TopicToEventName[docABI.Events[EventNameSignDocumentDOFF].Id()] = EventNameSignDocumentDOFF
	TopicToEventName[docABI.Events[EventNameSignDocumentDOBuyer].Id()] = EventNameSignDocumentDOBuyer
	TopicToEventName[docABI.Events[EventNameIssueDocumentINV].Id()] = EventNameIssueDocumentINV

}

/* #nosec */
const (
	EventNameIssueDocumentPO     = "IssueDocument_PO"
	EventNameSignDocumentPO      = "SignDocument_PO"
	EventNameIssueDocumentDO     = "IssueDocument_DO"
	EventNameSignDocumentDOFF    = "SignDocument_DO_FF"
	EventNameSignDocumentDOBuyer = "SignDocument_DO_Buyer"
	EventNameIssueDocumentINV    = "IssueDocument_INV"
)
