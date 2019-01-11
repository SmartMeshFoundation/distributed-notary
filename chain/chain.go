package chain

import (
	"crypto/ecdsa"
	"math/big"

	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// ErrorCallWrongChain : 调用错误
var ErrorCallWrongChain = errors.New("call wrong chain")

/*
Chain :
所有公链的统一接口
*/
type Chain interface {
	GetChainName() string
	GetEventChan() <-chan Event
	StartEventListener() error
	StopEventListener()
	RegisterEventListenContract(contractAddresses ...common.Address) error
	UnRegisterEventListenContract(contractAddresses ...common.Address)
	DeployContract(opts *bind.TransactOpts, params ...string) (contractAddress common.Address, err error)
	SetLastBlockNumber(lastBlockNumber uint64)
	GetContractProxy(contractAddress common.Address) ContractProxy

	Transfer10ToAccount(key *ecdsa.PrivateKey, accountTo common.Address, amount *big.Int) (err error) // for debug
}
