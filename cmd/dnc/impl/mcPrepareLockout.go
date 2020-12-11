package dnc

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"

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
}

func mcPrepareLockout(ctx *cli.Context) (err error) {
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call pli first")
		os.Exit(-1)
	}
	mcName := GlobalConfig.RunTime.MCName
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
	url := GlobalConfig.NotaryHost + "/api/1/user/mcpreparelockout2/" + scTokenInfo.SCToken.String()
	req := &userapi.MCPrepareLockoutRequest2{
		BaseReq:             api.NewBaseReq(userapi.APIUserNameMCPrepareLockout),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:  api.NewBaseReqWithSCToken(scTokenInfo.SCToken),
		SecretHash:          common.HexToHash(GlobalConfig.RunTime.SecretHash),
		SCUserAddress:       common.HexToAddress(GlobalConfig.SmcUserAddress),
	}
	key, err := getPrivateKey(GlobalConfig.EthUserAddress, GlobalConfig.EthUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	pubkey := crypto.CompressPubkey(&key.PublicKey) //[]byte类型
	req.Signer = hex.EncodeToString(pubkey)
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("req=%s\n", string(data))
	digest := crypto.Keccak256([]byte(data))
	//digest=f6b80c169ad021cbea0a8d225872bd56b8d41e15daec5173630734ea431edb48
	fmt.Printf("digest=%s", hex.EncodeToString(digest))
	signaturee, err := crypto.Sign(digest, key)
	if err != nil {
		panic(err)
		return
	}
	req.Signature = hex.EncodeToString(signaturee)

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
