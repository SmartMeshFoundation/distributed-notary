package events

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestCreateBaseEvent(t *testing.T) {
	var e LockoutSecretEvent
	e.BaseEvent = createBaseEventFromEthereumLog(LockedSpectrumLockoutSecretEventName, types.Log{
		BlockNumber: 1,
	})
	fmt.Println(utils.ToJSONStringFormat(e))
}
