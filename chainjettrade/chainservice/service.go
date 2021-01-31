package chainservice

import (
	"context"
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"

	"time"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/commons"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nkbai/log"
)

type ChainOperate interface {
	parserLogsToEventsAndSort(logs []types.Log) (es []chain.Event, err error)
	GetChainName() string
	//CreateNewBlockEvent(blockNumber uint64) chainjettrade.NewBlockEvent
}

// ChainService :
type ChainService struct {
	Name            string
	c               *chainjettrade.SafeEthClient
	host            string
	lastBlockNumber uint64
	ContractAddress common.Address

	connectStatus                  commons.ConnectStatus
	connectStatusChangeChanMap     map[string]chan commons.ConnectStatusChange
	connectStatusChangeChanMapLock sync.Mutex
	eventsDone                     map[string]uint64 // 事件处理历史记录 key: txhash+topic 一个tx中会发生好几个event
	eventChan                      chan chain.Event
	listenerQuitChan               chan struct{}
	co                             ChainOperate
}

func NewChainService(host, name string, contractAddress common.Address) (ss *ChainService, err error) {
	// init client
	var c *ethclient.Client
	ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.SMC.RPCTimeout)
	c, err = ethclient.DialContext(ctx, host)
	cancelFunc()
	if err != nil {
		return
	}
	ss = &ChainService{
		Name:                       name,
		c:                          chainjettrade.NewSafeClient(c),
		host:                       host,
		ContractAddress:            contractAddress,
		connectStatus:              commons.Disconnected,
		connectStatusChangeChanMap: make(map[string]chan commons.ConnectStatusChange),
		eventsDone:                 make(map[string]uint64),
		eventChan:                  make(chan chain.Event, 100),
	}
	err = ss.checkConnectStatus()
	if err != nil {
		return
	}
	ss.changeStatus(commons.Connected)
	return
}

// SetLastBlockNumber :
func (ss *ChainService) SetLastBlockNumber(lastBlockNumber uint64) {
	ss.lastBlockNumber = lastBlockNumber
}

// StartEventListener :
func (ss *ChainService) StartEventListener() error {
	ss.listenerQuitChan = make(chan struct{})
	go ss.loop()
	return nil
}

// StopEventListener :
func (ss *ChainService) StopEventListener() {
	if ss.listenerQuitChan != nil {
		close(ss.listenerQuitChan)
		ss.listenerQuitChan = nil
	}
}

// GetEventChan :
func (ss *ChainService) GetEventChan() <-chan chain.Event {
	return ss.eventChan
}

// RegisterConnectStatusChangeChan :
func (ss *ChainService) RegisterConnectStatusChangeChan(name string) <-chan commons.ConnectStatusChange {
	ch, ok := ss.connectStatusChangeChanMap[name]
	if ok {
		log.Warn("%s RegisterConnectStatusChangeChan should only call once", ss.Name)
		return ch
	}
	ch = make(chan commons.ConnectStatusChange, 1)
	ss.connectStatusChangeChanMapLock.Lock()
	ss.connectStatusChangeChanMap[name] = ch
	ss.connectStatusChangeChanMapLock.Unlock()
	return ch
}

// UnRegisterConnectStatusChangeChan :
func (ss *ChainService) UnRegisterConnectStatusChangeChan(name string) {
	ch, ok := ss.connectStatusChangeChanMap[name]
	ss.connectStatusChangeChanMapLock.Lock()
	delete(ss.connectStatusChangeChanMap, name)
	ss.connectStatusChangeChanMapLock.Unlock()
	if ok && ch != nil {
		close(ch)
	}
}

// RecoverDisconnect :
func (ss *ChainService) RecoverDisconnect() {
	var err error
	var c *ethclient.Client
	ss.changeStatus(commons.Reconnecting)
	if ss.c != nil && ss.c.Client != nil {
		ss.c.Client.Close()
	}
	for {
		log.Info("%s tyring to reconnect smc ...", ss.Name)
		select {
		case <-ss.listenerQuitChan:
			ss.changeStatus(commons.Closed)
			return
		default:
			//never block
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.SMC.RPCTimeout)
		c, err = ethclient.DialContext(ctx, ss.host)
		cancelFunc()
		ss.c = chainjettrade.NewSafeClient(c)
		if err == nil {
			err = ss.checkConnectStatus()
		}
		if err == nil {
			//reconnect ok
			ss.changeStatus(commons.Connected)
			return
		}
		log.Info(fmt.Sprintf("%s reconnect to %s error: %s", ss.Name, ss.host, err))
		time.Sleep(time.Second * 3)
	}
}

func (ss *ChainService) changeStatus(newStatus commons.ConnectStatus) {
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
	log.Info(fmt.Sprintf("%s connect status change from %d to %d", ss.Name, sc.OldStatus, sc.NewStatus))
}

// GetConn : impl chaintmp.Chain
func (ss *ChainService) GetConn() *ethclient.Client {
	return ss.c.Client
}

func (ss *ChainService) checkConnectStatus() (err error) {
	if ss.c == nil || ss.c.Client == nil {
		return chainjettrade.ErrNotConnected
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.SMC.RPCTimeout)
	defer cancelFunc()
	_, err = ss.c.HeaderByNumber(ctx, big.NewInt(1))
	if err != nil {
		return
	}
	return
}

// 事件监听主线程,理论上常驻,自动重连
func (ss *ChainService) loop() {
	log.Trace(fmt.Sprintf("%s.EventListener start getting lasted block number from blocknubmer=%d", ss.Name, ss.lastBlockNumber))
	currentBlock := ss.lastBlockNumber
	retryTime := 0
	for {
		ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.SMC.RPCTimeout)
		h, err := ss.c.HeaderByNumber(ctx, nil)
		if err != nil {
			log.Error(fmt.Sprintf("%s.EventListener HeaderByNumber err=%s", ss.Name, err))
			cancelFunc()
			if ss.listenerQuitChan != nil {
				go ss.RecoverDisconnect()
				// 阻塞等待重连成功,继续循环
				ch := ss.RegisterConnectStatusChangeChan("self")
				for {
					sc := <-ch
					if sc.NewStatus == commons.Closed {
						log.Info("%s.EventListener end because user closed SmcService", ss.Name)
						return
					}
					if sc.NewStatus == commons.Connected {
						ss.UnRegisterConnectStatusChangeChan("self")
						log.Trace(fmt.Sprintf("%s.EventListener reconnected success, start getting lasted block number from blocknubmer=%d", ss.Name, ss.lastBlockNumber))
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
			time.Sleep(cfg.SMC.BlockNumberPollPeriod / 2)
			retryTime++
			if retryTime > 10 {
				//log.Warn(fmt.Sprintf("SmcService.EventListener get same block number %d from chaintmp %d times,maybe something wrong with smc ...", lastedBlock, retryTime))
			}
			continue
		}
		retryTime = 0
		if lastedBlock != currentBlock+1 {
			log.Warn(fmt.Sprintf("%s.EventListener missed %d blocks", ss.Name, lastedBlock-currentBlock-1))
		}
		if lastedBlock%cfg.SMC.BlockNumberLogPeriod == 0 {
			log.Trace(fmt.Sprintf("Spectrum new block : %d", lastedBlock))
		}
		var fromBlockNumber, toBlockNumber uint64
		if currentBlock < 2*cfg.SMC.ConfirmBlockNumber {
			fromBlockNumber = 0
		} else {
			fromBlockNumber = currentBlock - 2*cfg.SMC.ConfirmBlockNumber
		}
		if lastedBlock < cfg.SMC.ConfirmBlockNumber {
			toBlockNumber = 0
		} else {
			toBlockNumber = lastedBlock - cfg.SMC.ConfirmBlockNumber
		}
		// get all events between currentBlock and confirmBlock
		es, err := ss.queryAllEvents(fromBlockNumber, toBlockNumber)
		if err != nil {
			log.Error(fmt.Sprintf("%s.EventListener queryAllStateChange err=%s", ss.Name, err))
			// 如果这里出现err,不能继续处理该blocknumber,否则会丢事件,直接从该块重新处理即可
			time.Sleep(cfg.SMC.BlockNumberPollPeriod / 2)
			continue
		}
		if len(es) > 0 {
			log.Trace(fmt.Sprintf("receive %d events of   between block %d - %d", len(es), currentBlock+1, lastedBlock))
		}

		// refresh block number and notify PhotonService
		currentBlock = lastedBlock
		ss.lastBlockNumber = currentBlock
		//block event是统一的,有原始跨链合约维护
		//ss.eventChan <- ss.co.CreateNewBlockEvent(currentBlock)

		// notify Photon service
		for _, e := range es {
			log.Info("%s send event,event=%s", ss.Name, log.StringInterface(e, 3))
			ss.eventChan <- e
		}
		// 清除过期流水
		for key, blockNumber := range ss.eventsDone {
			if blockNumber <= fromBlockNumber-5 { //清理5块之前的
				delete(ss.eventsDone, key)
			}
		}
		// wait to next time
		select {
		case <-time.After(cfg.SMC.BlockNumberPollPeriod):
		case <-ss.listenerQuitChan:
			ss.listenerQuitChan = nil
			log.Info(fmt.Sprintf("%s.EventListener quit complete", ss.Name))
			return
		}
	}
}

func (ss *ChainService) queryAllEvents(fromBlockNumber uint64, toBlockNumber uint64) (es []chain.Event, err error) {
	/*
		get all event of contract TokenNetworkRegistry, SecretRegistry , TokenNetwork
	*/
	logs, err := ss.getLogsFromChain(fromBlockNumber, toBlockNumber)
	if err != nil {
		return
	}
	return ss.co.parserLogsToEventsAndSort(logs)
}

func (ss *ChainService) getLogsFromChain(fromBlock uint64, toBlock uint64) (logs []types.Log, err error) {

	var contractsAddress = []common.Address{ss.ContractAddress}
	var q *ethereum.FilterQuery
	q, err = buildQueryBatch(contractsAddress, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	return ss.c.FilterLogs(getQueryContext(), *q)
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
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(cfg.SMC.RPCTimeout))
	if cf != nil {
	}
	return ctx
}
