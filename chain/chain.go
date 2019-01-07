package chain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

/*
Chain :

*/
type Chain interface {
	GetChainName() string
	GetEventChan() <-chan Event
	StartEventListener() error
	StopEventListener()
	RegisterEventListenContract(contractAddresses ...common.Address) error
	UnRegisterEventListenContract(contractAddresses ...common.Address)
	DeployContract(opts *bind.TransactOpts) (contractAddress common.Address, err error)
}
