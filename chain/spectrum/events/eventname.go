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

// ethereumTokenABI :
var ethereumTokenABI abi.ABI

// TopicToEventName :
var TopicToEventName map[common.Hash]string

func init() {
	var err error
	ethereumTokenABI, err = abi.JSON(strings.NewReader(contracts.EthereumTokenABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	TopicToEventName = make(map[common.Hash]string)
	TopicToEventName[ethereumTokenABI.Events[EthereumTokenPrepareLockinEventName].Id()] = EthereumTokenPrepareLockinEventName
	TopicToEventName[ethereumTokenABI.Events[EthereumTokenLockinSecretEventName].Id()] = EthereumTokenLockinSecretEventName
	TopicToEventName[ethereumTokenABI.Events[EthereumTokenPrePareLockedOutEventName].Id()] = EthereumTokenPrePareLockedOutEventName

}

/* #nosec */
const (
	EthereumTokenPrepareLockinEventName    = "EthereumTokenPrepareLockin"
	EthereumTokenLockinSecretEventName     = "EthereumTokenLockinSecret"
	EthereumTokenPrePareLockedOutEventName = "EthereumTokenPrePareLockedOut"
)
