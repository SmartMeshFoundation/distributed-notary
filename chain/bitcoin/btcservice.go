package bitcoin

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/Spectrum/log"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kataras/go-errors"
)

// MCChainName 公链名
var ChainName = "bitcoin"

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
}

// NewBTCService :
func NewBTCService(host, rpcUser, rpcPass, certFilePath string) (bs *BTCService, err error) {
	bs = &BTCService{
		chainName: ChainName,
		host:      host,
		rpcUser:   rpcUser,
		rpcPass:   rpcPass,
		eventChan: make(chan chain.Event, 100),
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
	bs.c, err = rpcclient.New(connCfg, nil)
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

func (bs *BTCService) loop() {
}
