package events

import (
	"fmt"
	"strings"

	"time"

	"github.com/SmartMeshFoundation/distributed-notary/ethereum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var LockedEthereumABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	LockedEthereumABI, err = abi.JSON(strings.NewReader(contracts.LockedEthereumABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	//event PrepareLockin(address indexed account,uint256 value);
	//event LockoutSecret(bytes32 secret);
	//event PrePareLockedOut(address indexed account, uint256 _value);
	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[LockedEthereumABI.Events[LockedEthereumPrepareLockinEventName].Id()] = LockedEthereumPrepareLockinEventName
	TopicToEventName[LockedEthereumABI.Events[LockedEthereumLockoutSecretEventName].Id()] = LockedEthereumLockoutSecretEventName
	TopicToEventName[LockedEthereumABI.Events[LockedEthereumPrePareLockedOutEventName].Id()] = LockedEthereumPrePareLockedOutEventName

}

// EventName :
type EventName string

const (
	// NewBlockEventName :
	NewBlockEventName = "NewBlockEvent"

	LockedEthereumPrepareLockinEventName    = "LockedEthereumPrepareLockin"
	LockedEthereumLockoutSecretEventName    = "LockedEthereumLockoutSecret"
	LockedEthereumPrePareLockedOutEventName = "LockedEthereumPrePareLockedOut"
)

// Event :
type Event interface{}

/*
BaseEvent :
*/
type BaseEvent struct {
	Name        EventName
	BlockNumber uint64
	Time        time.Time
}

func createBaseEventFromLog(e *BaseEvent, name EventName, log types.Log) {
	if e == nil {
		e = &BaseEvent{}
	}
	e.Name = name
	e.BlockNumber = log.BlockNumber
	e.Time = time.Now()
}
