package events

import (
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LockinSecretEvent :
type LockinSecretEvent struct {
	*chain.BaseEvent
	Secret common.Hash `json:"secret"` // 用户提交到合约的密码
}

// CreateLockinSecretEvent :
func CreateLockinSecretEvent(log types.Log) (event LockinSecretEvent, err error) {
	e := &contracts.AtmosphereTokenLockinSecret{}
	err = unpackLog(&atmosphereTokenABI, e, AtmosphereTokenLockinSecretEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromSpectrumLog(AtmosphereTokenLockinSecretEventName, log)
	// params
	event.Secret = e.Secret
	return
}
