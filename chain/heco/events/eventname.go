package events

import (
	"fmt"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/chain/heco/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// lockedEthereumABI :
var hecoTokenABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	hecoTokenABI, err = abi.JSON(strings.NewReader(contracts.HecoTokenABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[hecoTokenABI.Events[HecoTokenPrepareLockinEventName].Id()] = HecoTokenPrepareLockinEventName
	TopicToEventName[hecoTokenABI.Events[HecoTokenLockinSecretEventName].Id()] = HecoTokenLockinSecretEventName
	TopicToEventName[hecoTokenABI.Events[HecoTokenPrepareLockoutEventName].Id()] = HecoTokenPrepareLockoutEventName
	TopicToEventName[hecoTokenABI.Events[HecoTokenLockoutEventName].Id()] = HecoTokenLockoutEventName
	TopicToEventName[hecoTokenABI.Events[HecoTokenCancelLockinEventName].Id()] = HecoTokenCancelLockinEventName
	TopicToEventName[hecoTokenABI.Events[HecoTokenCancelLockoutEventName].Id()] = HecoTokenCancelLockoutEventName
}

/* #nosec */
const (
	HecoTokenPrepareLockinEventName  = "PrepareLockin"
	HecoTokenLockinSecretEventName   = "LockinSecret"
	HecoTokenPrepareLockoutEventName = "PrepareLockout"
	HecoTokenLockoutEventName        = "Lockout"
	HecoTokenCancelLockinEventName   = "CancelLockin"
	HecoTokenCancelLockoutEventName  = "CancelLockout"
)
