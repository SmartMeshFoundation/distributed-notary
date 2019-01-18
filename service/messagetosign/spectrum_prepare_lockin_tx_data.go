package messagetosign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

// SpectrumPrepareLockinTxDataName 用做消息传输时识别
const SpectrumPrepareLockinTxDataName = "SpectrumPrepareLockinTxData"

// SpectrumPrepareLockinTxData :
type SpectrumPrepareLockinTxData struct {
	BytesToSign  []byte                          `json:"bytes_to_sign"`
	UserRequest  *userapi.SCPrepareLockinRequest `json:"user_request"`  // 用户原始请求,验证用户签名
	SCExpiration uint64                          `json:"sc_expiration"` // 侧链超时块,由于公证人之间可能存在当前块误差,导致计算出来的侧链超时块不一致,所以在协商时传递
}

// NewSpectrumPrepareLockinTxData :
func NewSpectrumPrepareLockinTxData(scTokenProxy chain.ContractProxy, req *userapi.SCPrepareLockinRequest, callerAddress common.Address, scUserAddressHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (data *SpectrumPrepareLockinTxData) {
	data = &SpectrumPrepareLockinTxData{
		UserRequest:  req,
		SCExpiration: expiration,
	}
	transactor := &bind.TransactOpts{
		From: callerAddress,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != callerAddress {
				return nil, errors.New("not authorized to sign this account")
			}
			data.BytesToSign = signer.Hash(tx).Bytes()
			return nil, errShouldBe
		},
	}
	// 调用合约
	err := scTokenProxy.PrepareLockin(transactor, scUserAddressHex, secretHash, expiration, amount)
	if err != errShouldBe {
		// 这里不可能发生
		panic(err)
	}
	return
}

// GetSignBytes : impl MessageToSign
func (d *SpectrumPrepareLockinTxData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *SpectrumPrepareLockinTxData) GetName() string {
	return SpectrumPrepareLockinTxDataName
}

// GetTransportBytes : impl MessageToSign
func (d *SpectrumPrepareLockinTxData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *SpectrumPrepareLockinTxData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to SpectrumContractDeployTXData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData :
func (d *SpectrumPrepareLockinTxData) VerifySignData(scTokenProxy chain.ContractProxy, privateKeyInfo *models.PrivateKeyInfo, localLockinInfo *models.LockinInfo) (err error) {
	// 1. 校验本地lockinInfo状态
	if localLockinInfo.MCUserAddressHex != d.UserRequest.MCUserAddress.String() {
		err = fmt.Errorf("MCUserAddress wrong")
		return
	}
	if localLockinInfo.MCLockStatus != models.LockStatusLock {
		err = fmt.Errorf("MCLockStatus wrong")
		return
	}
	if localLockinInfo.SCLockStatus != models.LockStatusNone {
		err = fmt.Errorf("SCLockStatus wrong")
		return
	}
	if localLockinInfo.SCExpiration != d.SCExpiration {
		log.Warn("localLockinInfo.SCExpiration != request.SCExpiration, use request.SCExpiration")
		localLockinInfo.SCExpiration = d.SCExpiration
	}
	// 2. 校验用户原始请求签名,验证请求中的SCUserAddress有效性
	if !d.UserRequest.VerifySign() {
		err = fmt.Errorf("signature in user request does't wrign")
		return
	}
	// 3. 使用本地数据获取MsgToSign
	scUserAddressHex := d.UserRequest.SCUserAddress.String()
	scExpiration := localLockinInfo.SCExpiration
	secretHash := localLockinInfo.SecretHash
	amount := localLockinInfo.Amount
	var local *SpectrumPrepareLockinTxData
	local = NewSpectrumPrepareLockinTxData(scTokenProxy, d.UserRequest, privateKeyInfo.ToAddress(), scUserAddressHex, secretHash, scExpiration, amount)
	if bytes.Compare(local.GetSignBytes(), d.GetSignBytes()) != 0 {
		err = fmt.Errorf("SpectrumPrepareLockinTxData.VerifySignBytes() fail,maybe attack")
	}
	return
}
