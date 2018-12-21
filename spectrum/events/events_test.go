package events

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestCreateBaseEvent(t *testing.T) {
	var e LockinSecretEvent
	createBaseEventFromLog(&e.BaseEvent, EthereumTokenLockinSecretEventName, types.Log{
		BlockNumber: 1,
	})
	fmt.Println(utils.ToJsonStringFormat(e))
}
