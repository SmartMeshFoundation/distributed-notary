package service

import (
	"errors"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var errShouldBe = errors.New("should error")

// SpectrumContractDeployTXDataNameName 用做消息传输时识别
const SpectrumContractDeployTXDataNameName = "SpectrumContractDeployTXData"

// SpectrumContractDeployTXData :
type SpectrumContractDeployTXData struct {
	tx []byte
}

// NewSpectrumContractDeployTX :
func NewSpectrumContractDeployTX(c chain.Chain, callerAddress common.Address) (tx *SpectrumContractDeployTXData) {
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
	return &SpectrumContractDeployTXData{
		tx: txBytes,
	}
}

// GetBytes : impl MessageToSign
func (s *SpectrumContractDeployTXData) GetBytes() []byte {
	return s.tx
}

// GetName : impl MessageToSign
func (s *SpectrumContractDeployTXData) GetName() string {
	return SpectrumContractDeployTXDataNameName
}

// Parse : impl MessageToSign
func (s *SpectrumContractDeployTXData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to SpectrumContractDeployTXData")
	}
	s.tx = buf
	return nil
}
