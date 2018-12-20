package spectrum

import (
	"testing"

	"time"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestChain(t *testing.T) {
	// params
	spectrumHost := "http://127.0.0.1:9001"
	var spectrumContract1Address common.Address
	spectrumContract1Address = common.HexToAddress("0x63d6014616112d528A9cdc5e4A043267932E659d")
	// 1. 创建service
	smc, _ := NewSMCService(spectrumHost, 0, spectrumContract1Address)
	// 2. 注册需要监听的合约
	//smc.RegisterEventListenContract(spectrumContract1Address)
	//smc.UnRegisterEventListenContract(spectrumContract1Address)
	// 3. 启动service.listener
	smc.StartEventListener()
	go func() {
		for {
			e := <-smc.GetEventChan()
			fmt.Println("收到事件:\n", utils.ToJsonStringFormat(e))
		}
	}()
	time.Sleep(30 * time.Second)
	smc.StopEventListener()
}
