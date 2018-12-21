package spectrum

import (
	"context"
	"fmt"
	"sync"

	"time"

	"math/big"

	"errors"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/distributed-notary/commons"
	"github.com/SmartMeshFoundation/distributed-notary/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/ethereum/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ETHService :
type ETHService struct {
	c               *client.SafeEthClient
	host            string
	lastBlockNumber uint64

	lockedEthereumProxyMap map[common.Address]*proxy.LockedEthereumProxy
	tokenProxyMapLock      sync.Mutex

	connectStatus                  commons.ConnectStatus
	connectStatusChangeChanMap     map[string]chan commons.ConnectStatusChange
	connectStatusChangeChanMapLock sync.Mutex

	eventsDone map[common.Hash]uint64 // 事件处理历史记录

	eventChan        chan events.Event
	listenerQuitChan chan struct{}
}

// NewETHService :
func NewETHService(host string, lastBlockNumber uint64, contractAddresses ...common.Address) (ss *ETHService, err error) {
	// init client
	var c *ethclient.Client
	ctx, cancelFunc := context.WithTimeout(context.Background(), spectrumRPCTimeout)
	c, err = ethclient.DialContext(ctx, host)
	cancelFunc()
	if err != nil {
		return
	}
	ss = &ETHService{
		c:                          client.NewSafeClient(c),
		host:                       host,
		connectStatus:              commons.Disconnected,
		lastBlockNumber:            lastBlockNumber,
		lockedEthereumProxyMap:     make(map[common.Address]*proxy.LockedEthereumProxy),
		connectStatusChangeChanMap: make(map[string]chan commons.ConnectStatusChange),
		eventChan:                  make(chan events.Event, 100),
		eventsDone:                 make(map[common.Hash]uint64),
	}
	err = ss.checkConnectStatus()
	if err != nil {
		return
	}
	ss.changeStatus(commons.Connected)
	// init proxy
	if len(contractAddresses) > 0 {
		for _, addr := range contractAddresses {
			// init proxy
			p, err2 := proxy.NewLockedEthereumProxy(ss.c, addr)
			if err = err2; err != nil {
				return
			}
			ss.lockedEthereumProxyMap[addr] = p
		}
	}
	return
}

// RegisterEventListenContract :
func (ss *ETHService) RegisterEventListenContract(contractAddresses ...common.Address) error {
	if ss.connectStatus != commons.Connected {
		return errors.New("ETHService can not register when not connected")
	}
	ss.tokenProxyMapLock.Lock()
	for _, addr := range contractAddresses {
		if proxy, ok := ss.lockedEthereumProxyMap[addr]; ok && proxy != nil {
			continue
		}
		// init proxy
		p, err := proxy.NewLockedEthereumProxy(ss.c, addr)
		if err != nil {
			return err
		}
		ss.lockedEthereumProxyMap[addr] = p
	}
	ss.tokenProxyMapLock.Unlock()
	return nil
}

// UnRegisterEventListenContract :
func (ss *ETHService) UnRegisterEventListenContract(contractAddresses ...common.Address) {
	ss.tokenProxyMapLock.Lock()
	for _, addr := range contractAddresses {
		delete(ss.lockedEthereumProxyMap, addr)
	}
	ss.tokenProxyMapLock.Unlock()
}

// StartEventListener :
func (ss *ETHService) StartEventListener() error {
	ss.listenerQuitChan = make(chan struct{})
	go ss.loop()
	return nil
}

// StopEventListener :
func (ss *ETHService) StopEventListener() error {
	if ss.listenerQuitChan != nil {
		close(ss.listenerQuitChan)
		ss.listenerQuitChan = nil
	}
	return nil
}

// GetEventChan :
func (ss *ETHService) GetEventChan() <-chan events.Event {
	return ss.eventChan
}

// RegisterConnectStatusChangeChan :
func (ss *ETHService) RegisterConnectStatusChangeChan(name string) <-chan commons.ConnectStatusChange {
	ch, ok := ss.connectStatusChangeChanMap[name]
	if ok {
		log.Warn("ETHService RegisterConnectStatusChangeChan should only call once")
		return ch
	}
	ch = make(chan commons.ConnectStatusChange, 1)
	ss.connectStatusChangeChanMapLock.Lock()
	ss.connectStatusChangeChanMap[name] = ch
	ss.connectStatusChangeChanMapLock.Unlock()
	return ch
}

// UnRegisterConnectStatusChangeChan :
func (ss *ETHService) UnRegisterConnectStatusChangeChan(name string) {
	ch, ok := ss.connectStatusChangeChanMap[name]
	ss.connectStatusChangeChanMapLock.Lock()
	delete(ss.connectStatusChangeChanMap, name)
	ss.connectStatusChangeChanMapLock.Unlock()
	if ok && ch != nil {
		close(ch)
	}
}

// RecoverDisconnect :
func (ss *ETHService) RecoverDisconnect() {
	var err error
	var c *ethclient.Client
	ss.changeStatus(commons.Reconnecting)
	if ss.c != nil && ss.c.Client != nil {
		ss.c.Client.Close()
	}
	for {
		log.Info("ETHService tyring to reconnect geth ...")
		select {
		case <-ss.listenerQuitChan:
			ss.changeStatus(commons.Closed)
			return
		default:
			//never block
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), spectrumRPCTimeout)
		c, err = ethclient.DialContext(ctx, ss.host)
		cancelFunc()
		ss.c = client.NewSafeClient(c)
		if err == nil {
			err = ss.checkConnectStatus()
		}
		if err == nil {
			//reconnect ok
			ss.changeStatus(commons.Connected)
			return
		}
		log.Info(fmt.Sprintf("ETHService reconnect to %s error: %s", ss.host, err))
		time.Sleep(time.Second * 3)
	}
}

func (ss *ETHService) changeStatus(newStatus commons.ConnectStatus) {
	sc := &commons.ConnectStatusChange{
		OldStatus:  ss.connectStatus,
		NewStatus:  newStatus,
		ChangeTime: time.Now(),
	}
	ss.connectStatus = newStatus
	ss.connectStatusChangeChanMapLock.Lock()
	for _, ch := range ss.connectStatusChangeChanMap {
		select {
		case ch <- *sc:
		default:
			// never block
		}
	}
	ss.connectStatusChangeChanMapLock.Unlock()
	log.Info(fmt.Sprintf("ETHService connect status change from %d to %d", sc.OldStatus, sc.NewStatus))
}

func (ss *ETHService) checkConnectStatus() (err error) {
	if ss.c == nil || ss.c.Client == nil {
		return client.ErrNotConnected
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), spectrumRPCTimeout)
	defer cancelFunc()
	_, err = ss.c.HeaderByNumber(ctx, big.NewInt(1))
	if err != nil {
		return
	}
	return
}

// 事件监听主线程,理论上常驻,自动重连
func (ss *ETHService) loop() {
	log.Trace(fmt.Sprintf("ETHService.EventListener start getting lasted block number from blocknubmer=%d", ss.lastBlockNumber))
	currentBlock := ss.lastBlockNumber
	retryTime := 0
	for {
		ctx, cancelFunc := context.WithTimeout(context.Background(), spectrumRPCTimeout)
		h, err := ss.c.HeaderByNumber(ctx, nil)
		if err != nil {
			log.Error(fmt.Sprintf("ETHService.EventListener HeaderByNumber err=%s", err))
			cancelFunc()
			if ss.listenerQuitChan != nil {
				go ss.RecoverDisconnect()
				// 阻塞等待重连成功,继续循环
				ch := ss.RegisterConnectStatusChangeChan("self")
				for {
					sc := <-ch
					if sc.NewStatus == commons.Closed {
						log.Info("ETHService.EventListener end because user closed SmcService")
						return
					}
					if sc.NewStatus == commons.Connected {
						ss.UnRegisterConnectStatusChangeChan("self")
						log.Trace(fmt.Sprintf("ETHService.EventListener reconnected success, start getting lasted block number from blocknubmer=%d", ss.lastBlockNumber))
						// 重连成功,刷新proxy
						ss.refreshContractProxy()
						break
					}
				}
				continue
			}
		}
		cancelFunc()
		lastedBlock := h.Number.Uint64()
		// 这里如果出现切换公链导致获取到的新块比当前块更小的话,只需要等待即可
		if currentBlock >= lastedBlock {
			time.Sleep(pollPeriod / 2)
			retryTime++
			if retryTime > 10 {
				log.Warn(fmt.Sprintf("ETHService.EventListener get same block number %d from chain %d times,maybe something wrong with geth ...", lastedBlock, retryTime))
			}
			continue
		}
		retryTime = 0
		if lastedBlock != currentBlock+1 {
			log.Warn(fmt.Sprintf("ETHService.EventListener missed %d blocks", lastedBlock-currentBlock-1))
		}
		if lastedBlock%logPeriod == 0 {
			log.Trace(fmt.Sprintf("ETHService.EventListener new block : %d", lastedBlock))
		}
		var fromBlockNumber, toBlockNumber uint64
		if currentBlock < 2*confirmBlockNumber {
			fromBlockNumber = 0
		} else {
			fromBlockNumber = currentBlock - 2*confirmBlockNumber
		}
		if currentBlock < confirmBlockNumber {
			toBlockNumber = 0
		} else {
			toBlockNumber = currentBlock - confirmBlockNumber
		}
		// get all events between currentBlock and confirmBlock
		es, err := ss.queryAllEvents(fromBlockNumber, toBlockNumber)
		if err != nil {
			log.Error(fmt.Sprintf("ETHService.EventListener queryAllStateChange err=%s", err))
			// 如果这里出现err,不能继续处理该blocknumber,否则会丢事件,直接从该块重新处理即可
			time.Sleep(pollPeriod / 2)
			continue
		}
		if len(es) > 0 {
			log.Trace(fmt.Sprintf("receive %d events of %d contracts between block %d - %d", len(es), len(ss.lockedEthereumProxyMap), currentBlock+1, lastedBlock))
		}

		// refresh block number and notify PhotonService
		currentBlock = lastedBlock
		ss.lastBlockNumber = currentBlock
		ss.eventChan <- events.CreateNewBlockEvent(currentBlock)

		// notify Photon service
		for _, e := range es {
			ss.eventChan <- e
		}

		// 清除过期流水
		for key, blockNumber := range ss.eventsDone {
			if blockNumber <= fromBlockNumber {
				delete(ss.eventsDone, key)
			}
		}
		// wait to next time
		select {
		case <-time.After(pollPeriod):
		case <-ss.listenerQuitChan:
			ss.listenerQuitChan = nil
			log.Info(fmt.Sprintf("ETHService.EventListener quit complete"))
			return
		}
	}
}

func (ss *ETHService) refreshContractProxy() {
	ss.tokenProxyMapLock.Lock()
	for addr := range ss.lockedEthereumProxyMap {
		// rebuild proxy
		p, err := proxy.NewLockedEthereumProxy(ss.c, addr)
		if err != nil {
			log.Error(fmt.Sprintf("ETHService refreshContractProxy err : %s", err.Error()))
			continue
		}
		ss.lockedEthereumProxyMap[addr] = p
	}
	ss.tokenProxyMapLock.Unlock()
}

func (ss *ETHService) queryAllEvents(fromBlockNumber uint64, toBlockNumber uint64) (es []events.Event, err error) {
	/*
		get all event of contract TokenNetworkRegistry, SecretRegistry , TokenNetwork
	*/
	logs, err := ss.getLogsFromChain(fromBlockNumber, toBlockNumber)
	if err != nil {
		return
	}
	return ss.parserLogsToEventsAndSort(logs)
}

func (ss *ETHService) getLogsFromChain(fromBlock uint64, toBlock uint64) (logs []types.Log, err error) {
	if len(ss.lockedEthereumProxyMap) == 0 {
		return
	}
	ss.tokenProxyMapLock.Lock()
	defer ss.tokenProxyMapLock.Unlock()
	var contractsAddress []common.Address
	for key := range ss.lockedEthereumProxyMap {
		contractsAddress = append(contractsAddress, key)
	}
	var q *ethereum.FilterQuery
	q, err = buildQueryBatch(contractsAddress, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	return ss.c.FilterLogs(getQueryContext(), *q)
}

func (ss *ETHService) parserLogsToEventsAndSort(logs []types.Log) (es []events.Event, err error) {
	if len(logs) == 0 {
		return
	}
	for _, l := range logs {
		eventName := events.TopicToEventName[l.Topics[0]]
		// 根据已处理流水去重
		if doneBlockNumber, ok := ss.eventsDone[l.TxHash]; ok {
			if doneBlockNumber == l.BlockNumber {
				//log.Trace(fmt.Sprintf("get event txhash=%s repeated,ignore...", l.TxHash.String()))
				continue
			}
			log.Warn(fmt.Sprintf("ETHService.EventListener event tx=%s happened at %d, but now happend at %d ", l.TxHash.String(), doneBlockNumber, l.BlockNumber))
		}
		switch eventName {
		case events.LockedEthereumPrepareLockinEventName:
			e, err2 := events.CreatePrepareLockinEvent(l)
			if err = err2; err != nil {
				return
			}
			es = append(es, e)
		case events.LockedEthereumLockoutSecretEventName:
			e, err2 := events.CreateLockoutSecretEvent(l)
			if err = err2; err != nil {
				return
			}
			es = append(es, e)
		case events.LockedEthereumPrePareLockedOutEventName:
			e, err2 := events.CreatePrepareLockoutEvent(l)
			if err = err2; err != nil {
				return
			}
			es = append(es, e)
		default:
			log.Warn(fmt.Sprintf("ETHService.EventListener receive unkonwn type event from chain : \n%s\n", utils.ToJsonStringFormat(l)))
		}
		// 记录处理流水
		ss.eventsDone[l.TxHash] = l.BlockNumber
	}
	return
}

func buildQueryBatch(contractsAddress []common.Address, fromBlock uint64, toBlock uint64) (q *ethereum.FilterQuery, err error) {
	q = &ethereum.FilterQuery{}
	if fromBlock == 0 {
		q.FromBlock = nil
	} else {
		q.FromBlock = big.NewInt(int64(fromBlock))
	}
	if toBlock == 0 {
		q.ToBlock = nil
	} else {
		q.ToBlock = big.NewInt(int64(toBlock))
	}
	q.Addresses = contractsAddress
	return
}

func getQueryContext() context.Context {
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(defaultPollTimeout))
	if cf != nil {
	}
	return ctx
}
