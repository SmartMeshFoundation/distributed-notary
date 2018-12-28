package events

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestCreateLockinSecretEvent(t *testing.T) {
	var e LockinSecretEvent
	e.BaseEvent = createBaseEventFromSpectrumLog(EthereumTokenLockinSecretEventName, types.Log{
		Address:     utils.NewRandomAddress(),
		BlockNumber: 1,
	},
	)
	fmt.Println(utils.ToJSONStringFormat(e))
}
