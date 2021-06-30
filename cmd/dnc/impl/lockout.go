package dnc

import (
	"errors"
	"fmt"
	"net/http"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var locmd = cli.Command{
	Name:      "lock-out",
	ShortName: "lo",
	Usage:     "call main chain contract lock out",
	Action:    lockout,
}

func lockout(ctx *cli.Context) error {
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	mcName := GlobalConfig.RunTime.MCName
	fmt.Printf("start to lockout :\n ======> [chain=%s ]\n", mcName)
	if mcName == cfg.SMC.Name {
		return lockoutOnSpectrum(mcName)
	}
	return errors.New("unknown chain name")
}

func lockoutOnSpectrum(mcName string) error {
	// 1. get proxy
	_, cp := getMCContractProxy(mcName)
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call plo first")
		os.Exit(-1)
	}
	secret := common.HexToHash(GlobalConfig.RunTime.Secret)
	fmt.Printf("start to lockout :\n ======> [account=%s secret=%s secretHash=%s]\n", GlobalConfig.HecoUserAddress, secret.String(), utils.ShaSecret(secret[:]).String())

	// 3. get auth
	privateKey, err := getPrivateKey(GlobalConfig.SmcUserAddress, GlobalConfig.SmcUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call li
	auth := bind.NewKeyedTransactor(privateKey)

	err = cp.Lockout(auth, GlobalConfig.HecoUserAddress, secret)
	if err != nil {
		fmt.Println("lockout err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("Lockout SUCCESS")
	return nil
}

// getLockoutInfo :
func getLockoutInfo(scTokenAddressHex string, secretHash string) (lockoutInfo *models.LockoutInfo, err error) {
	if GlobalConfig.NotaryHost == "" {
		err = fmt.Errorf("must set globalConfig.NotaryHost first")
		fmt.Println(err)
		return
	}
	var resp api.BaseResponse
	url := GlobalConfig.NotaryHost + "/api/1/user/lockout/" + scTokenAddressHex + "/" + secretHash
	err = call(http.MethodGet, url, "", &resp)
	if err != nil {
		err = fmt.Errorf("call %s err : %s", url, err.Error())
		fmt.Println(err)
		return
	}
	lockoutInfo = &models.LockoutInfo{}
	err = resp.ParseData(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("parse data err : %s", err.Error())
		fmt.Println(err)
	}
	return
}
