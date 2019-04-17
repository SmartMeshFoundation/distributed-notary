package dnc

import (
	"context"
	"fmt"

	"time"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kataras/go-errors"
	"github.com/urfave/cli"
)

var pliCmd = cli.Command{
	Name:      "prepare-lock-in",
	ShortName: "pli",
	Usage:     "call main chain contract prepare lock in",
	Action:    prepareLockin,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "ethereum",
		},
		cli.Int64Flag{
			Name:  "amount",
			Usage: "amount of side chain token which you want to lockin, example amount=1 means 1eth",
		},
		cli.Uint64Flag{
			Name:  "expiration",
			Usage: "expiration of htlc",
			Value: 1000,
		},
	},
}

func prepareLockin(ctx *cli.Context) error {
	mcName := ctx.String("mcname")
	amount := ctx.Int64("amount")
	if amount == 0 {
		fmt.Println("pli must run with --amount")
		os.Exit(-1)
	}
	expiration := ctx.Uint64("expiration")
	fmt.Printf("start to prepare lockin :\n ======> [chain=%s amount=%d expiartion=%d]\n", mcName, amount, expiration)
	if mcName == events.ChainName {
		return prepareLockinOnEthereum(mcName, amount, expiration)
	} else if mcName == bitcoin.ChainName || mcName == "btc" {
		return prepareLockinOnBitcoin(amount, expiration)
	}
	return errors.New("unknown chain name")
}

func prepareLockinOnBitcoin(amount int64, expiration uint64) (err error) {
	//// 1. connect to btcd
	//certs, err := ioutil.ReadFile(filepath.Join("/home/chuck/.btcd", "rpc.cert"))
	//if err != nil {
	//	fmt.Println("get certs of btc err : ", err)
	//	os.Exit(-1)
	//}
	//connCfg := &rpcclient.ConnConfig{
	//	Host:         "192.168.124.13:18334",
	//	User:         "bai",
	//	Pass:         "bai",
	//	HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
	//	DisableTLS:   true, // Bitcoin core does not provide TLS by default
	//	Certificates: certs,
	//}
	//c, err := rpcclient.New(connCfg, nil)
	//if err != nil {
	//	fmt.Println("connect to btc err : ", err)
	//	os.Exit(-1)
	//}
	//// 2. get auth
	//privateKeyT, err := getPrivateKey(GlobalConfig.EthUserAddress, GlobalConfig.EthUserPassword)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	//privateKey := (*btcec.PrivateKey)(privateKeyT)
	//// 获取账号utxo集
	//wire.NewMsgTx(1)
	//c.CreateEncryptedWallet()
	//c.SendRawTransaction()
	return
}

func prepareLockinOnEthereum(mcName string, amount int64, expiration uint64) (err error) {
	contract := getMCContractAddressByMCName(mcName)
	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.EthRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	conn := client.NewSafeClient(c)

	// 2. init contract proxy
	cp, err := proxy.NewLockedEthereumProxy(conn, contract)
	if err != nil {
		fmt.Println("init contract proxy err : ", err)
		os.Exit(-1)
	}
	// 3. get auth
	privateKey, err := getPrivateKey(GlobalConfig.EthUserAddress, GlobalConfig.EthUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call pli
	auth := bind.NewKeyedTransactor(privateKey)
	secret := utils.NewRandomHash()
	secretHash := utils.ShaSecret(secret[:])
	expiration2 := getEthLastBlockNumber(conn) + expiration
	fmt.Printf(" ======> [secret=%s, secretHash=%s]\n", secret.String(), secretHash.String())
	GlobalConfig.RunTime = &runTime{
		Secret:     secret.String(),
		SecretHash: secretHash.String(),
	}
	updateConfigFile()
	err = cp.PrepareLockin(auth, "", secretHash, expiration2, eth2Wei(amount))
	if err != nil {
		fmt.Println("prepare lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("PrepareLockin SUCCESS")
	return nil
}
