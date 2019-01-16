package chain

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
