package dnc

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/params"

	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/client"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/urfave/cli"
)

var transferCmd = cli.Command{
	Name:      "transfer",
	ShortName: "t",
	Usage:     "transfer",
	Action:    transferForDebug,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to transfer ",
			Value: "ethereum",
		},
		cli.StringFlag{
			Name:  "from",
			Usage: "which acount to transfer from,default is account in configuration",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "which account to transfer to",
		},
		cli.StringFlag{
			Name:  "pass",
			Usage: "from's password,default is the pass in configuration",
		},
		cli.Int64Flag{
			Name:  "amount",
			Usage: "transfer amount, the unit is 1e15 wei",
		},
	},
}

func transferForDebug(ctx *cli.Context) {
	var err error
	mcName := ctx.String("mcname")
	amount := ctx.Int64("amount")
	if amount <= 0 {
		fmt.Println("amount must be int64 positive number ")
		os.Exit(-1)
	}
	fromAccountStr := ctx.String("from")
	fromPass := ctx.String("pass")
	var endpoint string
	if len(fromAccountStr) <= 0 {
		//从配置文件中获取
		if mcName == "ethereum" {
			fromAccountStr = GlobalConfig.EthUserAddress
			endpoint = GlobalConfig.EthRPCEndpoint
			if len(fromPass) <= 0 {
				fromPass = GlobalConfig.EthUserPassword
			}
		} else if mcName == "spectrum" {
			fromAccountStr = GlobalConfig.SmcUserAddress
			endpoint = GlobalConfig.SmcRPCEndpoint
			if len(fromPass) <= 0 {
				fromPass = GlobalConfig.SmcUserPassword
			}
		} else if mcName == "bitcoin" {
			//fromAccountStr = GlobalConfig.BtcUserAddress
			//if len(fromPass) <= 0 {
			//	fromPass = "123" //这个怎么不配置一个呢?
			//}
			fmt.Println("bitcoin not support")
			os.Exit(-1)
		} else {
			fmt.Printf("chain error, must be one of spectrum,ethereum,bitcoin:%s\n", mcName)
			os.Exit(-1)
		}
	}
	toAccountStr := ctx.String("to")
	if len(toAccountStr) <= 0 {
		fmt.Println("must specify to account")
		os.Exit(-1)
	}
	toAccount := common.HexToAddress(toAccountStr)

	privateKey, err := getPrivateKey(fromAccountStr, fromPass)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	c, err := ethclient.DialContext(context.Background(), endpoint)
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	conn := client.NewSafeClient(c)
	txCtx := context.Background()
	auth := bind.NewKeyedTransactor(privateKey)
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	var currentNonce uint64
	currentNonce, err = conn.NonceAt(txCtx, fromAddr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	amountBig := big.NewInt(amount)
	amountBig.Mul(amountBig, big.NewInt(params.Finney))
	msg := ethereum.CallMsg{From: fromAddr, To: &toAccount, Value: amountBig, Data: nil}
	gasLimit, err := conn.EstimateGas(txCtx, msg)
	if err != nil {
		fmt.Printf("failed to estimate gas needed: %v\n", err)
		os.Exit(-1)
	}
	gasPrice, err := conn.SuggestGasPrice(txCtx)
	if err != nil {
		fmt.Printf("failed to suggest gas price: %v\n", err)
		os.Exit(-1)
	}
	chainID, err := conn.NetworkID(txCtx)
	if err != nil {
		fmt.Printf("failed to get networkID : %v\n", err)
		os.Exit(-1)
	}
	rawTx := types.NewTransaction(currentNonce, toAccount, amountBig, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(types.NewEIP155Signer(chainID), auth.From, rawTx)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err = conn.SendTransaction(txCtx, signedTx); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("wait tx:\n%s", signedTx.String())
	_, err = bind.WaitMined(txCtx, conn, signedTx)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return

}
