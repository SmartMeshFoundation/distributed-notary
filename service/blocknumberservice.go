package service

import (
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/nkbai/log"
)

/*
BlockNumberService 保存所有链的最新高度
*/
type BlockNumberService struct {
	db    *models.DB
	m     map[string]uint64
	mLock sync.Mutex
}

// NewBlockNumberService :启动时初始化
func NewBlockNumberService(db *models.DB, chainMap map[string]chain.Chain) (bm *BlockNumberService, err error) {
	bm = &BlockNumberService{
		db: db,
		m:  make(map[string]uint64),
	}
	// 初始化m,同时设置各个Chain的LastBlockNumber
	bm.mLock.Lock()
	for _, c := range chainMap {
		lastBlockNumber := db.GetLastBlockNumber(c.GetChainName())
		c.SetLastBlockNumber(lastBlockNumber)
		bm.m[c.GetChainName()] = lastBlockNumber
	}
	bm.mLock.Unlock()
	return
}

// GetLastBlockNumber 获取一条链的最新块号
func (bm *BlockNumberService) GetLastBlockNumber(chainName string) uint64 {
	bm.mLock.Lock()
	if n, ok := bm.m[chainName]; ok {
		bm.mLock.Unlock()
		return n
	}
	bm.mLock.Unlock()
	return 0
}

// NewBlockNumber :
func (bm *BlockNumberService) NewBlockNumber(event chain.Event) {
	chainName := event.GetChainName()
	lastBlockNumber := event.GetBlockNumber()
	bm.mLock.Lock()
	err := bm.db.SaveLastBlockNumber(chainName, lastBlockNumber)
	if err != nil {
		log.Error("db.SaveLastBlockNumber err = %s, something must wrong", err.Error())
		return
	}
	bm.m[chainName] = lastBlockNumber
	bm.mLock.Unlock()
	return
}
