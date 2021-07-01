package dnc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/proxy"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

var cpliCmd = cli.Command{
	Name:      "cancel-prepare-lock-in",
	ShortName: "cpli",
	Usage:     "call main chain contract cancel prepare lock in",
	Action:    cancelPrepareLockin,
}

func cancelPrepareLockin(ctx *cli.Context) error {
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	mcName := GlobalConfig.RunTime.MCName
	if mcName == cfg.SMC.Name {
		return cancelPrepareLockin4SMT(mcName)
	}
	fmt.Println("Unknown chain name : ", mcName)
	os.Exit(-1)
	return nil
}

func cancelPrepareLockin4SMT(mcName string) error {
	contract := getMCContractAddressByMCName(mcName)
	fmt.Printf("start to cancel prepare lockin :\n ======> [contract=%s account=%s]\n", contract.String(), GlobalConfig.SmcUserAddress)

	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.SmcRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to spectrum chain fail : ", err)
		os.Exit(-1)
	}
	conn := client.NewSafeClient(c)

	// 2. init contract proxy
	cp, err := proxy.NewLockedSpectrumProxy(conn, contract)
	if err != nil {
		fmt.Println("cancelPrepareLockin4SMT init contract proxy err : ", err)
		os.Exit(-1)
	}
	// 3. get auth
	privateKey, err := getPrivateKey(GlobalConfig.SmcUserAddress, GlobalConfig.SmcUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call pli
	auth := bind.NewKeyedTransactor(privateKey)
	err = cp.CancelLockin(auth, GlobalConfig.SmcUserAddress)
	if err != nil {
		fmt.Println("cancel prepare lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("CancelPrepareLockin SUCCESS")
	return nil
}
