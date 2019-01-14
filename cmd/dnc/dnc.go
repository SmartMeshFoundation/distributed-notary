package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/SmartMeshFoundation/distributed-notary/accounts"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/nkbai/log"
	"github.com/urfave/cli"
)

/*
	用户client,demo性质
*/

//Version version of this build
var Version string

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		configCmd,
		pliCmd,
		cpliCmd,
	}
	app.Action = startMain
	app.Name = "dnc"
	app.Version = Version
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(-1)
	}
}

func startMain(ctx *cli.Context) {
}

func getEthPrivateKey() (privateKey *ecdsa.PrivateKey, err error) {
	am := accounts.NewAccountManager(globalConfig.Keystore)
	privateKeyBin, err := am.GetPrivateKey(common.HexToAddress(globalConfig.EthUserAddress), globalConfig.EthUserPassword)
	if err != nil {
		log.Error("load private key err : %s", err.Error())
		return
	}
	privateKey, err = crypto.ToECDSA(privateKeyBin)
	if err != nil {
		log.Error("load private key err : %s", err.Error())
		return
	}
	return
}

func getEthLastBlockNumber(c *client.SafeEthClient) uint64 {
	h, err := c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return h.Number.Uint64()
}

func getMCContractAddressByMCName(mcName string) common.Address {
	if globalConfig.SCTokenList == nil || len(globalConfig.SCTokenList) == 0 {
		fmt.Println("must run dnc config refresh first")
		os.Exit(-1)
	}
	for _, sctoken := range globalConfig.SCTokenList {
		if sctoken.MCName == mcName {
			return sctoken.MCLockedContractAddress
		}
	}
	fmt.Printf("can not found mc contract address of %s\n", mcName)
	os.Exit(-1)
	return utils.EmptyAddress
}

func getEther(amount int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(params.Ether)), big.NewInt(amount))
}
