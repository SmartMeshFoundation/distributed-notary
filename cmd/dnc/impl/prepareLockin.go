package dnc

import (
	"context"
	"fmt"
	"time"

	"os"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kataras/go-errors"
	"github.com/urfave/cli"
)

var pliCmd = cli.Command{
	Name:      "prepare-lock-in",
	ShortName: "pli",
	Usage:     "call main chain contract prepare lock in",
	Action:    prepareLockin,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "spectrum",
		},
		cli.Int64Flag{
			Name:  "amount",
			Usage: "amount of side chain token which you want to lockin, example amount=1 means 1 wei",
		},
		cli.Uint64Flag{
			Name:  "expiration",
			Usage: "expiration of htlc",
			Value: 900,
		},
	},
}

func prepareLockin(ctx *cli.Context) error {
	mcName := ctx.String("mcname")
	amountStr := ctx.String("amount")
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		fmt.Println("pli must run with --amount,amount format error")
		os.Exit(-1)
	}
	if amount == big.NewInt(0) {
		fmt.Println("pli must run with --amount")
		os.Exit(-1)
	}
	if mcName != cfg.SMC.Name {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	expiration := ctx.Uint64("expiration")
	fmt.Printf("start to prepare lockin :\n ======> [chain=%s amount=%d expiartion=%d]\n", mcName, amount, expiration)
	if mcName == cfg.SMC.Name {
		return prepareLockinOnSpectrum(mcName, amount, expiration)
	}
	return errors.New("unknown chain name")
}

func prepareLockinOnSpectrum(mcName string, amount *big.Int, expiration uint64) (err error) {
	contract := getMCContractAddressByMCName(mcName)
	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.SmcRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to Spectrum fail : ", err)
		os.Exit(-1)
	}
	conn := client.NewSafeClient(c)

	// 2. init contract proxy
	cp, err := proxy.NewLockedSpectrumProxy(conn, contract)
	if err != nil {
		fmt.Println("prepareLockinOnSpectrum init contract proxy err : ", err)
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
	secret := utils.NewRandomHash()
	secretHash := utils.ShaSecret(secret[:])
	expiration2 := getSmcLastBlockNumber(conn) + expiration
	fmt.Printf(" ======> [secret=%s, secretHash=%s]\n", secret.String(), secretHash.String())
	GlobalConfig.RunTime = &runTime{
		MCName:     cfg.SMC.Name,
		Secret:     secret.String(),
		SecretHash: secretHash.String(),
	}
	updateConfigFile()
	err = cp.PrepareLockin(auth, "", secretHash, expiration2, amount)
	if err != nil {
		fmt.Println("prepare lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("PrepareLockin on Spectrum SUCCESS")
	return nil
}
