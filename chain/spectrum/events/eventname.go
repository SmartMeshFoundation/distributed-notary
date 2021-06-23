package events

import (
	"fmt"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// lockedSpectrumABI :
var lockedSpectrumABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	lockedSpectrumABI, err = abi.JSON(strings.NewReader(contracts.LockedSpectrumABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[lockedSpectrumABI.Events[LockedSpectrumPrepareLockinEventName].Id()] = LockedSpectrumPrepareLockinEventName
	TopicToEventName[lockedSpectrumABI.Events[LockedSpectrumLockoutSecretEventName].Id()] = LockedSpectrumLockoutSecretEventName
	TopicToEventName[lockedSpectrumABI.Events[LockedSpectrumPrepareLockoutEventName].Id()] = LockedSpectrumPrepareLockoutEventName
	TopicToEventName[lockedSpectrumABI.Events[LockedSpectrumLockinEventName].Id()] = LockedSpectrumLockinEventName
	TopicToEventName[lockedSpectrumABI.Events[LockedSpectrumCancelLockinEventName].Id()] = LockedSpectrumCancelLockinEventName
	TopicToEventName[lockedSpectrumABI.Events[LockedSpectrumCancelLockoutEventName].Id()] = LockedSpectrumCancelLockoutEventName
}

/* #nosec */
const (
	LockedSpectrumPrepareLockinEventName  = "PrepareLockin"
	LockedSpectrumLockoutSecretEventName  = "LockoutSecret"
	LockedSpectrumPrepareLockoutEventName = "PrepareLockout"
	LockedSpectrumLockinEventName         = "Lockin"
	LockedSpectrumCancelLockinEventName   = "CancelLockin"
	LockedSpectrumCancelLockoutEventName  = "CancelLockout"
)
