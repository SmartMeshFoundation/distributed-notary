package service

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"errors"

	"math/big"

	"time"

	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
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
	scChain, err := dispatchService.getChainByName(cfg.SMC.Name)
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
	scUserAddressHex := req.GetSignerETHAddress().String()
	scExpiration := localLockinInfo.SCExpiration
	secretHash := localLockinInfo.SecretHash
	amount := new(big.Int).Sub(localLockinInfo.Amount, localLockinInfo.CrossFee) // 扣除手续费
	// 0. 获取nonce
	nonce, err := cs.dispatchService.applyNonceFromNonceServer(cfg.SMC.Name, privateKeyInfo.Key, req.SecretHash.String(), amount)
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
	if cs.meta.MCName == cfg.ETH.Name {
		// 以太坊直接使用自己私钥调用合约即可
		auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
		return cs.mcProxy.Lockin(auth, lockinInfo.MCUserAddressHex, lockinInfo.Secret)
	} else if cs.meta.MCName == cfg.BTC.Name {
		// 比特币的lockin操作需要分布式签名
		return cs.callBTCLockin(lockinInfo)
	}
	return errors.New("unknown chain")
}

func (cs *CrossChainService) callBTCLockin(lockinInfo *models.LockinInfo) (err error) {
	// 0. 获取bs
	bs := cs.mc.(*bitcoin.BTCService)
	// 1. privateKey状态校验
	privateKeyInfo, err := cs.lockinHandler.db.LoadPrivateKeyInfo(cs.meta.SCTokenOwnerKey)
	if err != nil {
		return
	}
	if privateKeyInfo.Status != models.PrivateKeyNegotiateStatusFinished {
		panic("never happen")
	}
	// 2. 估算主链交易手续费
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
	if localLockoutInfo.MCChainName == cfg.ETH.Name {
		return cs.callEthereumPrepareLockout(req, privateKeyInfo, localLockoutInfo)
	}
	if localLockoutInfo.MCChainName == cfg.BTC.Name {
		return cs.callBitcoinPrepareLockout(req, privateKeyInfo, localLockoutInfo)
	}
	err = fmt.Errorf("unknown chain : %s", localLockoutInfo.MCChainName)
	log.Error(err.Error())
	return
}
func (cs *CrossChainService) callBitcoinPrepareLockout(req *userapi.MCPrepareLockoutRequest, privateKeyInfo *models.PrivateKeyInfo, localLockoutInfo *models.LockoutInfo) (err error) {
	amount := localLockoutInfo.Amount
	// 0. 获取可用utxo
	utxoKeysStr, err := cs.dispatchService.applyUTXO(cs.meta.MCName, privateKeyInfo.Key, req.SecretHash.String(), amount)
	if err != nil {
		return
	}
	// 1. 获取bs
	bs := cs.mc.(*bitcoin.BTCService)
	// 1.5 获取fee
	fee := int64(1000)
	// 2. 对需要使用的每个utxo发起DSM
	utxoKeys := strings.Split(utxoKeysStr, "-")
	var tx *wire.MsgTx
	for index := range utxoKeys {
		//构造
		msgToSign, err2 := messagetosign.NewBitcoinPrepareLockoutTXData(req, bs, localLockoutInfo, privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()), cs.lockoutHandler.db, utxoKeysStr, fee, index)
		if err2 != nil {
			err = err2
			return
		}
		if tx == nil {
			// 保存一份originTx,用于收集填写Signature
			tx = msgToSign.GetOriginTxCopy()
		}
		// 签名
		var signature []byte
		signature, _, err = cs.dispatchService.getNotaryService().startDistributedSignAndWait(msgToSign, privateKeyInfo)
		if err != nil {
			return
		}
		// 按比特币标准调整签名
		tx.TxIn[index].SignatureScript = msgToSign.BuildBTCSignatureScript(signature)
		// 注册outpoint监听
		outpointToListen := msgToSign.GetOriginTxCopy().TxIn[index].PreviousOutPoint
		err = bs.RegisterOutpoint(outpointToListen, &bitcoin.BTCOutpointRelevantInfo{
			Use:           bitcoin.OutpointUseToPrepareLockout,
			SecretHash:    req.SecretHash,
			LockScriptHex: common.Bytes2Hex(privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()).PubKey().SerializeCompressed()),
			Data4PrepareLockout: &bitcoin.BTCOutpointRelevantInfo4PrepareLockout{
				// 保存用户主链取钱的地址
				UserAddressPublicKeyHashHex: req.GetSignerBTCPublicKey(bs.GetNetParam()).AddressPubKeyHash().String(),
				MCExpiration:                localLockoutInfo.MCExpiration,
				TxOutLockScriptHex:          msgToSign.GetLockScriptHex(),
			},
		})
		if err != nil {
			log.Error(err.Error())
		}
	}
	// 4. 发送交易
	log.Info("call PrepareLockout on bitcoin with account=%s", privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()))
	txHash, err := bs.GetBtcRPCClient().SendRawTransaction(tx, false)
	if err != nil {
		log.Error(fmt.Sprintf("callBTCLockin SendRawTransaction err : %s", err.Error()))
		return
	}
	log.Trace("callBTCPrepareLockout txHash=%s", txHash.String())
	return
}

func (cs *CrossChainService) callEthereumPrepareLockout(req *userapi.MCPrepareLockoutRequest, privateKeyInfo *models.PrivateKeyInfo, localLockoutInfo *models.LockoutInfo) (err error) {
	// 从本地获取调用合约的参数
	mcUserAddressHex := req.GetSignerETHAddress().String()
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
	if cs.meta.MCName == cfg.ETH.Name {
		return cs.callEthereumCancelLockout(userAddressHex)
	}
	//if cs.meta.MCName == cfg.BTC.Name {
	//	return cs.callBitcoinCancelLockout(userAddressHex)
	//}
	err = fmt.Errorf("unknown chain : %s", cs.meta.MCName)
	log.Error(err.Error())
	return
}

func (cs *CrossChainService) callEthereumCancelLockout(userAddressHex string) (err error) {
	// 以太坊 无需使用分布式签名,用自己的签名就好
	auth := bind.NewKeyedTransactor(cs.selfPrivateKey)
	return cs.scTokenProxy.CancelLockout(auth, userAddressHex)
}

func (cs *CrossChainService) callBitcoinCancelLockout(lockoutInfo *models.LockoutInfo) (err error) {
	// 0. privateKey状态校验
	privateKeyInfo, err := cs.lockinHandler.db.LoadPrivateKeyInfo(cs.meta.SCTokenOwnerKey)
	if err != nil {
		return
	}
	if privateKeyInfo.Status != models.PrivateKeyNegotiateStatusFinished {
		panic("never happen")
	}
	// 1. 估算手续费
	fee := int64(1000)
	// 2. 构造MsgToSign
	bs := cs.mc.(*bitcoin.BTCService)
	msgToSign := messagetosign.NewBitcoinCancelPrepareLockoutTXData(lockoutInfo, privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()), fee)
	// 3. 签名
	dsmSignature, _, err := cs.dispatchService.getNotaryService().startDistributedSignAndWait(msgToSign, privateKeyInfo)
	// 4. 构造RawTransaction
	tx, err := msgToSign.BuildRawTransaction(dsmSignature)
	if err != nil {
		return
	}
	// 6. 发送交易到链上
	log.Info("call CancelPrepareLockout on bitcoin with account=%s lockTime=%d", privateKeyInfo.ToBTCPubKeyAddress(bs.GetNetParam()), tx.LockTime)
	txHash, err := bs.GetBtcRPCClient().SendRawTransaction(tx, false)
	if err != nil {
		log.Error(fmt.Sprintf("callBitcoinCancelLockout SendRawTransaction err : %s", err.Error()))
		return
	}
	log.Trace("callBitcoinCancelLockout txHash=%s", txHash.String())
	return
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
	if cs.meta.MCName == cfg.ETH.Name {
		// 以太坊,收到用户请求时,lockinInfo必须已经存在,仅做校验工作
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
	if cs.meta.MCName == cfg.BTC.Name {
		btcService := cs.mc.(*bitcoin.BTCService)
		net := btcService.GetNetParam()
		mcUserAddress, err2 := btcutil.NewAddressPubKeyHash(req.MCUserAddress, net)
		if err2 != nil {
			log.Error(err2.Error())
			return nil, err2
		}
		notaryPublicKeyHash, err2 := btcutil.DecodeAddress(cs.meta.MCLockedPublicKeyHashStr, net)
		if err2 != nil {
			panic("never happen")
		}
		notaryAddress := notaryPublicKeyHash.(*btcutil.AddressPubKeyHash)
		// 比特币,收到用户请求之前,公证人无从得知任何信息,所以这里必须根据用户的请求验证链上数据,并构造lockinInfo
		// 1. 根据用户参数本地生成脚本hash
		builder := btcService.GetPrepareLockInScriptBuilder(mcUserAddress, notaryAddress, req.MCLockedAmount, req.SecretHash[:], req.MCExpiration)
		lockScript, lockAddr, _ := builder.GetPKScript()
		// 2. 从链上查询tx并在其中查找锁定地址为上一步生成的hash值的outpoint及交易发生的blockNumber
		mcTXHash, err2 := chainhash.NewHash(req.MCTXHash)
		if err2 != nil {
			log.Error(err2.Error())
			return nil, err2
		}
		btcPrepareLockinInfo, err2 := btcService.GetPrepareLockinInfo(*mcTXHash, lockAddr.String(), req.MCLockedAmount)
		if err2 != nil {
			return nil, err2
		}
		// 3. 校验mcExpiration
		mcExpiration := req.MCExpiration.Uint64()
		if mcExpiration-btcPrepareLockinInfo.BlockNumber <= cfg.GetMinExpirationBlock4User(cfg.BTC.Name) {
			err2 = fmt.Errorf("mcExpiration must bigger than %d", cfg.GetMinExpirationBlock4User(cfg.BTC.Name))
			return nil, err2
		}
		// 4 计算scExpiration
		scExpiration := cs.calculateSCExpiration(mcExpiration)
		// 5. 构造LockinInfo
		amount := big.NewInt(int64(req.MCLockedAmount))
		lockinInfo = &models.LockinInfo{
			MCChainName:      cs.meta.MCName,
			SecretHash:       req.SecretHash,
			Secret:           utils.EmptyHash,
			MCUserAddressHex: mcUserAddress.String(),
			SCUserAddress:    req.GetSignerETHAddress(),
			SCTokenAddress:   cs.meta.SCToken,
			Amount:           amount,
			MCExpiration:     mcExpiration,
			SCExpiration:     scExpiration,
			MCLockStatus:     models.LockStatusLock,
			SCLockStatus:     models.LockStatusNone,
			//Data:               data,
			NotaryIDInCharge:          models.UnknownNotaryIDInCharge,
			StartTime:                 btcPrepareLockinInfo.BlockNumberTime,
			StartMCBlockNumber:        btcPrepareLockinInfo.BlockNumber,
			BTCPrepareLockinTXHashHex: btcPrepareLockinInfo.TxHash.String(),
			BTCPrepareLockinVout:      uint32(btcPrepareLockinInfo.Index),
			CrossFee:                  cs.dispatchService.calculateCrossFee(cs.meta.MCName, amount), // 计算跨链手续费
		}
		// 6. 调用handler处理
		err2 = cs.lockinHandler.registerLockin(lockinInfo)
		if err2 != nil {
			err2 = fmt.Errorf("lockinHandler.registerLockin err2 = %s", err2.Error())
			return nil, err2
		}
		// 7. 注册需要监听的outpoint到BTCService
		err2 = btcService.RegisterOutpoint(wire.OutPoint{
			Hash:  btcPrepareLockinInfo.TxHash,
			Index: lockinInfo.BTCPrepareLockinVout,
		}, &bitcoin.BTCOutpointRelevantInfo{
			Use:           bitcoin.OutpointUseToLockinOrCancel,
			SecretHash:    req.SecretHash,
			LockScriptHex: common.Bytes2Hex(lockScript),
		})
		if err2 != nil {
			log.Error("lockinHandler.RegisterOutpoint err2 = %s", err2.Error())
			// 这里不返回错误,注册失败不能影响操作
		}
		return
	}
	err = fmt.Errorf("unknown chain : %s", cs.meta.MCName)
	return
}
