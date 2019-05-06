package dnc

import (
	"net/http"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var scpliCmd = cli.Command{
	Name:      "side-chain-prepare-lock-in",
	ShortName: "scpli",
	Usage:     "call SCPrepareLockin API of notary",
	Action:    scPrepareLockin,
	Flags:     []cli.Flag{},
}

func scPrepareLockin(ctx *cli.Context) (err error) {
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	mcName := GlobalConfig.RunTime.MCName
	if mcName == cfg.ETH.Name {
		scPrepareLockIn4Eth(mcName)
	} else if mcName == cfg.BTC.Name {
		scPrepareLockIn4Btc(mcName)
	}
	return
}

func scPrepareLockIn4Btc(mcName string) {
	scTokenInfo := getSCTokenByMCName(mcName)
	if scTokenInfo == nil {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	url := GlobalConfig.NotaryHost + "/api/1/user/scpreparelockin/" + scTokenInfo.SCToken.String()
	var mcTXHash chainhash.Hash
	err := chainhash.Decode(&mcTXHash, GlobalConfig.RunTime.BtcTXHash)
	if err != nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	req := &userapi.SCPrepareLockinRequest{
		BaseReq:              api.NewBaseReq(userapi.APIUserNameSCPrepareLockin),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:   api.NewBaseReqWithSCToken(scTokenInfo.SCToken),
		BaseReqWithSignature: api.NewBaseReqWithSignature(),
		SecretHash:           common.HexToHash(GlobalConfig.RunTime.SecretHash),
		MCUserAddress:        GlobalConfig.RunTime.BtcUserAddressBytes,
		//SCUserAddress:        common.HexToAddress(GlobalConfig.SmcUserAddress),
		MCTXHash:       mcTXHash,
		MCExpiration:   GlobalConfig.RunTime.BtcExpiration,
		MCLockedAmount: GlobalConfig.RunTime.BtcAmount,
	}
	privateKey, err := getPrivateKey(GlobalConfig.SmcUserAddress, GlobalConfig.SmcUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	req.Sign(req, privateKey)
	payload := utils.ToJSONString(req)
	var resp api.BaseResponse
	err = call(http.MethodPost, url, payload, &resp)
	if err != nil {
		fmt.Printf("call %s with payload=%s err :%s", url, payload, err.Error())
		os.Exit(-1)
	}
	fmt.Println("SCPrepareLockin SUCCESS")
	return
}

func scPrepareLockIn4Eth(mcName string) {
	scTokenInfo := getSCTokenByMCName(mcName)
	if scTokenInfo == nil {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	url := GlobalConfig.NotaryHost + "/api/1/user/scpreparelockin/" + scTokenInfo.SCToken.String()
	req := &userapi.SCPrepareLockinRequest{
		BaseReq:              api.NewBaseReq(userapi.APIUserNameSCPrepareLockin),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:   api.NewBaseReqWithSCToken(scTokenInfo.SCToken),
		BaseReqWithSignature: api.NewBaseReqWithSignature(),
		SecretHash:           common.HexToHash(GlobalConfig.RunTime.SecretHash),
		MCUserAddress:        common.HexToAddress(GlobalConfig.EthUserAddress).Bytes(),
		//SCUserAddress:        common.HexToAddress(GlobalConfig.SmcUserAddress),
	}
	privateKey, err := getPrivateKey(GlobalConfig.SmcUserAddress, GlobalConfig.SmcUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	req.Sign(req, privateKey)
	payload := utils.ToJSONString(req)
	var resp api.BaseResponse
	err = call(http.MethodPost, url, payload, &resp)
	if err != nil {
		fmt.Printf("call %s with payload=%s err :%s", url, payload, err.Error())
		os.Exit(-1)
	}
	fmt.Println("SCPrepareLockin SUCCESS")
	return
}
