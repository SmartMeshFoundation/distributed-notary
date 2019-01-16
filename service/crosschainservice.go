package service

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/kataras/iris/core/errors"
	"github.com/nkbai/log"
)

/*
CrossChainService :
负责一个SCToken的所有相关事件及用户请求
*/
type CrossChainService struct {
	selfPrivateKey *ecdsa.PrivateKey
	selfNotaryID   int
	meta           *models.SideChainTokenMetaInfo
	scTokenProxy   chain.ContractProxy
	mcProxy        chain.ContractProxy

	scLastedBlockNumber uint64
	mcLastedBlockNumber uint64

	dispatchService dispatchServiceBackend

	lockinHandler  *lockinHandler
	lockoutHandler *lockoutHandler
}

// NewCrossChainService :
func NewCrossChainService(db *models.DB, dispatchService dispatchServiceBackend, scTokenMetaInfo *models.SideChainTokenMetaInfo) *CrossChainService {
	scChain, err := dispatchService.getChainByName(smcevents.ChainName)
	if err != nil {
		panic("never happen")
	}
	mcChain, err := dispatchService.getChainByName(scTokenMetaInfo.MCName)
	if err != nil {
		panic("never happen")
	}
	return &CrossChainService{
		selfPrivateKey:  dispatchService.getSelfPrivateKey(),
		selfNotaryID:    dispatchService.getSelfNotaryInfo().ID,
		meta:            scTokenMetaInfo,
		dispatchService: dispatchService,
		lockinHandler:   newLockinHandler(db, scTokenMetaInfo.SCToken),
		lockoutHandler:  newLockoutHandler(db, scTokenMetaInfo.SCToken),
		scTokenProxy:    scChain.GetContractProxy(scTokenMetaInfo.SCToken),
		mcProxy:         mcChain.GetContractProxy(scTokenMetaInfo.MCLockedContractAddress),
	}
}

// getMCContractAddress 获取主链合约地址
func (cs *CrossChainService) getMCContractAddress() common.Address {
	return cs.meta.MCLockedContractAddress
}

/*
	contract calls about lockin
*/

// SCPLI 需使用分布式签名
func (cs *CrossChainService) callSCPrepareLockin(req *userapi.SCPrepareLockinRequest, privateKeyInfo *models.PrivateKeyInfo, localLockinInfo *models.LockinInfo) (err error) {
	// 从本地获取调用合约的参数
	scUserAddressHex := req.SCUserAddress.String()
	scExpiration := localLockinInfo.SCExpiration
	secretHash := localLockinInfo.SecretHash
	amount := localLockinInfo.Amount
	// 1. 构造MessageToSign
	var msgToSign mecdsa.MessageToSign
	msgToSign = messagetosign.NewSpectrumPrepareLockinTxData(cs.scTokenProxy, req, privateKeyInfo.ToAddress(), scUserAddressHex, secretHash, scExpiration, amount)
	// 2. 发起分布式签名
	var signature []byte
	var _ common.Hash
	signature, _, err = cs.dispatchService.getNotaryService().startDistributedSignAndWait(msgToSign, privateKeyInfo)
	if err != nil {
		return
	}
	log.Info("call PrepareLockin on spectrum with account=%s, signature=%s", privateKeyInfo.ToAddress().String(), common.Bytes2Hex(signature))
	// 3. 调用合约
	transactor := &bind.TransactOpts{
		From: privateKeyInfo.ToAddress(),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != privateKeyInfo.ToAddress() {
				return nil, errors.New("not authorized to sign this account")
			}
			msgToSign2 := signer.Hash(tx).Bytes()
			if bytes.Compare(msgToSign.GetSignBytes(), msgToSign2) != 0 {
				err = fmt.Errorf("txbytes when deploy contract step1 and step2 does't match")
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
	return cs.scTokenProxy.PrepareLockin(transactor, scUserAddressHex, secretHash, scExpiration, amount)
}

func (cs *CrossChainService) callMCLockin(userAddressHex string, secret common.Hash) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.mcProxy.Lockin(auth, userAddressHex, secret)
}

func (cs *CrossChainService) callSCCancelLockin(userAddressHex string) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.scTokenProxy.CancelLockin(auth, userAddressHex)
}

/*
	contract calls about lockout
*/

// MCPLO 需使用分布式签名
func (cs *CrossChainService) callMCPrepareLockout(req *userapi.MCPrepareLockoutRequest, privateKeyInfo *models.PrivateKeyInfo, localLockoutInfo *models.LockoutInfo) (err error) {
	// 从本地获取调用合约的参数
	mcUserAddressHex := req.MCUserAddress.String()
	mcExpiration := localLockoutInfo.MCExpiration
	secretHash := localLockoutInfo.SecretHash
	amount := localLockoutInfo.Amount
	// 1. 构造MessageToSign
	var msgToSign mecdsa.MessageToSign
	msgToSign = messagetosign.NewEthereumPrepareLockoutTxData(cs.scTokenProxy, req, privateKeyInfo.ToAddress(), mcUserAddressHex, secretHash, mcExpiration, amount)
	// 2. 发起分布式签名
	var signature []byte
	var _ common.Hash
	signature, _, err = cs.dispatchService.getNotaryService().startDistributedSignAndWait(msgToSign, privateKeyInfo)
	if err != nil {
		return
	}
	log.Info("call PrepareLockout on ethereum with account=%s, signature=%s", privateKeyInfo.ToAddress().String(), common.Bytes2Hex(signature))
	// 3. 调用合约
	transactor := &bind.TransactOpts{
		From: privateKeyInfo.ToAddress(),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != privateKeyInfo.ToAddress() {
				return nil, errors.New("not authorized to sign this account")
			}
			msgToSign2 := signer.Hash(tx).Bytes()
			if bytes.Compare(msgToSign.GetSignBytes(), msgToSign2) != 0 {
				err = fmt.Errorf("txbytes when deploy contract step1 and step2 does't match")
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
	return cs.mcProxy.PrepareLockout(transactor, mcUserAddressHex, secretHash, mcExpiration, amount)
}

func (cs *CrossChainService) callSCLockout(userAddressHex string, secret common.Hash) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.scTokenProxy.Lockout(auth, userAddressHex, secret)
}

func (cs *CrossChainService) callMCCancelLockout(userAddressHex string) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.scTokenProxy.CancelLockout(auth, userAddressHex)
}
