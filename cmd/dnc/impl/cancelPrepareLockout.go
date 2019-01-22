package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/proxy"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

var cploCmd = cli.Command{
	Name:      "cancel-prepare-lock-out",
	ShortName: "cplo",
	Usage:     "call side chain contract cancel prepare lock out",
	Action:    cancelPrepareLockout,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockout",
			Value: "ethereum",
		},
	},
}

func cancelPrepareLockout(ctx *cli.Context) error {
	contract := getSCContractAddressByMCName(ctx.String("mcname"))
	fmt.Printf("start to cancel prepare lockout :\n ======> [contract=%s account=%s]\n", contract.String(), globalConfig.EthUserAddress)

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
	cp, err := proxy.NewSideChainErc20TokenProxy(conn, contract)
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
	// 4. call pli
	auth := bind.NewKeyedTransactor(privateKey)
	err = cp.CancelLockout(auth, globalConfig.SmcUserAddress)
	if err != nil {
		fmt.Println("cancel prepare lockout err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("CancelPrepareLockout SUCCESS")
	return nil
}
