package service

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"

	"math/big"

	"time"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
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
	sc             chain.Chain
	mc             chain.Chain
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
	scChain, err := dispatchService.getChainByName(cfg.HECO.Name)
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
		sc:              scChain,
		mc:              mcChain,
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
	scUserAddressHex := req.GetSignerSMCAddress().String()
	scExpiration := localLockinInfo.SCExpiration
	secretHash := localLockinInfo.SecretHash
	amount := new(big.Int).Sub(localLockinInfo.Amount, localLockinInfo.CrossFee) // 扣除手续费
	// 0. 获取nonce
	nonce, err := cs.dispatchService.applyNonceFromNonceServer(cfg.HECO.Name, privateKeyInfo.Key, req.SecretHash.String(), amount)
	if err != nil {
		return
	}
	// 1. 构造MessageToSign
	var msgToSign messagetosign.MessageToSign
	msgToSign = messagetosign.NewHecoPrepareLockinTxData(cs.scTokenProxy, req, privateKeyInfo.ToAddress(), scUserAddressHex, secretHash, scExpiration, amount, nonce)
	// 2. 发起分布式签名
	var signature []byte
	var _ common.Hash
	signature, _, err = cs.dispatchService.getNotaryService().startDistributedSignAndWait(msgToSign, privateKeyInfo)
	if err != nil {
		return
	}
	log.Info("call PrepareLockin on heco with account=%s, signature=%s", privateKeyInfo.ToAddress().String(), common.Bytes2Hex(signature))
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
	//log.Trace(fmt.Sprintf("===>now callSCPrepareLockin,CrossChainService=%s", utils.StringInterface(cs, 5)))
	log.Trace(fmt.Sprintf("===>now callSCPrepareLockin,scUserAddressHex=%s ,secretHash=%s ,scExpiration=%d ,amount=%s", scUserAddressHex, secretHash.Hex(), scExpiration, amount.String()))
	return cs.scTokenProxy.PrepareLockin(transactor, scUserAddressHex, secretHash, scExpiration, amount)
}

func (cs *CrossChainService) callMCLockin(lockinInfo *models.LockinInfo) (err error) {
	// 无需使用分布式签名,用自己的签名就好
	if cs.meta.MCName == cfg.SMC.Name {
		// 以太坊直接使用自己私钥调用合约即可
		auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
		return cs.mcProxy.Lockin(auth, lockinInfo.MCUserAddressHex, lockinInfo.Secret)
	}
	return errors.New("unknown chain")
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
	if localLockoutInfo.MCChainName == cfg.SMC.Name {
		return cs.callSpectrumPrepareLockout(req, privateKeyInfo, localLockoutInfo)
	}
	err = fmt.Errorf("unknown chain : %s", localLockoutInfo.MCChainName)
	log.Error(err.Error())
	return
}

func (cs *CrossChainService) callSpectrumPrepareLockout(req *userapi.MCPrepareLockoutRequest, privateKeyInfo *models.PrivateKeyInfo, localLockoutInfo *models.LockoutInfo) (err error) {
	// 从本地获取调用合约的参数
	mcUserAddressHex := req.GetSignerSMCAddress().String()
	mcExpiration := localLockoutInfo.MCExpiration
	secretHash := localLockoutInfo.SecretHash
	//amount := localLockoutInfo.Amount
	amount := new(big.Int).Sub(localLockoutInfo.Amount, localLockoutInfo.CrossFee) // 扣除手续费
	// 0. 获取nonce ,这个地方不合理,应该使用侧链Txhash,因为密码可能会重复
	nonce, err := cs.dispatchService.applyNonceFromNonceServer(cs.meta.MCName, privateKeyInfo.Key, req.SecretHash.String(), amount)
	if err != nil {
		return
	}
	// 1. 构造MessageToSign
	var msgToSign messagetosign.MessageToSign
	msgToSign = messagetosign.NewSpectrumPrepareLockoutTxData(cs.mcProxy, req, privateKeyInfo.ToAddress(), mcUserAddressHex, secretHash, mcExpiration, amount, nonce)
	// 2. 发起分布式签名
	var signature []byte
	var _ common.Hash
	signature, _, err = cs.dispatchService.getNotaryService().startDistributedSignAndWait(msgToSign, privateKeyInfo)
	if err != nil {
		return
	}
	log.Info("call PrepareLockout on spectrum with account=%s, signature=%s", privateKeyInfo.ToAddress().String(), common.Bytes2Hex(signature))
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
	if cs.meta.MCName == cfg.SMC.Name {
		return cs.callSpectrumCancelLockout(userAddressHex)
	}
	//if cs.meta.MCName == cfg.BTC.Name {
	//	return cs.callBitcoinCancelLockout(userAddressHex)
	//}
	err = fmt.Errorf("unknown chain : %s", cs.meta.MCName)
	log.Error(err.Error())
	return
}

func (cs *CrossChainService) callSpectrumCancelLockout(userAddressHex string) (err error) {
	// spectrum 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.scTokenProxy.CancelLockout(auth, userAddressHex)
}

/*
根据侧链超时事件计算主链超时时间
*/
func (cs *CrossChainService) calculateMCExpiration(scExpiration uint64) (mcExpiration uint64) {
	scCfg := cfg.GetCfgByChainName(cs.sc.GetChainName())
	mcCfg := cfg.GetCfgByChainName(cs.mc.GetChainName())
	scExpirationSecond := time.Duration(scExpiration-cs.scLastedBlockNumber) * scCfg.BlockPeriod
	mcExpirationSecond := scExpirationSecond / 2
	mcExpirationBlocks := int64(mcExpirationSecond / mcCfg.BlockPeriod)
	mcExpiration = cs.mcLastedBlockNumber + uint64(mcExpirationBlocks)
	return
}

/*
根据主链超时事件计算侧链超时时间
*/
func (cs *CrossChainService) calculateSCExpiration(mcExpiration uint64) (scExpiration uint64) {
	scCfg := cfg.GetCfgByChainName(cs.sc.GetChainName())
	mcCfg := cfg.GetCfgByChainName(cs.mc.GetChainName())
	mcExpirationSecond := time.Duration(mcExpiration-cs.mcLastedBlockNumber) * mcCfg.BlockPeriod
	scExpirationSecond := mcExpirationSecond / 2
	scExpirationBlocks := int64(scExpirationSecond / scCfg.BlockPeriod)
	scExpiration = cs.scLastedBlockNumber + uint64(scExpirationBlocks)
	return
}

func (cs *CrossChainService) getLockInInfoBySCPrepareLockInRequest(req *userapi.SCPrepareLockinRequest) (lockinInfo *models.LockinInfo, err error) {
	if cs.meta.MCName == cfg.SMC.Name {
		// SMC,收到用户请求时,lockinInfo必须已经存在,仅做校验工作
		mcUserAddress := common.BytesToAddress(req.MCUserAddress)
		lockinInfo, err = cs.lockinHandler.getLockin(req.SecretHash)
		if err != nil {
			req.WriteErrorResponse(api.ErrorCodeException, err.Error())
			return
		}
		if lockinInfo.MCUserAddressHex != mcUserAddress.String() {
			err = errors.New("MCUserAddress wrong")
			return
		}
		if lockinInfo.MCLockStatus != models.LockStatusLock {
			err = fmt.Errorf("MCLockStatus %d wrong", lockinInfo.MCLockStatus)
			return
		}
		if lockinInfo.SCLockStatus != models.LockStatusNone {
			err = fmt.Errorf("SCLockStatus %d wrong", lockinInfo.SCLockStatus)
			return
		}
		return
	}
	err = fmt.Errorf("unknown chain : %s", cs.meta.MCName)
	return
}
