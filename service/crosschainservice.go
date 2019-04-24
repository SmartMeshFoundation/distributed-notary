package service

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"errors"

	"math/big"

	"time"

	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	scUserAddressHex := req.GetSignerETHAddress().String()
	scExpiration := localLockinInfo.SCExpiration
	secretHash := localLockinInfo.SecretHash
	amount := localLockinInfo.Amount
	// 0. 获取nonce
	nonce, err := cs.dispatchService.applyNonceFromNonceServer(smcevents.ChainName, privateKeyInfo.Key, req.SecretHash.String(), amount)
	if err != nil {
		return
	}
	// 1. 构造MessageToSign
	var msgToSign messagetosign.MessageToSign
	msgToSign = messagetosign.NewSpectrumPrepareLockinTxData(cs.scTokenProxy, req, privateKeyInfo.ToAddress(), scUserAddressHex, secretHash, scExpiration, amount, nonce)
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
		From:  privateKeyInfo.ToAddress(),
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != privateKeyInfo.ToAddress() {
				return nil, errors.New("not authorized to sign this account")
			}
			msgToSign2 := signer.Hash(tx).Bytes()
			fmt.Printf("======================from=%s nonce=%d\n", address.String(), tx.Nonce())
			if bytes.Compare(msgToSign.GetSignBytes(), msgToSign2) != 0 {
				err = fmt.Errorf("txbytes when deploy contract step1 and step2 does't match")
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
	return cs.scTokenProxy.PrepareLockin(transactor, scUserAddressHex, secretHash, scExpiration, amount)
}

func (cs *CrossChainService) callMCLockin(lockinInfo *models.LockinInfo) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	if cs.meta.MCName == events.ChainName {
		// 以太坊直接使用自己私钥调用合约即可
		auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
		return cs.mcProxy.Lockin(auth, lockinInfo.MCUserAddressHex, lockinInfo.Secret)
	} else if cs.meta.MCName == bitcoin.ChainName {
		// 比特币的lockin操作需要分布式签名
		return cs.callBTCLockin(lockinInfo)
	}
	return errors.New("unknown chain")
}

func (cs *CrossChainService) callBTCLockin(lockinInfo *models.LockinInfo) (err error) {
	// 0. 获取bs
	c, err := cs.dispatchService.getChainByName(bitcoin.ChainName)
	if err != nil {
		return
	}
	bs := c.(*bitcoin.BTCService)
	// 1. privateKey状态校验
	privateKeyInfo, err := cs.lockinHandler.db.LoadPrivateKeyInfo(cs.meta.SCTokenOwnerKey)
	if err != nil {
		return
	}
	if privateKeyInfo.Status != models.PrivateKeyNegotiateStatusFinished {
		panic("never happen")
	}
	// 2. 估算手续费
	fee := int64(1000)
	// 3. 构造MsgToSign
	msgToSign, err := messagetosign.NewBitcoinLockinTXData(bs, lockinInfo, privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()), fee)
	if err != nil {
		return
	}
	// 4. 签名
	var dsmSignature []byte
	// 考虑到公证人之间区块高度同步的误差导致其他公证人拒绝本次签名,如果签名出错,则重试
	for {
		dsmSignature, _, err = cs.dispatchService.getNotaryService().startDistributedSignAndWait(msgToSign, privateKeyInfo)
		if err != nil {
			log.Error(fmt.Sprintf("callBTCLockin startDistributedSignAndWait error, retry in 10 seconds"))
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}

	// 5. 构造RawTransaction
	tx, err := msgToSign.BuildRawTransaction(dsmSignature)
	if err != nil {
		return
	}
	// 6. 发送交易到链上
	log.Info("call PrepareLockin on bitcoin with account=%s", privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()))
	txHash, err := bs.GetBtcRPCClient().SendRawTransaction(tx, false)
	if err != nil {
		log.Error(fmt.Sprintf("callBTCLockin SendRawTransaction err : %s", err.Error()))
		return
	}
	log.Trace("callBTCLockin txHash=%s", txHash.String())
	return
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
	if localLockoutInfo.MCChainName == events.ChainName {
		return cs.callEthereumPrepareLockout(req, privateKeyInfo, localLockoutInfo)
	}
	if localLockoutInfo.MCChainName == bitcoin.ChainName {
		return cs.callBitcoinPrepareLockout(req, privateKeyInfo, localLockoutInfo)
	}
	err = fmt.Errorf("unknown chain : %s", localLockoutInfo.MCChainName)
	log.Error(err.Error())
	return
}
func (cs *CrossChainService) callBitcoinPrepareLockout(req *userapi.MCPrepareLockoutRequest, privateKeyInfo *models.PrivateKeyInfo, localLockoutInfo *models.LockoutInfo) (err error) {
	// 1. 获取可用utxo
	return
}

func (cs *CrossChainService) callEthereumPrepareLockout(req *userapi.MCPrepareLockoutRequest, privateKeyInfo *models.PrivateKeyInfo, localLockoutInfo *models.LockoutInfo) (err error) {
	// 从本地获取调用合约的参数
	mcUserAddressHex := req.GetSignerETHAddress().String()
	mcExpiration := localLockoutInfo.MCExpiration
	secretHash := localLockoutInfo.SecretHash
	amount := localLockoutInfo.Amount
	// 0. 获取nonce ,这个地方不合理,应该使用侧链Txhash,因为密码可能会重复
	nonce, err := cs.dispatchService.applyNonceFromNonceServer(cs.meta.MCName, privateKeyInfo.Key, req.SecretHash.String(), amount)
	if err != nil {
		return
	}
	// 1. 构造MessageToSign
	var msgToSign messagetosign.MessageToSign
	msgToSign = messagetosign.NewEthereumPrepareLockoutTxData(cs.mcProxy, req, privateKeyInfo.ToAddress(), mcUserAddressHex, secretHash, mcExpiration, amount, nonce)
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
		From:  privateKeyInfo.ToAddress(),
		Nonce: big.NewInt(int64(nonce)),
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
