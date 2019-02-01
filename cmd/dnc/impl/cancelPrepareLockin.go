package dnc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/proxy"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

var cpliCmd = cli.Command{
	Name:      "cancel-prepare-lock-in",
	ShortName: "cpli",
	Usage:     "call main chain contract cancel prepare lock in",
	Action:    cancelPrepareLockin,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "ethereum",
		},
	},
}

func cancelPrepareLockin(ctx *cli.Context) error {
	contract := getMCContractAddressByMCName(ctx.String("mcname"))
	fmt.Printf("start to cancel prepare lockin :\n ======> [contract=%s account=%s]\n", contract.String(), GlobalConfig.EthUserAddress)

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
	err = cp.CancelLockin(auth, GlobalConfig.EthUserAddress)
	if err != nil {
		fmt.Println("cancel prepare lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("CancelPrepareLockin SUCCESS")
	return nil
}
