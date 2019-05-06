package events

import (
	"fmt"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// lockedEthereumABI :
var lockedEthereumABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	lockedEthereumABI, err = abi.JSON(strings.NewReader(contracts.LockedEthereumABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[lockedEthereumABI.Events[LockedEthereumPrepareLockinEventName].Id()] = LockedEthereumPrepareLockinEventName
	TopicToEventName[lockedEthereumABI.Events[LockedEthereumLockoutSecretEventName].Id()] = LockedEthereumLockoutSecretEventName
	TopicToEventName[lockedEthereumABI.Events[LockedEthereumPrepareLockoutEventName].Id()] = LockedEthereumPrepareLockoutEventName
	TopicToEventName[lockedEthereumABI.Events[LockedEthereumLockinEventName].Id()] = LockedEthereumLockinEventName
	TopicToEventName[lockedEthereumABI.Events[LockedEthereumCancelLockinEventName].Id()] = LockedEthereumCancelLockinEventName
	TopicToEventName[lockedEthereumABI.Events[LockedEthereumCancelLockoutEventName].Id()] = LockedEthereumCancelLockoutEventName
}

/* #nosec */
const (
	LockedEthereumPrepareLockinEventName  = "PrepareLockin"
	LockedEthereumLockoutSecretEventName  = "LockoutSecret"
	LockedEthereumPrepareLockoutEventName = "PrepareLockout"
	LockedEthereumLockinEventName         = "Lockin"
	LockedEthereumCancelLockinEventName   = "CancelLockin"
	LockedEthereumCancelLockoutEventName  = "CancelLockout"
)
