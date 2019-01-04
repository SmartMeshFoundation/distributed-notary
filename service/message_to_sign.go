package service

import (
	"errors"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var errShouldBe = errors.New("should error")

// SpectrumContractDeployTX :
type SpectrumContractDeployTX struct {
	tx []byte
}

// NewSpectrumContractDeployTX :
func NewSpectrumContractDeployTX(c chain.Chain, callerAddress common.Address) (tx *SpectrumContractDeployTX) {
	var txBytes []byte
	transactor := &bind.TransactOpts{
		From: callerAddress,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != callerAddress {
				return nil, errors.New("not authorized to sign this account")
			}
			txBytes = signer.Hash(tx).Bytes()
			return nil, errShouldBe
		},
	}
	_, err := c.DeployContract(transactor)
	if err != errShouldBe {
		// 这里不可能发生
		panic(err)
	}
	return &SpectrumContractDeployTX{
		tx: txBytes,
	}
}

// GetHash : impl MessageToSign
func (s *SpectrumContractDeployTX) GetHash() common.Hash {
	return utils.Sha3(s.tx)
}

// GetBytes : impl MessageToSign
func (s *SpectrumContractDeployTX) GetBytes() []byte {
	return s.tx
}
