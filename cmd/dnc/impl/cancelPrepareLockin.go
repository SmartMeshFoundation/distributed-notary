package dnc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

var cpliCmd = cli.Command{
	Name:      "cancel-prepare-lock-in",
	ShortName: "cpli",
	Usage:     "call main chain contract cancel prepare lock in",
	Action:    cancelPrepareLockin,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "ethereum",
		},
	},
}

func cancelPrepareLockin(ctx *cli.Context) error {
	mcName := ctx.String("mcname")
	if mcName == events.ChainName {
		return cancelPrepareLockin4Eth(mcName)
	}
	if mcName == bitcoin.ChainName || mcName == "btc" {
		return cancelPrepareLocin4Btc()
	}
	fmt.Println("Unknown chain name : ", mcName)
	os.Exit(-1)
	return nil
}

func cancelPrepareLocin4Btc() (err error) {
	// 1. 构造钱包连接,复用BTCService
	bs, err := bitcoin.NewBTCService(GlobalConfig.BtcWalletRPCEndpoint, GlobalConfig.BtcRPCUser, GlobalConfig.BtcRPCPass, GlobalConfig.BtcWalletRPCCertFilePath)
	if err != nil {
		fmt.Println("NewBTCService err : ", err)
		os.Exit(-1)
	}
	// 2. 构造ScriptBuilder
	userAddress, notaryAddress := getBtcAddresses(bs.GetNetParam())
	secertHash := common.HexToHash(GlobalConfig.RunTime.SecretHash)
	builder := bs.GetPrepareLockInScriptBuilder(userAddress.(*btcutil.AddressPubKeyHash), notaryAddress.(*btcutil.AddressPubKeyHash), GlobalConfig.RunTime.BtcAmount, secertHash.Bytes(), GlobalConfig.RunTime.BtcExpiration)

	// 3. 构造交易
	c := bs.GetBtcRPCClient()
	var txHash chainhash.Hash
	err = chainhash.Decode(&txHash, GlobalConfig.RunTime.BtcTXHash)
	tx := buildBtcCancelPrepareLockinTx(c, &txHash, builder, userAddress)
	newTxHash, err := c.SendRawTransaction(tx, true)
	if err != nil {
		fmt.Println("SendRawTransaction err : ", err)
		os.Exit(-1)
	}
	fmt.Println("cancel tx hash : ", newTxHash.String())
	utils.PrintBTCBalanceOfAccount(c, "default")
	return
}

func cancelPrepareLockin4Eth(mcName string) error {
	contract := getMCContractAddressByMCName(mcName)
	fmt.Printf("start to cancel prepare lockin :\n ======> [contract=%s account=%s]\n", contract.String(), GlobalConfig.EthUserAddress)

	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.EthRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	conn := client.NewSafeClient(c)

	// 2. init contract proxy
	cp, err := proxy.NewLockedEthereumProxy(conn, contract)
	if err != nil {
		fmt.Println("init contract proxy err : ", err)
		os.Exit(-1)
	}
	// 3. get auth
	privateKey, err := getPrivateKey(GlobalConfig.EthUserAddress, GlobalConfig.EthUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// 4. call pli
	auth := bind.NewKeyedTransactor(privateKey)
	err = cp.CancelLockin(auth, GlobalConfig.EthUserAddress)
	if err != nil {
		fmt.Println("cancel prepare lockin err : ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("CancelPrepareLockin SUCCESS")
	return nil
}

func buildBtcCancelPrepareLockinTx(c *rpcclient.Client, txHash *chainhash.Hash, builder *bitcoin.PrepareLockInScriptBuilder, userAddress btcutil.Address) (tx *wire.MsgTx) {
	// 1. 解锁钱包
	err := c.WalletPassphrase("123", 1000)
	if err != nil {
		fmt.Println("WalletPassphrase err : ", err)
		os.Exit(-1)
	}
	// 2.
	key, err := c.DumpPrivKey(userAddress)
	if err != nil {
		fmt.Println("DumpPrivKey err : ", err)
		os.Exit(-1)
	}
	tx = wire.NewMsgTx(wire.TxVersion)
	//
	// txIn
	prevOut := wire.NewOutPoint(txHash, 0)
	txIn := wire.NewTxIn(prevOut, nil, nil)
	tx.AddTxIn(txIn)
	// txout
	pkScript, err := txscript.PayToAddrScript(userAddress)
	if err != nil {
		fmt.Println("GetTransaction err : ", err)
		os.Exit(-1)
	}
	txOut := wire.NewTxOut(int64(GlobalConfig.RunTime.BtcAmount)-1000, pkScript)
	tx.AddTxOut(txOut)
	tx.LockTime = uint32(GlobalConfig.RunTime.BtcExpiration.Uint64())
	tx.TxIn[0].Sequence = 0
	// 签名
	sigScript, err := txscript.SignatureScript(tx, 0, GlobalConfig.RunTime.BtcLockScript, txscript.SigHashAll, key.PrivKey, true)
	if err != nil {
		fmt.Println("SignatureScript err : ", err)
		os.Exit(-1)
	}
	sb := txscript.NewScriptBuilder()
	sb.AddOps(sigScript)
	sb.AddOps(builder.GetSigScriptForUser())
	sb.AddData(GlobalConfig.RunTime.BtcLockScript)
	tx.TxIn[0].SignatureScript, err = sb.Script()
	if err != nil {
		fmt.Println("Script err : ", err)
		os.Exit(-1)
	}
	return
}
