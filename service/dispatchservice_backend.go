package service

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	spectrumevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/nkbai/log"
)

/*
	其他service回调DispatchService的入口
*/
type dispatchServiceBackend interface {
	getChainByName(chainName string) (c chain.Chain, err error)

	/*
		notaryService在部署合约之后调用,原则上除此和启动时,其余地方不能调用
	*/
	registerNewSCToken(scTokenMetaInfo *models.SideChainTokenMetaInfo) (err error)
}

func (ds *DispatchService) getChainByName(chainName string) (c chain.Chain, err error) {
	var ok bool
	c, ok = ds.chainMap[chainName]
	if !ok {
		err = fmt.Errorf("can not find chain %s,something must wrong", chainName)
		return
	}
	return
}

func (ds *DispatchService) registerNewSCToken(scTokenMetaInfo *models.SideChainTokenMetaInfo) (err error) {
	// 6. 构造CrossChainService开始提供服务
	ds.scToken2CrossChainServiceMapLock.Lock()
	defer ds.scToken2CrossChainServiceMapLock.Unlock()
	ds.scToken2CrossChainServiceMap[scTokenMetaInfo.SCToken] = NewCrossChainService(ds.db, ds.notaryService.self, scTokenMetaInfo)
	// 注册侧链合约:
	err = ds.chainMap[spectrumevents.ChainName].RegisterEventListenContract(scTokenMetaInfo.SCToken)
	if err != nil {
		log.Error("RegisterEventListenContract on side chain err : %s", err.Error())
		return
	}
	// 注册主链合约:
	mc, ok := ds.chainMap[scTokenMetaInfo.MCName]
	if !ok {
		log.Error("can not find chain %s,something must wrong", scTokenMetaInfo.MCName)
		return
	}
	err = mc.RegisterEventListenContract(scTokenMetaInfo.MCLockedContractAddress)
	if err != nil {
		log.Error("RegisterEventListenContract on main chain %s err : %s", scTokenMetaInfo.MCName, err.Error())
		return
	}
	return
}
