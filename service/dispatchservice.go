package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/accounts"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum"
	ethereumevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum"
	spectrumevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nkbai/log"
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
	adminService                     *AdminService
	notaryService                    *NotaryService
	scToken2CrossChainServiceMap     map[common.Address]*CrossChainService
	scToken2CrossChainServiceMapLock sync.Mutex
}

// NewDispatchService :
func NewDispatchService(cfg *params.Config) (ds *DispatchService, err error) {
	// 1. 加载私钥
	am := accounts.NewAccountManager(cfg.KeystorePath)
	privateKeyBin, err := am.GetPrivateKey(cfg.Address, cfg.Password)
	if err != nil {
		log.Error("load private key err : %s", err.Error())
		return
	}
	privateKey, err := crypto.ToECDSA(privateKeyBin)
	if err != nil {
		log.Error("load private key err : %s", err.Error())
		return
	}
	// 2. 打开db,如果是第一次启动,读取notary.conf并写入db
	db := models.SetUpDB("sqlite3", cfg.DataBasePath)
	notaries, err := db.GetNotaryInfo()
	if err != nil {
		log.Error("get notary info from db err : %s", err.Error())
		return
	}
	if len(notaries) == 0 {
		// first start
		notaries, err = db.NewNotaryInfoFromConfFile(cfg.NotaryConfFilePath)
		if err != nil {
			log.Error("get notary info from file %s err : %s", cfg.NotaryConfFilePath, err.Error())
			return
		}
		if len(notaries) < 2 {
			err = errors.New("notary num should not < 2")
			log.Error("notary num should not < 2")
			return
		}
	}
	// 2.5 根据notaries数量初始化 ShareCount及ThresholdCount
	params.ShareCount = len(notaries)
	params.ThresholdCount = params.ShareCount / 3 * 2
	if params.ShareCount%3 > 1 {
		params.ThresholdCount++
	}
	log.Info("ShareCount=%d ThresholdCount=%d", params.ShareCount, params.ThresholdCount)
	// 3. init dispatch service
	ds = &DispatchService{
		userAPI:                      userapi.NewUserAPI(cfg.UserAPIListen),
		notaryAPI:                    notaryapi.NewNotaryAPI(cfg.NotaryAPIListen),
		chainMap:                     make(map[string]chain.Chain),
		db:                           db,
		quitChan:                     make(chan struct{}),
		scToken2CrossChainServiceMap: make(map[common.Address]*CrossChainService),
	}
	// 4. 初始化侧链事件监听
	chainName := spectrumevents.ChainName
	ds.chainMap[chainName], err = spectrum.NewSMCService(cfg.SmcRPCEndPoint, db.GetLastBlockNumber(chainName))
	if err != nil {
		log.Error("new SMCService err : %s", err.Error())
		return
	}
	// 5. 初始化主链事件监听
	chainName = ethereumevents.ChainName
	ds.chainMap[chainName], err = ethereum.NewETHService(cfg.EthRPCEndPoint, db.GetLastBlockNumber(chainName))
	if err != nil {
		log.Error("new ETHService err : %s", err.Error())
		return
	}
	// 5. 初始化NotaryService
	ds.notaryService, err = NewNotaryService(db, privateKey, notaries)
	if err != nil {
		log.Error("init NotaryService err : %s", err.Error())
		return
	}
	// 6. 初始化AdminService
	ds.adminService, err = NewAdminService(db, ds.notaryService, ds.chainMap)
	if err != nil {
		log.Error("init AdminService err : %s", err.Error())
		return
	}
	// 7. 根据SCTokenMetaInfo初始化CrossChainService,并将所有合约地址注册到对应链的监听器中
	scTokenMetaInfoList, err := ds.db.GetSCTokenMetaInfoList()
	if err != nil {
		log.Error("GetSCTokenMetaInfoList err : %s", err.Error())
		return
	}
	ds.scToken2CrossChainServiceMapLock.Lock()
	for _, scTokenMetaInfo := range scTokenMetaInfoList {
		ds.scToken2CrossChainServiceMap[scTokenMetaInfo.SCToken] = NewCrossChainService(scTokenMetaInfo)
		// 注册侧链合约:
		err = ds.chainMap[spectrumevents.ChainName].RegisterEventListenContract(scTokenMetaInfo.SCToken)
		if err != nil {
			log.Error("RegisterEventListenContract on side chain err : %s", err.Error())
			return
		}
		// 注册主链合约:
		err = ds.chainMap[scTokenMetaInfo.MCName].RegisterEventListenContract(scTokenMetaInfo.MCLockedContractAddress)
		if err != nil {
			log.Error("RegisterEventListenContract on main chain %s err : %s", scTokenMetaInfo.MCName, err.Error())
			return
		}
	}
	ds.scToken2CrossChainServiceMapLock.Unlock()
	return
}

// Start :
func (ds *DispatchService) Start() (err error) {
	// 1. 启动restful请求调度线程,包含用户及公证人节点的restful请求
	go ds.restfulRequestDispatcherLoop()
	// 2. 对每条链启动事件调度线程
	for _, v := range ds.chainMap {
		go ds.chainEventDispatcherLoop(v)
	}
	// 3. 启动所有链的事件监听
	for _, c := range ds.chainMap {
		err = c.StartEventListener()
		if err != nil {
			return
		}
	}
	// 4. 启动API
	ds.notaryAPI.Start(false)
	ds.userAPI.Start(true)
	log.Info("DispatchService start complete ...")
	return
}

// Stop :
func (ds *DispatchService) Stop() {
	/*
		关闭所有链的事件监听
	*/
	for _, c := range ds.chainMap {
		c.StopEventListener()
	}
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
	switch r := req.(type) {
	case api.CrossChainRequest:
		service, ok := ds.scToken2CrossChainServiceMap[r.GetSCTokenAddress()]
		if !ok {
			log.Error(fmt.Sprintf("%s receive request with out notary service : \n%s\n", logPrefix, utils.ToJSONStringFormat(req)))
			// 返回api错误
			req.WriteErrorResponse(api.ErrorCodeException)
			return
		}
		go service.OnRequest(req)
		return
	case api.NotaryRequest:
		go ds.notaryService.OnRequest(req)
	case api.Request:
		go ds.adminService.OnRequest(req)
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
				// 新块事件
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
	// 通知所有Service
	ds.adminService.OnEvent(e)
	ds.notaryService.OnEvent(e)
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
