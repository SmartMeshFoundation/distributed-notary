package dnc

import (
	"fmt"
	"os"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
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
	if GlobalConfig.SCTokenList == nil {
		fmt.Println("must run dnc config refresh and get a SCToken first")
		os.Exit(-1)
	}
	mcName := ctx.String("mcname")
	//scToken := getSCTokenByMCName(mcName)
	// 区块高度信息查询
	fmt.Println("\n===> MC/SC Lasted BlockNumber info :")
	//mconn, mcp := getMCContractProxy(mcName)
	_, scp := getSCContractProxy(mcName)
	// 主测链账户余额查询
	fmt.Println("\n===> MC/SC User account info :")
	//mcUserBalance, err := mconn.BalanceAt(context.Background(), common.HexToAddress(GlobalConfig.EthUserAddress), nil)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	//fmt.Printf("[MC]user %s account balance : %d\n", utils.APex(common.HexToAddress(GlobalConfig.EthUserAddress)), wei2Eth(mcUserBalance))
	scUserTokenBalance, err := scp.Contract.BalanceOf(nil, common.HexToAddress(GlobalConfig.SmcUserAddress))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("[SC]user %s sctoken balance : %d\n", utils.APex(common.HexToAddress(GlobalConfig.SmcUserAddress)), wei2Eth(scUserTokenBalance))
	// 主侧链合约余额查询
	//fmt.Println("\n===> MC/SC Contract account info :")
	//mBalance, err := mconn.BalanceAt(context.Background(), scToken.MCLockedContractAddress, nil)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	//fmt.Printf("[MC]contract %s account balance : %d\n", utils.APex(scToken.MCLockedContractAddress), wei2Eth(mBalance))
	//scSupply, err := scp.Contract.TotalSupply(nil)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	//fmt.Printf("[SC]contract %s total supply : %d\n", utils.APex(scToken.SCToken), wei2Eth(scSupply))
	// 主侧链合约锁定查询
	fmt.Println("\n===> MC/SC Contract data info :")
	//if ctx.Bool("mcli") || ctx.Bool("all") {
	//	sh, e, a, err := mcp.QueryLockin(GlobalConfig.EthUserAddress)
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(-1)
	//	}
	//	fmt.Printf("[MC]data of lockin  : ")
	//	if sh != utils.EmptyHash {
	//		fmt.Println("\n\t account    = ", GlobalConfig.EthUserAddress)
	//		fmt.Println("\t secretHash = ", sh.String())
	//		fmt.Println("\t expiration = ", e)
	//		fmt.Println("\t amount     = ", a)
	//	} else {
	//		fmt.Println("Empty")
	//	}
	//}
	if ctx.Bool("scli") || ctx.Bool("all") {
		sh, e, a, err := scp.QueryLockin(GlobalConfig.SmcUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Printf("[SC]data of lockin  : ")
		if sh != utils.EmptyHash {
			fmt.Println("\n\t account    = ", GlobalConfig.SmcUserAddress)
			fmt.Println("\t secretHash = ", sh.String())
			fmt.Println("\t expiration = ", e)
			fmt.Println("\t amount     = ", a)
		} else {
			fmt.Println("Empty")
		}
	}
	//if ctx.Bool("mclo") || ctx.Bool("all") {
	//	sh, e, a, err := mcp.QueryLockout(GlobalConfig.EthUserAddress)
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(-1)
	//	}
	//	fmt.Printf("[MC]data of lockout : ")
	//	if sh != utils.EmptyHash {
	//		fmt.Println("\n\t account    = ", GlobalConfig.EthUserAddress)
	//		fmt.Println("\t secretHash = ", sh.String())
	//		fmt.Println("\t expiration = ", e)
	//		fmt.Println("\t amount     = ", a)
	//	} else {
	//		fmt.Println("Empty")
	//	}
	//}
	if ctx.Bool("sclo") || ctx.Bool("all") {
		sh, e, a, err := scp.QueryLockout(GlobalConfig.SmcUserAddress)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Printf("[SC]data of lockout : ")
		if sh != utils.EmptyHash {
			fmt.Println("\n\t account    = ", GlobalConfig.SmcUserAddress)
			fmt.Println("\t secretHash = ", sh.String())
			fmt.Println("\t expiration = ", e)
			fmt.Println("\t amount     = ", a)
		} else {
			fmt.Println("Empty")
		}
	}
}
