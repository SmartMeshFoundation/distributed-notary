package events

import (
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestCreateLockinSecretEvent(t *testing.T) {
	var e chainjettrade.IssueDocumentPOEvent
	e.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentPO, types.Log{
		Address:     utils.NewRandomAddress(),
		BlockNumber: 1,
	},
	)
	fmt.Println(utils.ToJSONStringFormat(e))
}