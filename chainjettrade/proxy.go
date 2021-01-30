package chainjettrade

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// ContractProxy :
type ContractProxy interface {
	QueryLockin(accountHex string) (secretHash common.Hash, expiration uint64, amount *big.Int, err error)
	QueryLockout(accountHex string) (secretHash common.Hash, expiration uint64, amount *big.Int, err error)

	PrepareLockin(opts *bind.TransactOpts, accountHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (err error)
	Lockin(opts *bind.TransactOpts, accountHex string, secret common.Hash) (err error)
	CancelLockin(opts *bind.TransactOpts, accountHex string) (err error)

	PrepareLockout(opts *bind.TransactOpts, accountHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (err error)
	Lockout(opts *bind.TransactOpts, accountHex string, secret common.Hash) (err error)
	CancelLockout(opts *bind.TransactOpts, accountHex string) (err error)
}

type ContractProxy2 interface {
	issueDO() (err error)
	singDOFF() (err error)
	signDOBuyer() (err error)
	//issuePO,在spectrum上时,需要做一些列工作,然后才能继续
	issuePO() (err error)
	issueINV() (err error)
}
