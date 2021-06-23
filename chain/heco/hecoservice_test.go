package heco

import (
	"testing"

	"time"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestChain(t *testing.T) {
	// params
	HecoHost := "http://106.52.171.12:8545"
	var hecoContract1Address common.Address
	hecoContract1Address = common.HexToAddress("0x0f75Cc3e01d6802bca296094cEcdBb88fc50e0a6")
	// 1. 创建service
	heco, err := NewHECOService(HecoHost, hecoContract1Address)
	if err != nil {
		t.Error(err)
		return
	}
	// 2. 注册需要监听的合约
	//smc.RegisterEventListenContract(spectrumContract1Address)
	//smc.UnRegisterEventListenContract(spectrumContract1Address)
	// 3. 启动service.listener
	heco.StartEventListener()
	go func() {
		for {
			e := <-heco.GetEventChan()
			fmt.Println("收到事件:\n", utils.ToJSONStringFormat(e))
		}
	}()
	proxy := heco.GetProxyByTokenAddress(hecoContract1Address)
	name, err := proxy.Contract.Name(nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("name : ", name)
	// end
	time.Sleep(30 * time.Second)
	heco.StopEventListener()
}
