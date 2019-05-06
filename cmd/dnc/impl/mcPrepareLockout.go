package dnc

import (
	"errors"
	"net/http"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
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
	mcName := ctx.String("mcname")
	if mcName == cfg.ETH.Name {
		return mcPrepareLockout4Eth(mcName)
	}
	if mcName == cfg.BTC.Name {
		return mcPrepareLockout4Btc(mcName)
	}
	return errors.New("unknown chain name")
}
func mcPrepareLockout4Btc(mcName string) (err error) {
	scTokenInfo := getSCTokenByMCName(mcName)
	if scTokenInfo == nil {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call plo first")
		os.Exit(-1)
	}
	url := GlobalConfig.NotaryHost + "/api/1/user/mcpreparelockout/" + scTokenInfo.SCToken.String()
	req := &userapi.MCPrepareLockoutRequest{
		BaseReq:              api.NewBaseReq(userapi.APIUserNameMCPrepareLockout),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:   api.NewBaseReqWithSCToken(scTokenInfo.SCToken),
		BaseReqWithSignature: api.NewBaseReqWithSignature(),
		SecretHash:           common.HexToHash(GlobalConfig.RunTime.SecretHash),
		SCUserAddress:        common.HexToAddress(GlobalConfig.SmcUserAddress),
	}
	// 1. 构造钱包连接,复用BTCService
	bs, err := bitcoin.NewBTCService(GlobalConfig.BtcWalletRPCEndpoint, GlobalConfig.BtcRPCUser, GlobalConfig.BtcRPCPass, GlobalConfig.BtcWalletRPCCertFilePath)
	if err != nil {
		fmt.Println("NewBTCService err : ", err)
		os.Exit(-1)
	}
	c := bs.GetBtcRPCClient()
	// 2. 解锁钱包
	err = c.WalletPassphrase("123", 1000)
	if err != nil {
		fmt.Println("WalletPassphrase err : ", err)
		os.Exit(-1)
	}
	// 3. 获取双方地址,并导出私钥
	userAddress, _ := getBtcAddresses(bs.GetNetParam())
	userWIF, err := c.DumpPrivKey(userAddress)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 签名请求
	req.Sign(req, userWIF.PrivKey.ToECDSA())
	payload := utils.ToJSONString(req)
	// 调用MCPrepareLockout
	var resp api.BaseResponse
	err = call(http.MethodPost, url, payload, &resp)
	if err != nil {
		fmt.Printf("call %s with payload=%s err :%s", url, payload, err.Error())
		os.Exit(-1)
	}
	fmt.Println("MCPrepareLockout SUCCESS")
	return
}

func mcPrepareLockout4Eth(mcName string) (err error) {
	scTokenInfo := getSCTokenByMCName(mcName)
	if scTokenInfo == nil {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call plo first")
		os.Exit(-1)
	}
	url := GlobalConfig.NotaryHost + "/api/1/user/mcpreparelockout/" + scTokenInfo.SCToken.String()
	req := &userapi.MCPrepareLockoutRequest{
		BaseReq:              api.NewBaseReq(userapi.APIUserNameMCPrepareLockout),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:   api.NewBaseReqWithSCToken(scTokenInfo.SCToken),
		BaseReqWithSignature: api.NewBaseReqWithSignature(),
		SecretHash:           common.HexToHash(GlobalConfig.RunTime.SecretHash),
		SCUserAddress:        common.HexToAddress(GlobalConfig.SmcUserAddress),
	}
	privateKey, err := getPrivateKey(GlobalConfig.EthUserAddress, GlobalConfig.EthUserPassword)
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
	fmt.Println("MCPrepareLockout SUCCESS")
	// 记录数据方便Lockout
	fmt.Println(utils.ToJSONStringFormat(resp))
	err = resp.ParseData(GlobalConfig.RunTime.LockoutInfo)
	return
}
