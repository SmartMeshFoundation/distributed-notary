package messagetosign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"math/big"

	"context"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthereumCancelNonceTxDataName 用做消息传输时识别
const EthereumCancelNonceTxDataName = "EthereumPrepareLockoutTxData"

// EthereumCancelNonceTxData :
type EthereumCancelNonceTxData struct {
	BytesToSign []byte `json:"bytes_to_sign"`
	Nonce       uint64 `json:"nonce"`
}

// NewEthereumCancelNonceTxData :
func NewEthereumCancelNonceTxData(c chain.Chain, account common.Address, nonce uint64) (data *EthereumCancelNonceTxData, chainID *big.Int, rawTx *types.Transaction, err error) {
	data = &EthereumCancelNonceTxData{
		Nonce: nonce,
	}
	conn := c.GetConn()
	ctx := context.Background()
	amount := big.NewInt(1)
	msg := ethereum.CallMsg{From: account, To: &account, Value: amount, Data: nil}
	gasLimit, err := conn.EstimateGas(ctx, msg)
	if err != nil {
		return
	}
	gasPrice, err := conn.SuggestGasPrice(ctx)
	if err != nil {
		return
	}
	chainID, err = conn.NetworkID(ctx)
	if err != nil {
		return
	}
	rawTx = types.NewTransaction(nonce, account, amount, gasLimit, gasPrice, nil)
	//signer := func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	//	data.BytesToSign = signer.Hash(tx).Bytes()
	//	return nil, errShouldBe
	//}
	//_, err = signer(types.NewEIP155Signer(chainID), account, rawTx)
	//if err != errShouldBe {
	//	return
	//}
	data.BytesToSign = types.NewEIP155Signer(chainID).Hash(rawTx).Bytes()
	return
}

// GetSignBytes : impl MessageToSign
func (d *EthereumCancelNonceTxData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *EthereumCancelNonceTxData) GetName() string {
	return EthereumCancelNonceTxDataName
}

// GetTransportBytes : impl MessageToSign
func (d *EthereumCancelNonceTxData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *EthereumCancelNonceTxData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to EthereumCancelNonceTxData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData :
func (d *EthereumCancelNonceTxData) VerifySignData(c chain.Chain, account common.Address) (err error) {
	local, _, _, err := NewEthereumCancelNonceTxData(c, account, d.Nonce)
	if err != nil {
		return
	}
	if bytes.Compare(local.GetSignBytes(), d.GetSignBytes()) != 0 {
		err = fmt.Errorf("EthereumCancelNonceTxData.VerifySignBytes() fail,maybe attack")
	}
	return
}
