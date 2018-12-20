package events

import (
	"fmt"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/spectrum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var EthereumTokenABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	EthereumTokenABI, err = abi.JSON(strings.NewReader(contracts.EthereumTokenABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}

	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[EthereumTokenABI.Events[EthereumTokenPrepareLockinEventName].Id()] = EthereumTokenPrepareLockinEventName
	TopicToEventName[EthereumTokenABI.Events[EthereumTokenLockinSecretEventName].Id()] = EthereumTokenLockinSecretEventName
	TopicToEventName[EthereumTokenABI.Events[EthereumTokenPrePareLockedOutEventName].Id()] = EthereumTokenPrePareLockedOutEventName

}

// EventName :
type EventName string

const (
	// NewBlockEventName :
	NewBlockEventName = "NewBlockEvent"

	EthereumTokenPrepareLockinEventName    = "EthereumTokenPrepareLockin"
	EthereumTokenLockinSecretEventName     = "EthereumTokenLockinSecret"
	EthereumTokenPrePareLockedOutEventName = "EthereumTokenPrePareLockedOut"
)

// Event :
type Event interface{}

/*
BaseEvent :
*/
type BaseEvent struct {
	Name        EventName
	BlockNumber uint64
}
