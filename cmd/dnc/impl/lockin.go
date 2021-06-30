package dnc

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
}

func lockin(ctx *cli.Context) error {
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	_, cp := getSCContractProxy(GlobalConfig.RunTime.MCName)
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	secret := common.HexToHash(GlobalConfig.RunTime.Secret)
	fmt.Printf("start to lockin :\n ======> [account=%s secret=%s secretHash=%s]\n", GlobalConfig.HecoUserAddress, secret.String(), utils.ShaSecret(secret[:]).String())

	// 3. get auth
	privateKey, err := getPrivateKey(GlobalConfig.HecoUserAddress, GlobalConfig.HecoUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call li
	auth := bind.NewKeyedTransactor(privateKey)

	err = cp.Lockin(auth, GlobalConfig.HecoUserAddress, secret)
	if err != nil {
		fmt.Println("lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("Lockin SUCCESS")
	return nil
}
