package service

import (
	"fmt"

	"sync"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum"
	ethereumevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum"
	spectrumevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
DispatchService :
核心调度service,管理所有chains, notaryService及所有消息分发, 单例,常驻3个线程,分别为:
	1. chainEventDispatcherLoop 事件调度线程,负责监听所有链的所有事件并分发至对应的notaryService
	2. APIRequestDispatcherLoop 用户请求调度线程
	3. NotaryRequestDispatcherLoop 节点消息调度线程

调度器除异常崩溃外,不参与错误处理,仅记录日志,调度器本身支持高并发,业务逻辑的线程安全及顺序性等问题,交由业务service处理
*/
type DispatchService struct {
	// restful api
	userAPI   *userapi.UserAPI
	notaryAPI *notaryapi.NotaryAPI

	// 区块链事件监听
	chainMap map[string]chain.Chain

	// 数据库
	db *models.DB

	// 杂项
	quitChan chan struct{}

	/*
		业务处理service
	*/
	systemService                    *SystemService
	notaryService                    *NotaryService
	scToken2CrossChainServiceMap     map[common.Address]*CrossChainService
	scToken2CrossChainServiceMapLock sync.Mutex
}

// NewDispatchService :
func NewDispatchService(userAPI *userapi.UserAPI, notaryAPI *notaryapi.NotaryAPI, db *models.DB) (ds *DispatchService, err error) {
	// 校验

	ds = &DispatchService{
		userAPI:                      userAPI,
		notaryAPI:                    notaryAPI,
		chainMap:                     make(map[string]chain.Chain),
		db:                           db,
		quitChan:                     make(chan struct{}),
		scToken2CrossChainServiceMap: make(map[common.Address]*CrossChainService),
	}
	/*
		初始化所有service TODO
	*/
	ds.chainMap[spectrumevents.ChainName], err = spectrum.NewSMCService("", 0)
	ds.chainMap[ethereumevents.ChainName], err = ethereum.NewETHService("", 0)
	return
}

// Start :
func (ds *DispatchService) Start() {
	/*
		启动restful请求调度线程,包含用户及公证人节点的restful请求
	*/
	go ds.restfulRequestDispatcherLoop()
	/*
		对每条链启动事件调度线程
	*/
	for _, v := range ds.chainMap {
		go ds.chainEventDispatcherLoop(v)
	}
	log.Info("DispatchService start complete ...")
}

// StopAll :
func (ds *DispatchService) StopAll() {
	close(ds.quitChan)
}

/*
常驻线程
restful请求调度线程,负责监听所有的api调用并分发至对应的service进行处理
*/
func (ds *DispatchService) restfulRequestDispatcherLoop() {
	logPrefix := "RestfulRequest : "
	log.Info(fmt.Sprintf("%s dispather start ...", logPrefix))
	notaryRequestChan := ds.notaryAPI.GetRequestChan()
	userRequestChan := ds.userAPI.GetRequestChan()
	for {
		select {
		case req, ok := <-notaryRequestChan:
			if !ok {
				err := fmt.Errorf("%s read notary request chan err ", logPrefix)
				panic(err)
			}
			ds.dispatchRestfulRequest(req)
		case req, ok := <-userRequestChan:
			if !ok {
				err := fmt.Errorf("%s read user request chan err ", logPrefix)
				panic(err)
			}
			ds.dispatchRestfulRequest(req)
		case <-ds.quitChan:
			log.Info(fmt.Sprintf("%s dispather stop success", logPrefix))
			return
		}
	}
}

func (ds *DispatchService) dispatchRestfulRequest(req api.Request) {
	logPrefix := "NotaryRequest : "
	/*
		restful 请求调度规则如下,优先级从高到低:
		1. CrossChainRequest 带SCToken的请求,下发至对应的CrossChainService,如果找不到,返回错误
		2. NotaryRequest 带key且不带SCToken的请求,一定是公证人之间的请求,比如私钥生成过程中的消息交互,下发至NotaryService
		3. Request 不带key且不带SCToken的请求,一定为管理用户的非交易请求,下发至SystemService
		TODO
	*/
	// 跨链交易相关请求,下发至对应service
	switch r := req.(type) {
	case api.CrossChainRequest:
		// 跨链交易相关请求,下发至对应service
		service, ok := ds.scToken2CrossChainServiceMap[r.GetSCTokenAddress()]
		if !ok {
			log.Error(fmt.Sprintf("%s receive request with out notary service : \n%s\n", logPrefix, utils.ToJsonStringFormat(req)))
			// 返回api错误 TODO
			return
		}
		go service.OnRequest(req)
		return
	case api.NotaryRequest:
		// TODO
	case api.Request:
		// TODO
	}
}

/*
常驻线程
事件调度线程,负责监听所有链的所有事件并分发至对应的notaryService
*/
func (ds *DispatchService) chainEventDispatcherLoop(c chain.Chain) {
	logPrefix := fmt.Sprintf("Chain %s : ", c.GetChainName())
	log.Info(fmt.Sprintf("%s dispather start ...", logPrefix))
	for {
		select {
		case e, ok := <-c.GetEventChan():
			if !ok {
				err := fmt.Errorf("%s read event chan err ", logPrefix)
				panic(err)
			}
			if e.GetEventName() == chain.NewBlockNumberEventName {
				// 新块事件,dispatch至所有service
				ds.dispatchEvent2All(e)
			} else {
				// 合约事件,根据合约地址调度
				ds.dispatchEvent(e)
			}
		case <-ds.quitChan:
			log.Info(fmt.Sprintf("%s dispather stop success", logPrefix))
			return
		}
	}
}

func (ds *DispatchService) dispatchEvent2All(e chain.Event) {
	for _, service := range ds.scToken2CrossChainServiceMap {
		// 这里在处理区块高度的时候,会不会导致协程数量过大 TODO
		go service.OnEvent(e)
	}
}

func (ds *DispatchService) dispatchEvent(e chain.Event) {
	if e.GetSCTokenAddress() == utils.EmptyAddress {
		// 主链事件,根据主链合约地址FromAddress调度,遍历,后续可优化,维护一个主链合约地址-SCToken地址的map即可
		for _, service := range ds.scToken2CrossChainServiceMap {
			if service.GetMCContractAddress() == e.GetFromAddress() {
				// 事件业务逻辑处理
				go service.OnEvent(e)
				// 每个事件应该只有一个对应service,所以这里处理完毕直接return
				return
			}
		}
	} else {
		// 侧链事件,直接根据SCToken地址调度
		service, ok := ds.scToken2CrossChainServiceMap[e.GetSCTokenAddress()]
		if !ok {
			// never happen
			panic(fmt.Errorf("can not find CrossChainService with SCToken %s", e.GetSCTokenAddress().String()))
		}
		// 事件业务逻辑处理
		go service.OnEvent(e)
	}
}
