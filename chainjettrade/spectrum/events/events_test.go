package events

import (
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/stretchr/testify/assert"

	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestCreateLockinSecretEvent(t *testing.T) {
	ast := assert.New(t)
	var e chainjettrade.IssueDocumentPOEvent
	e.BaseEvent = chainjettrade.CreateBaseEventFromLog(EventNameIssueDocumentPO, cfg.SMC.Name, types.Log{
		Address:     utils.NewRandomAddress(),
		BlockNumber: 1,
	},
	)
	ast.True(e.IsJettradeEvent())
	var e2 chain.Event = &e
	_, isJe := e2.(chainjettrade.IsJettradeEvent)
	ast.True(isJe)
	e2 = e
	_, isJe = e2.(chainjettrade.IsJettradeEvent)
	ast.True(isJe)
	fmt.Println(utils.ToJSONStringFormat(e))
}
