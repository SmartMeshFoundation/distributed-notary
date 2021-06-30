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
	if mcName == cfg.SMC.Name {
		return mcPrepareLockout4SMT(mcName)
	}
	return errors.New("unknown chain name")
}

func mcPrepareLockout4SMT(mcName string) (err error) {
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
	key, err := getPrivateKey(GlobalConfig.SmcUserAddress, GlobalConfig.SmcUserPassword)
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
