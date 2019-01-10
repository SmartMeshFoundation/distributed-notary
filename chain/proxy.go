package chain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ContractProxy :
type ContractProxy interface {
	QueryLockin(account common.Address) (secretHash common.Hash, expiration uint64, amount *big.Int, data []byte, err error)
	QueryLockout(account common.Address) (secretHash common.Hash, expiration uint64, amount *big.Int, err error)
}
