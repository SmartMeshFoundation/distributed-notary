package dnc

import (
	"context"
	"fmt"
	"os"
	"time"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain/heco/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/heco/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

var ploCmd = cli.Command{
	Name:      "prepare-lock-out",
	ShortName: "plo",
	Usage:     "call side chain contract prepare lock out",
	Action:    prepareLockout,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockout",
			Value: "spectrum",
		},
		cli.Int64Flag{
			Name:  "amount",
			Usage: "amount of side chain token which you want to lockout, example amount=1 means 1 wei",
		},
		cli.Uint64Flag{
			Name:  "expiration",
			Usage: "expiration of htlc",
			Value: 1000,
		},
	},
}

func prepareLockout(ctx *cli.Context) error {
	mcName := ctx.String("mcname")
	if mcName != cfg.SMC.Name {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	contract := getSCContractAddressByMCName(mcName)
	amount := ctx.Int64("amount")
	if amount == 0 {
		fmt.Println("plo must run with --amount")
		os.Exit(-1)
	}
	expiration := ctx.Uint64("expiration")
	fmt.Printf("start to prepare lockout :\n ======> [contract=%s amount=%d expiartion=%d]\n", contract.String(), amount, expiration)

	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.HecoRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	conn := client.NewSafeClient(c)

	// 2. init contract proxy
	cp, err := proxy.NewSideChainErc20TokenProxy(conn, contract)
	if err != nil {
		fmt.Println("init contract proxy err : ", err)
		os.Exit(-1)
	}
	// 3. get auth
	privateKey, err := getPrivateKey(GlobalConfig.HecoUserAddress, GlobalConfig.HecoUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call plo
	auth := bind.NewKeyedTransactor(privateKey)
	secret := utils.NewRandomHash()
	secretHash := utils.ShaSecret(secret[:])
	expiration2 := getHecoLastBlockNumber(conn) + expiration
	fmt.Printf(" ======> [secret=%s, secretHash=%s]\n", secret.String(), secretHash.String())
	GlobalConfig.RunTime = &runTime{
		MCName:     mcName,
		Secret:     secret.String(),
		SecretHash: secretHash.String(),
	}
	updateConfigFile()
	_, scToken := getSCContractProxy(mcName)
	balance, err := scToken.Contract.BalanceOf(nil, common.HexToAddress(GlobalConfig.HecoUserAddress))
	fmt.Printf(" ======> [token balance=%d]\n", balance.Uint64())
	fmt.Println("call params :")
	fmt.Println("callerAddress = ", auth.From.String())
	fmt.Println("secretHash    = ", secretHash.String())
	fmt.Println("expiration    = ", expiration2)
	fmt.Println("amount        = ", amount)
	err = cp.PrepareLockout(auth, "", secretHash, expiration2, big.NewInt(amount))
	if err != nil {
		fmt.Println("prepare lockout err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("PrepareLockout SUCCESS")
	return nil
}
