package dnc

import (
	"errors"
	"fmt"
	"net/http"

	"os"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var locmd = cli.Command{
	Name:      "lock-out",
	ShortName: "lo",
	Usage:     "call main chain contract lock out",
	Action:    lockout,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "ethereum",
		},
	},
}

func lockout(ctx *cli.Context) error {
	mcName := ctx.String("mcname")
	fmt.Printf("start to lockout :\n ======> [chain=%s ]\n", mcName)
	if mcName == cfg.ETH.Name {
		return lockoutOnEthereum(mcName)
	} else if mcName == cfg.BTC.Name || mcName == "btc" {
		return lockoutOnBitcoin(mcName)
	}
	return errors.New("unknown chain name")
}
func lockoutOnBitcoin(mcName string) error {
	scToken := getSCTokenByMCName(mcName)
	lockoutInfo, err := getLockoutInfo(scToken.SCToken.String(), GlobalConfig.RunTime.SecretHash)
	if err != nil || lockoutInfo == nil {
		fmt.Println("need call mcplo first ", err)
		os.Exit(-1)
	}
	if lockoutInfo.BTCLockScriptHex == "" {
		fmt.Println("wait for notary to confirm prepare lockout, retry after 6 blocks ")
		os.Exit(-1)
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
	// 4. 构造交易
	tx := wire.NewMsgTx(wire.TxVersion)
	// txIn
	prepareLockoutTxHash, err := chainhash.NewHashFromStr(lockoutInfo.BTCPrepareLockoutTXHashHex)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	prepareLockinOutPoint := &wire.OutPoint{
		Hash:  *prepareLockoutTxHash,
		Index: lockoutInfo.BTCPrepareLockoutVout,
	}
	txIn := wire.NewTxIn(prepareLockinOutPoint, nil, nil)
	tx.AddTxIn(txIn)
	// txOut
	pkScript, err := txscript.PayToAddrScript(userAddress)
	txOut := wire.NewTxOut(lockoutInfo.Amount.Int64()-lockoutInfo.CrossFee.Int64()-1000, pkScript)
	tx.AddTxOut(txOut)
	// 5. 签名
	lockScript := common.Hex2Bytes(lockoutInfo.BTCLockScriptHex)
	sigScript, err := txscript.SignatureScript(tx, 0, lockScript, txscript.SigHashAll, userWIF.PrivKey, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 构造SignatureScript
	sb := txscript.NewScriptBuilder()
	sb.AddOps(sigScript)
	sb.AddData(common.HexToHash(GlobalConfig.RunTime.Secret).Bytes())
	sb.AddOp(txscript.OP_TRUE)
	sb.AddData(lockScript)
	tx.TxIn[0].SignatureScript, err = sb.Script()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 6. 发送交易
	txHash, err := c.SendRawTransaction(tx, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf(" ======> [txHash=%s]\n", txHash.String())
	return nil
}

func lockoutOnEthereum(mcName string) error {
	// 1. get proxy
	_, cp := getMCContractProxy(mcName)
	if GlobalConfig.RunTime == nil {
		fmt.Println("must call plo first")
		os.Exit(-1)
	}
	secret := common.HexToHash(GlobalConfig.RunTime.Secret)
	fmt.Printf("start to lockout :\n ======> [account=%s secret=%s secretHash=%s]\n", GlobalConfig.SmcUserAddress, secret.String(), utils.ShaSecret(secret[:]).String())

	// 3. get auth
	privateKey, err := getPrivateKey(GlobalConfig.EthUserAddress, GlobalConfig.EthUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call li
	auth := bind.NewKeyedTransactor(privateKey)

	err = cp.Lockout(auth, GlobalConfig.SmcUserAddress, secret)
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
