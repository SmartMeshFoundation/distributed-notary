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

// SpectrumCancelNonceTxDataName 用做消息传输时识别
const SpectrumCancelNonceTxDataName = "SpectrumCancelNonceTxData"

// SpectrumCancelNonceTxData :
type SpectrumCancelNonceTxData struct {
	BytesToSign []byte `json:"bytes_to_sign"`
	ChainName   string `json:"chain_name"`
	Account     string `json:"account"`
	Nonce       uint64 `json:"nonce"`
}

// NewEthereumCancelNonceTxData :
func NewSpectrumCancelNonceTxData(c chain.Chain, account common.Address, nonce uint64) (data *SpectrumCancelNonceTxData, chainID *big.Int, rawTx *types.Transaction, err error) {
	data = &SpectrumCancelNonceTxData{
		ChainName: c.GetChainName(),
		Account:   account.String(),
		Nonce:     nonce,
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
func (d *SpectrumCancelNonceTxData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *SpectrumCancelNonceTxData) GetName() string {
	return SpectrumCancelNonceTxDataName
}

// GetTransportBytes : impl MessageToSign
func (d *SpectrumCancelNonceTxData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *SpectrumCancelNonceTxData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to SpectrumCancelNonceTxData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData :
func (d *SpectrumCancelNonceTxData) VerifySignData(c chain.Chain, account common.Address) (err error) {
	local, _, _, err := NewSpectrumCancelNonceTxData(c, account, d.Nonce)
	if err != nil {
		return
	}
	if bytes.Compare(local.GetSignBytes(), d.GetSignBytes()) != 0 {
		err = fmt.Errorf("SpectrumCancelNonceTxData.VerifySignBytes() fail,maybe attack")
	}
	return
}
