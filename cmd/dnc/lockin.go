package main

import (
	"fmt"
	"os"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	cp := getSCContractProxy(ctx.String("mcname"))
	if globalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	secret := common.HexToHash(globalConfig.RunTime.Secret)
	fmt.Printf("start to lockin :\n ======> [account=%s secret=%s secretHash=%s]\n", globalConfig.SmcUserAddress, secret.String(), utils.ShaSecret(secret[:]).String())

	// 3. get auth
	privateKey, err := getPrivateKey(globalConfig.SmcUserAddress, globalConfig.SmcUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call li
	auth := bind.NewKeyedTransactor(privateKey)

	err = cp.Lockin(auth, globalConfig.SmcUserAddress, secret)
	if err != nil {
		fmt.Println("lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("Lockin SUCCESS")
	return nil
}
