package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/client"
	proxy2 "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

var licmd = cli.Command{
	Name:      "lock-in",
	ShortName: "li",
	Usage:     "call side chain contract lock in",
	Action:    lockin,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "ethereum",
		},
	},
}

func lockin(ctx *cli.Context) error {
	contract := getSCContractAddressByMCName(ctx.String("mcname"))
	if globalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	secret := common.HexToHash(globalConfig.RunTime.Secret)
	fmt.Printf("start to lockin :\n ======> [contract=%s account=%s secret=%s secretHash=%s]\n", contract.String(), globalConfig.SmcUserAddress, secret.String(), utils.ShaSecret(secret[:]).String())

	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, globalConfig.SmcRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	conn := client.NewSafeClient(c)

	// 2. init contract proxy
	cp, err := proxy2.NewSideChainErc20TokenProxy(conn, contract)
	if err != nil {
		fmt.Println("init contract proxy err : ", err)
		os.Exit(-1)
	}
	// 3. get auth
	privateKey, err := getPrivateKey(globalConfig.SmcUserAddress, globalConfig.SmcUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call li
	auth := bind.NewKeyedTransactor(privateKey)

	sh, e, a, err := cp.QueryLockin(globalConfig.SmcUserAddress)
	fmt.Println("contract data :")
	fmt.Println("secretHash = ", sh.String())
	fmt.Println("expiration = ", e)
	fmt.Println("amount     = ", a)

	err = cp.Lockin(auth, globalConfig.SmcUserAddress, secret)
	if err != nil {
		fmt.Println("lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("Lockin SUCCESS")
	return nil
}
