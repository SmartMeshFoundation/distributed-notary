package messagetosign

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/api/userapi/jettradeapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/chainservice"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const JettradeTxIssuePODataName = "IssuePO"

//只是为了满足接口,并不真的做校验
type JettradeTxIssuePOData struct {
	UserRequest *jettradeapi.IssuePOOnSpectrumRequest
	Nonce       uint64
	BytesToSign []byte //必须是有意义的,因为需要其他公证人签署,所以非常关键
}

func NewJettradeTxIssuePOData(req *jettradeapi.IssuePOOnSpectrumRequest, nonce uint64, callerAddress common.Address, c *chainservice.SMTProxy) (data *JettradeTxIssuePOData, err error) {
	data = &JettradeTxIssuePOData{
		Nonce:       nonce,
		BytesToSign: []byte("padding"),
		UserRequest: req,
	}
	transactor := &bind.TransactOpts{
		From:  callerAddress,
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != callerAddress {
				return nil, errors.New("not authorized to sign this account")
			}
			data.BytesToSign = signer.Hash(tx).Bytes()
			return nil, errShouldBe
		},
	}
	// 调用合约
	_, err = c.CreateAndIssuePO2(transactor, req.TokenId, req.DocumentInfo, req.PONUm, req.Buyer, req.Farmer)
	if err != errShouldBe {
		// 这里不可能发生
		err = fmt.Errorf("CreateAndIssuePO2 err=%w", err)
		return
	}
	return data, nil
}

// GetSignBytes : impl MessageToSign
func (d *JettradeTxIssuePOData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *JettradeTxIssuePOData) GetName() string {
	return JettradeTxIssuePODataName
}

// GetTransportBytes : impl MessageToSign
func (d *JettradeTxIssuePOData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *JettradeTxIssuePOData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to EthereumPrepareLockoutTxData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData :
func (d *JettradeTxIssuePOData) VerifySignData(c chain.Chain, account common.Address) (err error) {
	//不做校验,跳过
	return nil
}
