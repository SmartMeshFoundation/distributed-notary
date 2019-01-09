package events

import (
	"fmt"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ChainName 公链名
var ChainName = "spectrum"

// atmosphereTokenABI :
var atmosphereTokenABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	atmosphereTokenABI, err = abi.JSON(strings.NewReader(contracts.AtmosphereTokenABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[atmosphereTokenABI.Events[AtmosphereTokenPrepareLockinEventName].Id()] = AtmosphereTokenPrepareLockinEventName
	TopicToEventName[atmosphereTokenABI.Events[AtmosphereTokenLockinSecretEventName].Id()] = AtmosphereTokenLockinSecretEventName
	TopicToEventName[atmosphereTokenABI.Events[AtmosphereTokenPrepareLockoutEventName].Id()] = AtmosphereTokenPrepareLockoutEventName
	TopicToEventName[atmosphereTokenABI.Events[AtmosphereTokenLockoutEventName].Id()] = AtmosphereTokenLockoutEventName
	TopicToEventName[atmosphereTokenABI.Events[AtmosphereTokenCancelLockinEventName].Id()] = AtmosphereTokenCancelLockinEventName
	TopicToEventName[atmosphereTokenABI.Events[AtmosphereTokenCancelLockoutEventName].Id()] = AtmosphereTokenCancelLockoutEventName

}

/* #nosec */
const (
	AtmosphereTokenPrepareLockinEventName  = "AtmosphereTokenPrepareLockin"
	AtmosphereTokenLockinSecretEventName   = "AtmosphereTokenLockinSecret"
	AtmosphereTokenPrepareLockoutEventName = "AtmosphereTokenPrepareLockout"
	AtmosphereTokenLockoutEventName        = "AtmosphereTokenLockout"
	AtmosphereTokenCancelLockinEventName   = "AtmosphereTokenCancelLockin"
	AtmosphereTokenCancelLockoutEventName  = "AtmosphereTokenCancelLockout"
)
