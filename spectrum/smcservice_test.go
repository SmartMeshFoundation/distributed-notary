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
	spectrumHost := "http://192.168.124.13:28545"
	var spectrumContract1Address common.Address
	spectrumContract1Address = common.HexToAddress("0x0f75Cc3e01d6802bca296094cEcdBb88fc50e0a6")
	// 1. 创建service
	smc, err := NewSMCService(spectrumHost, 0, spectrumContract1Address)
	if err != nil {
		t.Error(err)
		return
	}
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
