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
	Secret common.Hash // 用户提交到合约的密码
}

// CreateLockinSecretEvent :
func CreateLockinSecretEvent(log types.Log) (event LockinSecretEvent, err error) {
	e := &contracts.EthereumTokenLockinSecret{}
	err = unpackLog(&ethereumTokenABI, e, EthereumTokenLockinSecretEventName, &log)
	if err != nil {
		return
	}
	event.BaseEvent = createBaseEventFromSpectrumLog(EthereumTokenLockinSecretEventName, log)
	// params
	event.Secret = e.Secret
	return
}