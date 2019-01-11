package ethereum

import (
	"testing"

	"time"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestChain(t *testing.T) {
	// params
	ethHost := "http://127.0.0.1:9001"
	var contractAddress common.Address
	contractAddress = common.HexToAddress("0x720bF7a52fDb3f656E0E653E09C4e57DC1e655eE")
	// 1. 创建service
	eth, err := NewETHService(ethHost, contractAddress)
	if err != nil {
		t.Error(err)
		return
	}
	// 2. 注册需要监听的合约
	//smc.RegisterEventListenContract(spectrumContract1Address)
	//smc.UnRegisterEventListenContract(spectrumContract1Address)
	// 3. 启动service.listener
	eth.StartEventListener()
	go func() {
		for {
			e := <-eth.GetEventChan()
			fmt.Println("收到事件:\n", utils.ToJSONStringFormat(e))
		}
	}()

	proxy := eth.GetProxyByLockedEthereumAddress(contractAddress)
	name, err := proxy.Contract.Name(nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("name : ", name)
	// end
	time.Sleep(30 * time.Second)
	eth.StopEventListener()
}
