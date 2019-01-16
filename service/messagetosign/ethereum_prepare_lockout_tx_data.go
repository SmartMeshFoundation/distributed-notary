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

// EthereumPrepareLockoutTxDataName 用做消息传输时识别
const EthereumPrepareLockoutTxDataName = "EthereumPrepareLockoutTxData"

// EthereumPrepareLockoutTxData : TODO
type EthereumPrepareLockoutTxData struct {
	BytesToSign  []byte                           `json:"bytes_to_sign"`
	UserRequest  *userapi.MCPrepareLockoutRequest `json:"user_request"`  // 用户原始请求,验证用户签名
	MCExpiration uint64                           `json:"mc_expiration"` // 主链超时块,由于公证人之间可能存在当前块误差,导致计算出来的主链超时块不一致,所以在协商时传递
}

// NewEthereumPrepareLockoutTxData :
func NewEthereumPrepareLockoutTxData(mcProxy chain.ContractProxy, req *userapi.MCPrepareLockoutRequest, callerAddress common.Address, mcUserAddressHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (data *EthereumPrepareLockoutTxData) {
	data = &EthereumPrepareLockoutTxData{
		UserRequest:  req,
		MCExpiration: expiration,
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
	err := mcProxy.PrepareLockout(transactor, mcUserAddressHex, secretHash, expiration, amount)
	if err != errShouldBe {
		// 这里不可能发生
		panic(err)
	}
	return
}

// GetSignBytes : impl MessageToSign
func (d *EthereumPrepareLockoutTxData) GetSignBytes() []byte {
	return d.BytesToSign
}

// GetName : impl MessageToSign
func (d *EthereumPrepareLockoutTxData) GetName() string {
	return EthereumPrepareLockoutTxDataName
}

// GetTransportBytes : impl MessageToSign
func (d *EthereumPrepareLockoutTxData) GetTransportBytes() []byte {
	buf, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return buf
}

// Parse : impl MessageToSign
func (d *EthereumPrepareLockoutTxData) Parse(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return errors.New("can not parse empty data to EthereumPrepareLockoutTxData")
	}
	return json.Unmarshal(buf, d)
}

// VerifySignData :
func (d *EthereumPrepareLockoutTxData) VerifySignData(scTokenProxy chain.ContractProxy, privateKeyInfo *models.PrivateKeyInfo, localLockoutInfo *models.LockoutInfo) (err error) {
	// 1. 校验本地lockinInfo状态
	if localLockoutInfo.SCUserAddress != d.UserRequest.SCUserAddress {
		err = fmt.Errorf("SCUserAddress wrong")
		return
	}
	if localLockoutInfo.SCLockStatus != models.LockStatusLock {
		err = fmt.Errorf("SCLockStatus wrong")
		return
	}
	if localLockoutInfo.MCLockStatus != models.LockStatusNone {
		err = fmt.Errorf("MCLockStatus wrong")
		return
	}
	if localLockoutInfo.MCExpiration != d.MCExpiration {
		log.Warn("localLockoutInfo.MCExpiration != request.MCExpiration, use request.MCExpiration")
		localLockoutInfo.SCExpiration = d.MCExpiration
	}
	// 2. 校验用户原始请求签名,验证请求中的SCUserAddress有效性
	if !d.UserRequest.VerifySign() {
		err = fmt.Errorf("signature in user request does't wrign")
		return
	}
	// 3. 使用本地数据获取MsgToSign
	mcUserAddressHex := d.UserRequest.MCUserAddress.String()
	mcExpiration := localLockoutInfo.MCExpiration
	secretHash := localLockoutInfo.SecretHash
	amount := localLockoutInfo.Amount
	var local *EthereumPrepareLockoutTxData
	local = NewEthereumPrepareLockoutTxData(scTokenProxy, d.UserRequest, privateKeyInfo.ToAddress(), mcUserAddressHex, secretHash, mcExpiration, amount)
	if bytes.Compare(local.GetSignBytes(), d.GetSignBytes()) != 0 {
		err = fmt.Errorf("EthereumPrepareLockoutTxData.VerifySignBytes() fail,maybe attack")
	}
	return
}
