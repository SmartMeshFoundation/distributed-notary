package main

import (
	"fmt"
	"os"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
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
	mcp := getMCContractProxy(mcName)
	scp := getSCContractProxy(mcName)
	if ctx.Bool("mcli") || ctx.Bool("all") {
		sh, e, a, err := mcp.QueryLockin(globalConfig.EthUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of mc lockin:")
		if sh != utils.EmptyHash {
			fmt.Println("\t account    = ", globalConfig.EthUserAddress)
			fmt.Println("\t secretHash = ", sh.String())
			fmt.Println("\t expiration = ", e)
			fmt.Println("\t amount     = ", a)
		} else {
			fmt.Println("Empty")
		}
	}
	if ctx.Bool("scli") || ctx.Bool("all") {
		sh, e, a, err := scp.QueryLockin(globalConfig.SmcUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of sc lockin:")
		if sh != utils.EmptyHash {
			fmt.Println("\t account    = ", globalConfig.SmcUserAddress)
			fmt.Println("\t secretHash = ", sh.String())
			fmt.Println("\t expiration = ", e)
			fmt.Println("\t amount     = ", a)
		} else {
			fmt.Println("Empty")
		}
	}
	if ctx.Bool("mclo") || ctx.Bool("all") {
		sh, e, a, err := mcp.QueryLockout(globalConfig.EthUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of mc lockout:")
		if sh != utils.EmptyHash {
			fmt.Println("\t account    = ", globalConfig.EthUserAddress)
			fmt.Println("\t secretHash = ", sh.String())
			fmt.Println("\t expiration = ", e)
			fmt.Println("\t amount     = ", a)
		} else {
			fmt.Println("Empty")
		}
	}
	if ctx.Bool("sclo") || ctx.Bool("all") {
		sh, e, a, err := scp.QueryLockout(globalConfig.SmcUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println("===> contract data of sc lockout:")
		if sh != utils.EmptyHash {
			fmt.Println("\t account    = ", globalConfig.SmcUserAddress)
			fmt.Println("\t secretHash = ", sh.String())
			fmt.Println("\t expiration = ", e)
			fmt.Println("\t amount     = ", a)
		} else {
			fmt.Println("Empty")
		}
	}
}
