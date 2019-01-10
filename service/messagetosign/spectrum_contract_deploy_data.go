package messagetosign

import (
	"errors"

	"encoding/json"

	"bytes"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var errShouldBe = errors.New("should error")

// SpectrumContractDeployTXDataName 用做消息传输时识别
const SpectrumContractDeployTXDataName = "SpectrumContractDeployTXData"

// SpectrumContractDeployTXData :
type SpectrumContractDeployTXData struct {
	BytesToSign     []byte `json:"bytes_to_sign"`
	DeployChainName string `json:"deploy_chain_name"`
	TokenName       string `json:"token_name"` // 如果为侧链token,需要token名
}

// NewSpectrumContractDeployTX :
func NewSpectrumContractDeployTX(c chain.Chain, callerAddress common.Address, params ...string) (tx *SpectrumContractDeployTXData) {
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
	_, err := c.DeployContract(transactor, params...)
	if err != errShouldBe {
		// 这里不可能发生
		panic(err)
	}
	tx = &SpectrumContractDeployTXData{
		BytesToSign:     txBytes,
		DeployChainName: c.GetChainName(),
	}
	if c.GetChainName() == events.ChainName {
		tx.TokenName = params[0]
	}
	return
}

// GetSignBytes : impl MessageToSign
func (s *SpectrumContractDeployTXData) GetSignBytes() []byte {
	return s.BytesToSign
}

// GetName : impl MessageToSign
func (s *SpectrumContractDeployTXData) GetName() string {
	return SpectrumContractDeployTXDataName
}

// GetTransportBytes : impl MessageToSign
func (s *SpectrumContractDeployTXData) GetTransportBytes() []byte {
	buf, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (s *SpectrumContractDeployTXData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to SpectrumContractDeployTXData")
	}
	return json.Unmarshal(buf, s)
}

// VerifySignBytes :
func (s *SpectrumContractDeployTXData) VerifySignBytes(c chain.Chain, callerAddress common.Address) (err error) {
	var local *SpectrumContractDeployTXData
	if s.DeployChainName == events.ChainName {
		local = NewSpectrumContractDeployTX(c, callerAddress, s.TokenName)
	} else {
		local = NewSpectrumContractDeployTX(c, callerAddress)
	}
	if bytes.Compare(local.GetSignBytes(), s.GetSignBytes()) != 0 {
		err = fmt.Errorf("SpectrumContractDeployTXData.VerifySignBytes() fail,maybe attack")
	}
	return
}
