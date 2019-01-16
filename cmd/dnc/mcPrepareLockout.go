package main

import (
	"net/http"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var mcploCmd = cli.Command{
	Name:      "main-chain-prepare-lock-out",
	ShortName: "mcplo",
	Usage:     "call MCPrepareLockout API of notary",
	Action:    mcPrepareLockout,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockout",
			Value: "ethereum",
		},
	},
}

func mcPrepareLockout(ctx *cli.Context) (err error) {
	scTokenInfo := getSCTokenByMCName(ctx.String("mcname"))
	if scTokenInfo == nil {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	if globalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	url := globalConfig.NotaryHost + "/api/1/user/mcpreparelockout/" + scTokenInfo.SCToken.String()
	req := &userapi.MCPrepareLockoutRequest{
		BaseRequest:           api.NewBaseRequest(userapi.APIUserNameMCPrepareLockout),
		BaseCrossChainRequest: api.NewBaseCrossChainRequest(scTokenInfo.SCToken),
		SecretHash:            common.HexToHash(globalConfig.RunTime.SecretHash),
		MCUserAddress:         common.HexToAddress(globalConfig.EthUserAddress),
		SCUserAddress:         common.HexToAddress(globalConfig.SmcUserAddress),
	}
	privateKey, err := getPrivateKey(globalConfig.EthUserAddress, globalConfig.EthUserAddress)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	req.Sign(privateKey)
	payload := utils.ToJSONString(req)
	var resp api.BaseResponse
	err = call(http.MethodPost, url, payload, &resp)
	if err != nil {
		fmt.Printf("call %s with payload=%s err :%s", url, payload, err.Error())
		os.Exit(-1)
	}
	return
}
