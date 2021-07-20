package main

import (
	"context"

	"sync"

	"time"

	"net/http"

	"encoding/json"

	"bytes"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

var defaultTimeoutBlock uint64 = 15

type nonceManager struct {
	account                 common.Address
	chainName               string
	nextNonce               uint64
	c                       chainClient
	usedNonceToCancelURLMap *sync.Map
	lock                    sync.Mutex
}

func newNonceManager(account common.Address, chainName string, c chainClient) *nonceManager {
	nm := &nonceManager{
		account:                 account,
		chainName:               chainName,
		c:                       c,
		usedNonceToCancelURLMap: new(sync.Map),
	}
	var err error
	nm.nextNonce, err = c.NonceAt(context.Background(), account, nil)
	if err != nil {
		log.Error("NonceAt error %s", err.Error())
		nm.nextNonce = 0
	}
	return nm
}

func (nm *nonceManager) applyNonce(cancelURL string) uint64 {
	if nm.nextNonce == 0 {
		var err error
		nm.nextNonce, err = nm.c.PendingNonceAt(context.Background(), nm.account)
		if err != nil {
			log.Error("PendingNonceAt error %s", err.Error())
			nm.nextNonce = 0
		}
	}
	nm.lock.Lock()
	nonceToUse := nm.nextNonce
	nm.nextNonce++
	nm.usedNonceToCancelURLMap.Store(nonceToUse, cancelURL)
	nm.lock.Unlock()
	//go nm.confirmLoop(nonceToUse)
	return nonceToUse
}

/*
nonce确认策略:
1. 轮询间隔1秒
2. 如果15块之后,该nonce没有被使用,则调用map中保存的cancelUrl发起一笔无效交易消耗该nonce
*/
func (nm *nonceManager) confirmLoop(nonceUsed uint64) {
	startBlock, err := nm.c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Error("HeaderByNumber err %s", err.Error())
		return
	}
	startBlockNumber := startBlock.Number.Uint64()
	for {
		time.Sleep(time.Second)
		nonce, err2 := nm.c.NonceAt(context.Background(), nm.account, nil)
		if err2 != nil {
			log.Error("NonceAt err %s", err2.Error())
			continue
		}
		if nonce > nonceUsed {
			// 确认
			nm.confirmNonce(nonceUsed)
			return
		}
		block, err2 := nm.c.HeaderByNumber(context.Background(), nil)
		if err2 != nil {
			log.Error("HeaderByNumber err %s", err2.Error())
			continue
		}
		blockNumber := block.Number.Uint64()
		if blockNumber-startBlockNumber >= defaultTimeoutBlock {
			nm.lock.Lock()
			if nm.nextNonce-nonceUsed == 1 {
				// 说明这个nonce之后我没有分配新的nonce,且该nonce没有被使用
				nm.nextNonce = nonceUsed
				nm.reuseNonce(nonceUsed)
				nm.lock.Unlock()
				return
			}
			nm.lock.Unlock()
			//说明在这个nonce之后我又分配了后续的nonce,所以很可能有大于该nonce的交易在排队,为了不妨碍后续交易,需要消耗掉该nonce
			nm.cancelNonce(nonceUsed)
			return
		}
	}
}

func (nm *nonceManager) reuseNonce(nonceUsed uint64) {
	nm.usedNonceToCancelURLMap.Delete(nonceUsed)
	log.Info("account=%s chain=%s nonce=%d reuse", nm.account.String(), nm.chainName, nonceUsed)
}

func (nm *nonceManager) confirmNonce(nonceUsed uint64) {
	nm.usedNonceToCancelURLMap.Delete(nonceUsed)
	log.Info("account=%s chain=%s nonce=%d confirm", nm.account.String(), nm.chainName, nonceUsed)
}

func (nm *nonceManager) cancelNonce(nonceUsed uint64) {
	cancelURLInterface, ok := nm.usedNonceToCancelURLMap.Load(nonceUsed)
	if !ok {
		panic("never happen")
	}
	cancelURL := cancelURLInterface.(string)
	req := userapi.NewCancelNonceRequest(nm.chainName, nm.account, nonceUsed)
	payload, err := json.Marshal(req)
	if err != nil {
		log.Error("account=%s chain=%s nonce=%d cancelNonce error %s", nm.account.String(), nm.chainName, nonceUsed, err.Error())
		return
	}
	/* #nosec */
	resp, err := http.Post(cancelURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Error("account=%s chain=%s nonce=%d cancelNonce error %s", nm.account.String(), nm.chainName, nonceUsed, err.Error())
		return
	}
	var buf [4096 * 1024]byte
	n := 0
	n, err = resp.Body.Read(buf[:])
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	var response api.BaseResponse
	err = json.Unmarshal(buf[:n], &response)
	if err != nil {
		log.Error("account=%s chain=%s nonce=%d cancelNonce error %s", nm.account.String(), nm.chainName, nonceUsed, err.Error())
		return
	}
	if response.GetErrorCode() != api.ErrorCodeSuccess {
		log.Error("account=%s chain=%s nonce=%d cancelNonce error : errorCode=%s errorMsg=%s", nm.account.String(), nm.chainName, nonceUsed, response.ErrorCode, response.ErrorMsg)
		return
	}
	nm.usedNonceToCancelURLMap.Delete(nonceUsed)
	log.Info("account=%s chain=%s nonce=%d cancel", nm.account.String(), nm.chainName, nonceUsed)
}
