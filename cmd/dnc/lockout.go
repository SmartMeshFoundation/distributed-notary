package main

import (
	"fmt"
	"os"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var locmd = cli.Command{
	Name:      "lock-out",
	ShortName: "lo",
	Usage:     "call main chain contract lock out",
	Action:    lockout,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "ethereum",
		},
	},
}

func lockout(ctx *cli.Context) error {
	// 1. get proxy
	_, cp := getMCContractProxy(ctx.String("mcname"))
	if globalConfig.RunTime == nil {
		fmt.Println("must call plo first")
		os.Exit(-1)
	}
	secret := common.HexToHash(globalConfig.RunTime.Secret)
	fmt.Printf("start to lockout :\n ======> [account=%s secret=%s secretHash=%s]\n", globalConfig.SmcUserAddress, secret.String(), utils.ShaSecret(secret[:]).String())

	// 3. get auth
	privateKey, err := getPrivateKey(globalConfig.EthUserAddress, globalConfig.EthUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call li
	auth := bind.NewKeyedTransactor(privateKey)

	err = cp.Lockout(auth, globalConfig.SmcUserAddress, secret)
	if err != nil {
		fmt.Println("lockout err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("Lockout SUCCESS")
	return nil
}
