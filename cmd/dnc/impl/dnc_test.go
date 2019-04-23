package dnc

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"

	"time"

	"github.com/ethereum/go-ethereum/params"
)

func TestEther(t *testing.T) {
	fmt.Println(params.Ether)
	fmt.Println(uint64(params.Ether))
	fmt.Println(int64(params.Ether))
}

func TestVerify(t *testing.T) {

	s := `{"name":"User-SCPrepareLockin","request_id":"183e","sc_token_address":"0x88148d2f9e23769a143396d6124121a13d5a7c39","signer":"0x56771612bfae7fda173cb89579cc67876e34d6e7","secret_hash":"0x6de53701dcd7916f01f11687c3d0066ebe932ea59980e25f0191e03e693679a4","mc_user_address":"0x56771612bfae7fda173cb89579cc67876e34d6e7","sc_user_address":"0x56771612bfae7fda173cb89579cc67876e34d6e7","signature":"R7m0aUkoDmulgbDTQwzTrvllD4zxZ7N3vbiMuG1IPpR8PluLeHc3KjDD9cRDA+TIKWiBUndJjX7ua2RIW33vigA="}`

	var req userapi.SCPrepareLockinRequest
	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		t.Error(err)
		return
	}
	b := req.VerifySign(&req)
	if !b {
		t.Error("should success")
	}
}

func TestLoadTxFilter(t *testing.T) {
	//ntfnHandlers := &rpcclient.NotificationHandlers{
	//	OnRelevantTxAccepted: func(transaction []byte) {
	//		fmt.Println("====> OnRelevantTxAccepted Got :")
	//		fmt.Println(common.Bytes2Hex(transaction))
	//		var msg wire.MsgTx
	//		msg.Deserialize(bytes.NewReader(transaction))
	//		fmt.Println(utils.ToJSONStringFormat(msg))
	//	},
	//	//OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*btcutil.Tx) {
	//	//	fmt.Println("OnFilteredBlockConnected new block :", height)
	//	//},
	//	OnClientConnected: func() {
	//		fmt.Println("==================== connect SUCCESS ")
	//	},
	//}
	// 0. 构造btcd连接
	bs0, err := bitcoin.NewBTCService(GlobalConfig.BtcRPCEndpoint, GlobalConfig.BtcRPCUser, GlobalConfig.BtcRPCPass, GlobalConfig.BtcRPCCertFilePath)
	if err != nil {
		fmt.Println("NewBTCService err : ", err)
		os.Exit(-1)
	}
	// 2. 获取双方地址
	//notaryAddress, err := btcutil.DecodeAddress("SRgR4wx8UyauNDWpjPSrjhbPLuEgyhJJUJ", bs0.GetNetParam())
	//if err != nil {
	//	fmt.Println("DecodeAddress err : ", err)
	//	os.Exit(-1)
	//}
	// 3. 启动监听
	//err = bs0.GetBtcRPCClient().LoadTxFilter(false, []btcutil.Address{notaryAddress}, nil)
	//if err != nil {
	//	fmt.Println("LoadTxFilter err : ", err)
	//	os.Exit(-1)
	//}
	bs0.GetBtcRPCClient().NotifyBlocks()
	c := bs0.GetBtcRPCClient()
	time.Sleep(5 * time.Second)
	c.Disconnect()
	fmt.Println("Disconnect")
	time.Sleep(5 * time.Second)
	fmt.Println("Connect")
	c.Connect(5)
	//// 1. 构造钱包连接,复用BTCService
	//bs, err := bitcoin.NewBTCService(GlobalConfig.BtcWalletRPCEndpoint, GlobalConfig.BtcRPCUser, GlobalConfig.BtcRPCPass, GlobalConfig.BtcWalletRPCCertFilePath)
	//if err != nil {
	//	fmt.Println("NewBTCService err : ", err)
	//	os.Exit(-1)
	//}
	//c := bs.GetBtcRPCClient()
	//// 3. 解锁钱包
	//err = c.WalletPassphrase("123", 1000)
	//if err != nil {
	//	fmt.Println("WalletPassphrase err : ", err)
	//	os.Exit(-1)
	//}
	//// 4. 发送交易
	//fmt.Println("发送交易...")
	//txHash, err := c.SendToAddress(notaryAddress, btcutil.Amount(100000000))
	//if err != nil {
	//	fmt.Println("SendToAddress err : ", err)
	//	os.Exit(-1)
	//}
	//fmt.Println("交易发送完成,txHash=", txHash.String())
	time.Sleep(10000 * time.Second)
}
