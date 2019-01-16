package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var queryCmd = cli.Command{
	Name:      "query",
	ShortName: "q",
	Usage:     "query lockin/lockout info on sc/mc",
	Action:    queryContract,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockout",
			Value: "ethereum",
		},
		cli.BoolFlag{
			Name:  "mcli",
			Usage: "query mc lockin info if true",
		},
		cli.BoolFlag{
			Name:  "mclo",
			Usage: "query mc lockout info if true",
		},
		cli.BoolFlag{
			Name:  "scli",
			Usage: "query sc lockin info if true",
		},
		cli.BoolFlag{
			Name:  "sclo",
			Usage: "query sc lockout info if true",
		},
		cli.BoolFlag{
			Name:  "all",
			Usage: "query all above",
		},
	},
}

func queryContract(ctx *cli.Context) {
	if globalConfig.SCTokenList == nil {
		fmt.Println("must run dnc config refresh first")
		os.Exit(-1)
	}
	mcName := ctx.String("mcname")
	if ctx.Bool("mcli") || ctx.Bool("all") {
		cp := getMCContractProxy(mcName)
		sh, e, a, err := cp.QueryLockin(globalConfig.EthUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of mc lockin:")
		fmt.Println("account    = ", globalConfig.EthUserAddress)
		fmt.Println("secretHash = ", sh.String())
		fmt.Println("expiration = ", e)
		fmt.Println("amount     = ", a)
	}
	if ctx.Bool("mclo") || ctx.Bool("all") {
		cp := getMCContractProxy(mcName)
		sh, e, a, err := cp.QueryLockout(globalConfig.EthUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of mc lockout:")
		fmt.Println("account    = ", globalConfig.EthUserAddress)
		fmt.Println("secretHash = ", sh.String())
		fmt.Println("expiration = ", e)
		fmt.Println("amount     = ", a)
	}
	if ctx.Bool("scli") || ctx.Bool("all") {
		cp := getSCContractProxy(mcName)
		sh, e, a, err := cp.QueryLockin(globalConfig.SmcUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of sc lockin:")
		fmt.Println("account    = ", globalConfig.SmcUserAddress)
		fmt.Println("secretHash = ", sh.String())
		fmt.Println("expiration = ", e)
		fmt.Println("amount     = ", a)
	}
	if ctx.Bool("sclo") || ctx.Bool("all") {
		cp := getSCContractProxy(mcName)
		sh, e, a, err := cp.QueryLockout(globalConfig.SmcUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of sc lockout:")
		fmt.Println("account    = ", globalConfig.SmcUserAddress)
		fmt.Println("secretHash = ", sh.String())
		fmt.Println("expiration = ", e)
		fmt.Println("amount     = ", a)
	}
}
