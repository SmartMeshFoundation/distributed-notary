package dnc

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"os"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
			Usage: "amount of side chain token which you want to lockin, example amount=1 means 1 wei",
		},
		cli.Uint64Flag{
			Name:  "expiration",
			Usage: "expiration of htlc",
			Value: 500,
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
	//account := "default"
	// 1. 构造钱包连接,复用BTCService
	bs, err := bitcoin.NewBTCService(GlobalConfig.BtcWalletRPCEndpoint, GlobalConfig.BtcRPCUser, GlobalConfig.BtcRPCPass, GlobalConfig.BtcWalletRPCCertFilePath)
	if err != nil {
		fmt.Println("NewBTCService err : ", err)
		os.Exit(-1)
	}
	c := bs.GetBtcRPCClient()
	// 2. 解锁钱包
	err = c.WalletPassphrase("123", 1000)
	if err != nil {
		fmt.Println("WalletPassphrase err : ", err)
		os.Exit(-1)
	}
	// 3. 获取双方地址
	userAddress, notaryAddress := getBtcAddresses(bs.GetNetParam())
	// 4. 生成密码及其他参数
	amount2 := btcutil.Amount(amount)
	secret := utils.NewRandomHash()
	secretHash := utils.ShaSecret(secret[:])
	expiration2 := big.NewInt(int64(getBtcLastBlockNumber(c) + expiration))
	fmt.Printf(" ======> [secret=%s, secretHash=%s]\n", secret.String(), secretHash.String())
	// 5. 构造锁定脚本的地址
	scriptBuilder := bs.GetPrepareLockInScriptBuilder(userAddress.(*btcutil.AddressPubKeyHash), notaryAddress.(*btcutil.AddressPubKeyHash), btcutil.Amount(amount), secretHash[:], expiration2)
	lockScript, lockScriptAddr, _ := scriptBuilder.GetPKScript()
	// 6. 发送交易
	err = c.WalletPassphrase("123", 100)
	if err != nil {
		fmt.Println("WalletPassphrase err : ", err)
		os.Exit(-1)
	}
	//utils.PrintBTCBalanceOfAccount(c, "default")
	fmt.Println(amount2)
	txHash, err := c.SendToAddress(lockScriptAddr, amount2)
	if err != nil {
		fmt.Println("SendToAddress err : ", err)
		os.Exit(-1)
	}
	fmt.Printf(" ======> [LockScriptHash=%s, txHash=%s %s]\n", lockScriptAddr.String(), txHash.String(), common.Hash(*txHash).String())
	//time.Sleep(time.Second * 6) // 等待确认
	//utils.PrintBTCBalanceOfAccount(c, "default")

	// 记录runtime数据
	GlobalConfig.RunTime = &runTime{
		Secret:              secret.String(),
		SecretHash:          secretHash.String(),
		BtcLockScript:       lockScript,
		BtcUserAddressBytes: userAddress.ScriptAddress(),
		BtcTXHash:           txHash.String(),
		BtcExpiration:       expiration2,
		BtcAmount:           amount2,
	}
	updateConfigFile()
	fmt.Println("PrepareLockin on bitcoin SUCCESS")
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
	err = cp.PrepareLockin(auth, "", secretHash, expiration2, big.NewInt(amount))
	if err != nil {
		fmt.Println("prepare lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("PrepareLockin on ethereum SUCCESS")
	return nil
}

func getBtcAddresses(net *chaincfg.Params) (userAddress, notaryAddress btcutil.Address) {
	userAddress, err := btcutil.DecodeAddress(GlobalConfig.BtcUserAddress, net)
	if err != nil {
		fmt.Println("DecodeAddress err : ", err)
		os.Exit(-1)
	}
	fmt.Printf(" ======> [userAddress=%s type=%s]\n", userAddress.String(), reflect.TypeOf(userAddress))
	userAddress.ScriptAddress()
	notaryAddress, err = btcutil.DecodeAddress(getSCTokenByMCName(bitcoin.ChainName).MCLockedPublicKeyHashStr, net)
	if err != nil {
		fmt.Println("DecodeAddress err : ", err)
		os.Exit(-1)
	}
	fmt.Printf(" ======> [notaryAddress=%s type=%s]\n", notaryAddress.String(), reflect.TypeOf(notaryAddress))
	return
}
