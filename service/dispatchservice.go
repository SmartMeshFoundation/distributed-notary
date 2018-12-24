package service

import (
	"fmt"

	"sync"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum"
	ethereumevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum"
	spectrumevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
DispatchService :
核心调度service,管理所有chains, notaryService及所有消息分发, 单例,常驻3个线程,分别为:
	1. chainEventDispatcherLoop 事件调度线程,负责监听所有链的所有事件并分发至对应的notaryService
	2. APIRequestDispatcherLoop 用户请求调度线程
	3. NotaryMessageDispatcherLoop 节点消息调度线程

调度器不参与错误处理,仅记录日志
*/
type DispatchService struct {
	/*
		保存并维护所有接入的链,key为ChainName
	*/
	chainMap map[string]chain.Chain
	/*
		当前正在进行的公证行为
	*/
	scToken2NotaryServiceMap     map[common.Address]*NotaryService
	scToken2NotaryServiceMapLock sync.Mutex

	quitChan chan struct{}
}

// NewDispatchService :
func NewDispatchService() (ds *DispatchService, err error) {
	ds = &DispatchService{
		chainMap:                 make(map[string]chain.Chain),
		scToken2NotaryServiceMap: make(map[common.Address]*NotaryService),
		quitChan:                 make(chan struct{}),
	}
	// 应该怎么初始化???
	ds.chainMap[spectrumevents.ChainName], err = spectrum.NewSMCService("", 0)
	ds.chainMap[ethereumevents.ChainName], err = ethereum.NewETHService("", 0)
	return
}

// DispatchAPIRequest API调用调度
func (ds *DispatchService) DispatchAPIRequest() {

}

// DispatchNotaryMessage 公证人消息调度
func (ds *DispatchService) DispatchNotaryMessage() {

}

// StartChainEventDispatcher 链上事件调度
func (ds *DispatchService) StartChainEventDispatcher() {
	/*
		对每条链启动事件处理线程
	*/
	for _, v := range ds.chainMap {
		go ds.chainEventDispatcherLoop(v)
	}
}

// StopChainEventDispatcher :
func (ds *DispatchService) StopChainEventDispatcher() {
	close(ds.quitChan)
}

func (ds *DispatchService) chainEventDispatcherLoop(c chain.Chain) {
	logPrefix := fmt.Sprintf("Chain %s : ", c.GetChainName())
	log.Info(fmt.Sprintf("%s dispather start ...", logPrefix))
	for {
		select {
		case e, ok := <-c.GetEventChan():
			if !ok {
				log.Error(fmt.Sprintf("%s read event chan err ", logPrefix))
				continue
			}
			if e.GetEventName() == chain.NewBlockNumberEventName {
				// 新块事件,dispatch至所有service
				ds.dispatchEvent2All(e)
			} else {
				// 合约事件,根据SCToken调度
				ds.dispatchEvent(e)
			}
		case <-ds.quitChan:
			log.Info(fmt.Sprintf("%s dispather stop success", logPrefix))
			return
		}
	}
}

func (ds *DispatchService) dispatchEvent2All(e chain.Event) {
	logPrefix := fmt.Sprintf("Chain %s : ", e.GetChainName())
	var err error
	var needRemove bool
	for scTokenAddress, service := range ds.scToken2NotaryServiceMap {
		needRemove, err = service.OnChainEvent(e)
		if err != nil {
			log.Error(fmt.Sprintf("%s notary service err when deal event: err=%s,event:\n%s\n", logPrefix, err.Error(), utils.ToJsonStringFormat(e)))
		}
		if needRemove {
			ds.scToken2NotaryServiceMapLock.Lock()
			delete(ds.scToken2NotaryServiceMap, scTokenAddress)
			ds.scToken2NotaryServiceMapLock.Unlock()
			log.Info(fmt.Sprintf("%s remove notary service of SCToken %s", logPrefix, scTokenAddress.String()))
		}
	}
}

func (ds *DispatchService) dispatchEvent(e chain.Event) {
	logPrefix := fmt.Sprintf("Chain %s : ", e.GetChainName())
	if e.GetSCTokenAddress() == utils.EmptyAddress {
		// 主链事件,根据主链合约地址FromAddress调度,遍历,后续可优化,维护一个主链合约地址-SCToken地址的map即可
		for scTokenAddress, service := range ds.scToken2NotaryServiceMap {
			if service.GetMCContractAddress() == e.GetFromAddress() {
				// 事件业务逻辑处理
				needRemove, err := service.OnChainEvent(e)
				if err != nil {
					log.Error(fmt.Sprintf("%s notary service err when deal event: err=%s,event:\n%s\n", logPrefix, err.Error(), utils.ToJsonStringFormat(e)))
				}
				if needRemove {
					ds.scToken2NotaryServiceMapLock.Lock()
					delete(ds.scToken2NotaryServiceMap, scTokenAddress)
					ds.scToken2NotaryServiceMapLock.Unlock()
					log.Info(fmt.Sprintf("%s remove notary service of SCToken %s", logPrefix, scTokenAddress.String()))
				}
				// 每个事件应该只有一个对应service,所以这里处理完毕直接return ???
				return
			}
		}
	} else {
		// 侧链事件,直接根据SCToken地址调度
		service, ok := ds.scToken2NotaryServiceMap[e.GetSCTokenAddress()]
		if !ok {
			/*
				是否存在什么公证人收到什么事件的时候需要new一个notary的情况???
			*/
			log.Error(fmt.Sprintf("%s get event with out notary service : \n%s\n", logPrefix, utils.ToJsonStringFormat(e)))
			return
		}
		// 事件业务逻辑处理
		needRemove, err := service.OnChainEvent(e)
		if err != nil {
			log.Error(fmt.Sprintf("%s notary service err when deal event: err=%s,event:\n%s\n", logPrefix, err.Error(), utils.ToJsonStringFormat(e)))
		}
		if needRemove {
			ds.scToken2NotaryServiceMapLock.Lock()
			delete(ds.scToken2NotaryServiceMap, e.GetSCTokenAddress())
			ds.scToken2NotaryServiceMapLock.Unlock()
			log.Info(fmt.Sprintf("%s remove notary service of SCToken %s", logPrefix, e.GetSCTokenAddress().String()))
		}
	}
}
