package bitcoin

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kataras/go-errors"
)

// ChainName 公链名
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
		User:         bs.rpcUser,
		Pass:         bs.rpcPass,
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		DisableTLS:   false, // Bitcoin core does not provide TLS by default
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
	return errors.New("TODO")
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

func (bs *BTCService) loop() {
}
