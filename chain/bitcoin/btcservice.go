package bitcoin

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"time"

	"fmt"

	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kataras/go-errors"
	"github.com/nkbai/log"
)

// ChainName 公链名
var ChainName = "bitcoin"

// BTCOutpointRelevantInfo4PrepareLockout :
type BTCOutpointRelevantInfo4PrepareLockout struct {
	UserAddressPublicKeyHashHex string `json:"user_address_public_key_hash_hex"`
	MCExpiration                uint64 `json:"mc_expiration"`
}

// BTCOutpointRelevantInfo 存储需要监听的outpoint的相关信息
type BTCOutpointRelevantInfo struct {
	SecretHash          common.Hash
	LockScriptHex       string
	Data4PrepareLockout *BTCOutpointRelevantInfo4PrepareLockout
}

// BTCService :
type BTCService struct {
	chainName    string
	host         string
	rpcUser      string
	rpcPass      string
	certificates []byte
	c            *rpcclient.Client
	net          chaincfg.Params

	eventChan             chan chain.Event
	lastBlockNumber       uint64
	eventListenerQuitChan chan struct{}
	confirmMap            map[int32][]*btcutil.Tx
	outpoint2VerifyHexMap map[string]*BTCOutpointRelevantInfo
}

// NewBTCService :
func NewBTCService(host, rpcUser, rpcPass, certFilePath string) (bs *BTCService, err error) {
	bs = &BTCService{
		chainName:             ChainName,
		host:                  host,
		rpcUser:               rpcUser,
		rpcPass:               rpcPass,
		eventChan:             make(chan chain.Event, 100),
		confirmMap:            make(map[int32][]*btcutil.Tx),
		outpoint2VerifyHexMap: make(map[string]*BTCOutpointRelevantInfo),
	}
	// #nosec
	certs, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return
	}
	bs.certificates = certs
	connCfg := &rpcclient.ConnConfig{
		Host:         bs.host,
		Endpoint:     "ws",
		User:         bs.rpcUser,
		Pass:         bs.rpcPass,
		Certificates: certs,
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	bs.c, err = rpcclient.New(connCfg, &rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: bs.onFilteredBlockConnected,
		OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			fmt.Println("disconnected")
		},
	})
	if err != nil {
		return
	}
	info, err := bs.c.GetBlockChainInfo()
	if err != nil {
		return
	}
	switch info.Chain {
	case chaincfg.MainNetParams.Name:
		bs.net = chaincfg.MainNetParams
	case chaincfg.TestNet3Params.Name:
		bs.net = chaincfg.TestNet3Params
	case chaincfg.SimNetParams.Name:
		bs.net = chaincfg.SimNetParams
	case chaincfg.RegressionNetParams.Name:
		bs.net = chaincfg.RegressionNetParams
	default:
		err = fmt.Errorf("unknown bitcoin network : %s", info.Chain)
	}
	// 开启BlockNumber订阅
	err = bs.c.NotifyBlocks()
	if err != nil {
		return
	}
	return
}

// GetChainName impl chain.Chain
func (bs *BTCService) GetChainName() string {
	return bs.chainName
}

// GetEventChan impl chain.Chain
func (bs *BTCService) GetEventChan() <-chan chain.Event {
	return bs.eventChan
}

// StartEventListener impl chain.Chain TODO
func (bs *BTCService) StartEventListener() error {
	if bs.eventListenerQuitChan != nil {
		return errors.New("event listener already start")
	}
	bs.eventListenerQuitChan = make(chan struct{})
	return nil
}

// StopEventListener impl chain.Chain
func (bs *BTCService) StopEventListener() {
	if bs.eventListenerQuitChan != nil {
		close(bs.eventListenerQuitChan)
		bs.eventListenerQuitChan = nil
	}
	return
}

// RegisterEventListenContract impl chain.Chain
func (bs *BTCService) RegisterEventListenContract(contractAddresses ...common.Address) error {
	// do nothing
	return nil
}

// UnRegisterEventListenContract impl chain.Chain
func (bs *BTCService) UnRegisterEventListenContract(contractAddresses ...common.Address) {
	// do nothing
}

// DeployContract impl chain.Chain
func (bs *BTCService) DeployContract(opts *bind.TransactOpts, params ...string) (contractAddress common.Address, err error) {
	// do nothing
	return
}

// SetLastBlockNumber impl chain.Chain
func (bs *BTCService) SetLastBlockNumber(lastBlockNumber uint64) {
	bs.lastBlockNumber = lastBlockNumber
}

// GetContractProxy impl chain.Chain TODO
func (bs *BTCService) GetContractProxy(contractAddress common.Address) chain.ContractProxy {
	return nil
}

// GetConn impl chain.Chain
func (bs *BTCService) GetConn() *ethclient.Client {
	panic(chain.ErrorCallWrongChain.Error())
}

// Transfer10ToAccount impl chain.Chain, for debug
func (bs *BTCService) Transfer10ToAccount(key *ecdsa.PrivateKey, accountTo common.Address, amount *big.Int, nonce ...int) (err error) {
	// do nothing
	return
}

// GetNetParam :
func (bs *BTCService) GetNetParam() *chaincfg.Params {
	return &bs.net
}

// GetBtcRPCClient :
func (bs *BTCService) GetBtcRPCClient() *rpcclient.Client {
	return bs.c
}

// RegisterOutpoint 注册outpoint监听
func (bs *BTCService) RegisterOutpoint(outpoint wire.OutPoint, relevantInfo *BTCOutpointRelevantInfo) (err error) {
	outPointKey := getOutpointKey(outpoint)
	bs.outpoint2VerifyHexMap[outPointKey] = relevantInfo
	return bs.c.LoadTxFilter(false, nil, []wire.OutPoint{outpoint})
}

// BtcPrepareLockinInfo :
type BtcPrepareLockinInfo struct {
	TxHash          chainhash.Hash `json:"tx_hash"`
	Index           int            `json:"index"`
	BlockNumber     uint64         `json:"block_number"`
	BlockNumberTime int64          `json:"block_number_time"`
}

/*
GetPrepareLockinInfo :
*/
func (bs *BTCService) GetPrepareLockinInfo(txHash chainhash.Hash, lockAddr string, lockAmount btcutil.Amount) (res *BtcPrepareLockinInfo, err error) {
	tx, err := bs.c.GetRawTransactionVerbose(&txHash)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, txOut := range tx.Vout {
		// 这里tx里面的amount单位为btc
		if txOut.Value == lockAmount.ToBTC() && txOut.ScriptPubKey.Addresses[0] == lockAddr {
			res = &BtcPrepareLockinInfo{
				TxHash: txHash,
				Index:  int(txOut.N),
			}
			var blockHash chainhash.Hash
			err2 := chainhash.Decode(&blockHash, tx.BlockHash)
			if err2 != nil {
				log.Error(err2.Error())
				return nil, err2
			}
			block, err2 := bs.c.GetBlockVerbose(&blockHash)
			if err2 != nil {
				log.Error(err2.Error())
				return nil, err2
			}
			res.BlockNumberTime = block.Time
			res.BlockNumber = uint64(block.Height)
			return
		}
	}
	err = fmt.Errorf("can not found PrepareLockinInfo on bitcoin : txHash=%s, lockAddr=%s, lockAmount=%d tx=\n%s", txHash.String(), lockAddr, lockAmount, utils.ToJSONStringFormat(tx))
	return
}

/*
块销毁处理
*/
func (bs *BTCService) onFilteredBlockDisconnected(height int32, header *wire.BlockHeader) {
	txs, ok := bs.confirmMap[height]
	if ok {
		log.Info("block %d disconnected, %d txs delete in confirm map", height, len(txs))
		delete(bs.confirmMap, height)
	} else {
		log.Info("block %d disconnected")
	}
}

/*
新块处理
*/
func (bs *BTCService) onFilteredBlockConnected(height int32, header *wire.BlockHeader, txs []*btcutil.Tx) {
	if height%logPeriod == 0 {
		log.Trace(fmt.Sprintf("Bitcoin  new block : %d", height))
	}
	// 1. 生成NewBlock事件
	bs.lastBlockNumber = uint64(height)
	bs.eventChan <- createNewBlockEvent(bs.lastBlockNumber)
	// 2. 保存tx到确认池
	if len(txs) > 0 {
		bs.confirmMap[height] = txs
		log.Trace("Bitcoin  new block : %d relevant tx num: %d", height, len(txs))
	}
	// 3. 确认事件
	var eventsToSend []chain.Event
	var confirmBlockNumbers []int32
	for blockNumber := range bs.confirmMap {
		if blockNumber <= height-confirmBlockNumber {
			for _, tx := range bs.confirmMap[blockNumber] {
				eventsToSend = append(eventsToSend, bs.createEventsFromTx(blockNumber, tx)...)
			}
			confirmBlockNumbers = append(confirmBlockNumbers, blockNumber)
		}
	}
	for _, n := range confirmBlockNumbers {
		delete(bs.confirmMap, n)
	}
	// 4. 投递事件
	for _, e := range eventsToSend {
		bs.eventChan <- e
	}
}

/*
已确认的交易处理
*/
func (bs *BTCService) createEventsFromTx(blockNumber int32, tx *btcutil.Tx) (events []chain.Event) {
	// 1. 获取utxo对应的RelevantInfo
	var outpointRelevantInfo *BTCOutpointRelevantInfo
	var txIn *wire.TxIn
	ok := false
	for _, in := range tx.MsgTx().TxIn {
		outpointKey := getOutpointKey(in.PreviousOutPoint)
		outpointRelevantInfo, ok = bs.outpoint2VerifyHexMap[outpointKey]
		if ok {
			txIn = in
			break
		}
	}
	if !ok {
		log.Error("can not find BTCOutpointRelevantInfo of tx :\n%s", utils.ToJSONStringFormat(tx))
		return
	}
	// 判断tx对应的事件
	if bs.isLockin(txIn, outpointRelevantInfo) {
		// 用户-主链lockin
		e := createLockinEvent(uint64(blockNumber), tx.Hash().String(), outpointRelevantInfo.SecretHash, tx.MsgTx().TxOut)
		events = append(events, e)
		return
	}
	if bs.isCancelPrepareLockin(txIn, outpointRelevantInfo) {
		// 自己-主链CancelLockin
		e := createCancelLockinEvent(uint64(blockNumber), outpointRelevantInfo.SecretHash)
		events = append(events, e)
		return
	}
	if bs.isPrepareLockout(tx.MsgTx(), txIn, outpointRelevantInfo) {
		// 自己-主链PrepareLockout
		e := createPrepareLockoutEvent(uint64(blockNumber), tx.MsgTx().TxHash(), outpointRelevantInfo)
		events = append(events, e)
		return
	}
	//if secret, isLockout := bs.isLockout(txIn, outpointRelevantInfo); isLockout {
	//	// 主链lockout
	//	log.Info("收到BTC Lockout 事件,secret=%s secretHash=%s", secret.String(), utils.ShaSecret(secret[:]).String())
	//	return
	//}
	//if bs.isCancelPrepareLockout(tx.MsgTx()) {
	//	return
	//}
	log.Error("receive unknown tx : \n%s", utils.ToJSONStringFormat(tx))
	return
}

/*
用户CancelPrepareLockin的tx
SigScript : SIG {{用户PKH}} 0 {{锁定脚本}}
*/
func (bs *BTCService) isCancelPrepareLockin(txIn *wire.TxIn, outpointRelevantInfo *BTCOutpointRelevantInfo) (isCancelPrepareLockin bool) {
	// 验证部分: 0 {{锁定脚本}}
	verifyStr := fmt.Sprintf("0 %s", outpointRelevantInfo.LockScriptHex)
	signatureScriptStr, err := txscript.DisasmString(txIn.SignatureScript)
	if err != nil {
		log.Error("DisasmString SignatureScript err : %s", err.Error())
		return
	}
	isCancelPrepareLockin = strings.Contains(signatureScriptStr, verifyStr)
	return
}

/*
自己Lockin的tx
txIn 只有一个 SigScript : SIG {{分布式私钥的PKH}} {{SECRET}} 1 {{锁定脚本}}
*/
func (bs *BTCService) isLockin(txIn *wire.TxIn, outpointRelevantInfo *BTCOutpointRelevantInfo) (ok bool) {
	// 验证部分: 1 {{锁定脚本}}
	verifyStr := fmt.Sprintf("1 %s", outpointRelevantInfo.LockScriptHex)
	signatureScriptStr, err := txscript.DisasmString(txIn.SignatureScript)
	if err != nil {
		log.Error("DisasmString SignatureScript err : %s", err.Error())
		return
	}
	if !strings.Contains(signatureScriptStr, verifyStr) {
		return
	}
	ok = true
	return
}

/*
自己PrepareLockout的tx
txIns 数量不确定,但没有P2SH,SigScipt : SIG {{分布式私钥的PKH}}
*/
func (bs *BTCService) isPrepareLockout(tx *wire.MsgTx, txIn *wire.TxIn, outpointRelevantInfo *BTCOutpointRelevantInfo) bool {
	if len(tx.TxOut) != 1 {
		return false
	}
	// 验证部分: {{分布式私钥的PKH}}
	signatureScriptStr, err := txscript.DisasmString(txIn.SignatureScript)
	if err != nil {
		log.Error("DisasmString SignatureScript err : %s", err.Error())
		return false
	}
	if !strings.Contains(signatureScriptStr, outpointRelevantInfo.LockScriptHex) {
		return false
	}
	return true
}

/*
自己CancenlPrepareLockout的tx
SigScript : SIG {{分布式私钥的PKH}} 0 {{锁定脚本}}
*/
func (bs *BTCService) isCancelPrepareLockout(txIn *wire.TxIn, outpointRelevantInfo *BTCOutpointRelevantInfo) bool {
	// 验证部分: 1 {{锁定脚本}}
	verifyStr := fmt.Sprintf("0 %s", outpointRelevantInfo.LockScriptHex)
	signatureScriptStr, err := txscript.DisasmString(txIn.SignatureScript)
	if err != nil {
		log.Error("DisasmString SignatureScript err : %s", err.Error())
		return false
	}
	return strings.Contains(signatureScriptStr, verifyStr)
}

/*
用户Lockout的tx
SigScript : SIG {{用户的PKH}} {{SECRET}} 1 {{锁定脚本}}
*/
func (bs *BTCService) isLockout(txIn *wire.TxIn, outpointRelevantInfo *BTCOutpointRelevantInfo) (secret common.Hash, ok bool) {
	// 验证部分: 1 {{锁定脚本}}
	verifyStr := fmt.Sprintf("1 %s", outpointRelevantInfo.LockScriptHex)
	signatureScriptStr, err := txscript.DisasmString(txIn.SignatureScript)
	if err != nil {
		log.Error("DisasmString SignatureScript err : %s", err.Error())
		return
	}
	if !strings.Contains(signatureScriptStr, verifyStr) {
		return
	}
	// 验证通过,解析secret
	ss := strings.Split(signatureScriptStr, " ")
	secret = common.HexToHash(ss[2])
	ok = true
	return
}

func getOutpointKey(outpoint wire.OutPoint) string {
	return fmt.Sprintf("%s-%d", outpoint.Hash.String(), outpoint.Index)
}
