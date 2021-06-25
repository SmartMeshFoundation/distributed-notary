package messagetosign

/*
	以太坊合约部署及specrum合约部署均使用该消息体
*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var errShouldBe = errors.New("should error")

// HecoContractDeployTXDataName 用做消息传输时识别
const HecoContractDeployTXDataName = "HecoContractDeployTXData"

// HecoContractDeployTXData :
type HecoContractDeployTXData struct {
	BytesToSign     []byte `json:"bytes_to_sign"`
	Nonce           uint64 `json:"nonce"`
	DeployChainName string `json:"deploy_chain_name"`
	TokenName       string `json:"token_name"` // 如果为侧链token,需要token名
}

// NewHecoContractDeployTX :
func NewHecoContractDeployTX(c chain.Chain, callerAddress common.Address, nonce uint64, params ...string) (tx *HecoContractDeployTXData) {
	var txBytes []byte
	transactor := &bind.TransactOpts{
		From:  callerAddress,
		Nonce: big.NewInt(int64(nonce)),
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
	tx = &HecoContractDeployTXData{
		Nonce:           nonce,
		BytesToSign:     txBytes,
		DeployChainName: c.GetChainName(),
	}
	if c.GetChainName() == cfg.SMC.Name {
		tx.TokenName = params[0]
	}
	return
}

// GetSignBytes : impl MessageToSign
func (s *HecoContractDeployTXData) GetSignBytes() []byte {
	return s.BytesToSign
}

// GetName : impl MessageToSign
func (s *HecoContractDeployTXData) GetName() string {
	return HecoContractDeployTXDataName
}

// GetTransportBytes : impl MessageToSign
func (s *HecoContractDeployTXData) GetTransportBytes() []byte {
	buf, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (s *HecoContractDeployTXData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to HecoContractDeployTXData")
	}
	return json.Unmarshal(buf, s)
}

// VerifySignBytes :
func (s *HecoContractDeployTXData) VerifySignBytes(c chain.Chain, callerAddress common.Address) (err error) {
	var local *HecoContractDeployTXData
	if s.DeployChainName == cfg.SMC.Name {
		local = NewHecoContractDeployTX(c, callerAddress, s.Nonce, s.TokenName)
	} else {
		local = NewHecoContractDeployTX(c, callerAddress, s.Nonce)
	}
	if bytes.Compare(local.GetSignBytes(), s.GetSignBytes()) != 0 {
		err = fmt.Errorf("HecoContractDeployTXData.VerifySignBytes() fail,maybe attack")
	}
	return
}
