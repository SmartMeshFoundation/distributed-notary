package main

import (
	"context"
	"math/big"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
	smtcontracts "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/ethereum/go-ethereum/ethclient"
)

const ETHRPC = "http://106.52.171.12:18003"
const SMTRPC = "http://106.52.171.12:17004"

const ETHCONTRACT = "0x326dee230e67e5c124e9c36eae2126c2158bf361"
const SMTCONTRACT = "0x326dee230e67e5c124e9c36eae2126c2158bf361"
const NOTARYADDR = "http://transport01.smartmesh.cn:8032/api/1"

var CallTimeout = time.Second * 5

func bigInt2Ether(n *big.Int) float64 {
	n2 := new(big.Int).Set(n)
	n2 = n2.Div(n2, big.NewInt(1e15))
	n3 := float64(n2.Int64())
	return n3 / 1000.0
}
func main() {
	var ethc *ethclient.Client
	ctx, _ := context.WithTimeout(context.Background(), cfg.ETH.RPCTimeout)
	ethc, err := ethclient.DialContext(ctx, ETHRPC)
	if err != nil {
		panic(err)
	}
	smtClient, err := ethclient.Dial(SMTRPC)
	if err != nil {
		panic(err)
	}
	ethContract, err := contracts.NewLockedEthereum(common.HexToAddress(ETHCONTRACT), ethc)
	if err != nil {
		panic(err)
	}
	name, err := ethContract.Name(nil)
	if err != nil {
		panic(err)
	}
	log.Infof("eth name=%s", name)
	smtContract, err := smtcontracts.NewAtmosphereToken(common.HexToAddress(SMTCONTRACT), smtClient)
	if err != nil {
		panic(err)
	}
	name, err = smtContract.Name(nil)
	if err != nil {
		panic(err)
	}
	log.Infof("smt name=%s", name)
	account := common.HexToAddress("0x201b20123b3c489b47fde27ce5b451a0fa55fd60")
	balance, err := smtClient.BalanceAt(ctx, account, nil)
	if err != nil {
		panic(err)
	}
	ethbalance, err := ethc.BalanceAt(ctx, account, nil)
	if err != nil {
		panic(err)
	}
	log.Infof("account %s eth test balance=%f eth, smt test  balance =%f eth ", account, bigInt2Ether(ethbalance), bigInt2Ether(balance))
	secretHash, expiration, value, err := ethContract.QueryLockin(nil, account)
	if err != nil {
		panic(err)
	}
	lastestHeader, _ := ethc.HeaderByNumber(context.Background(), nil)
	log.Infof("expiration=%s,currentBlockNumber=%s value=%s,secret hash=%s", expiration, lastestHeader.Number, value, common.BytesToHash(secretHash[:]).String())
}
