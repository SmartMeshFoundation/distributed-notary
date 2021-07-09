package dnc

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	utilss "github.com/nkbai/goutils"
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
	if mcName == cfg.SMC.Name {
		scPrepareLockIn4Eth(mcName)
	}
	return
}

func scPrepareLockIn4Eth(mcName string) {
	scTokenInfo := getSCTokenByMCName(mcName)
	if scTokenInfo == nil {
		fmt.Println("wrong mcname")
		os.Exit(-1)
	}
	fmt.Println(fmt.Sprintf("scTokenInfo=%s", utilss.StringInterface(scTokenInfo, 3)))
	fmt.Println(mcName)
	url := GlobalConfig.NotaryHost + "/api/1/user/scpreparelockin2/" + scTokenInfo.SCToken.String()
	req := &userapi.SCPrepareLockinRequest2{
		BaseReq:             api.NewBaseReq(userapi.APIUserNameSCPrepareLockin),
		BaseReqWithResponse: api.NewBaseReqWithResponse(),
		BaseReqWithSCToken:  api.NewBaseReqWithSCToken(scTokenInfo.SCToken),

		SecretHash:    common.HexToHash(GlobalConfig.RunTime.SecretHash),
		MCUserAddress: common.HexToAddress(GlobalConfig.SmcUserAddress),
		//SCUserAddress:        common.HexToAddress(GlobalConfig.SmcUserAddress),
	}
	privateKey, err := getPrivateKey(GlobalConfig.HecoUserAddress, GlobalConfig.HecoUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	key := privateKey
	pubkey := crypto.CompressPubkey(&key.PublicKey) //[]byte类型
	req.Signer = hex.EncodeToString(pubkey)
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("req=%s\n", string(data))
	digest := crypto.Keccak256([]byte(data))
	//digest=73aa7baa7ee1416b76aa69c3605c104e1995e4762a83eecfeea61862cb08d616
	fmt.Printf("digest=%s\n", hex.EncodeToString(digest))
	signaturee, err := crypto.Sign(digest, key)
	if err != nil {
		panic(err)
		return
	}
	req.Signature = hex.EncodeToString(signaturee)
	data, err = json.Marshal(req)
	if err != nil {
		panic(err)
		return
	}
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
