package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"

	"github.com/SmartMeshFoundation/distributed-notary/pbft/pbft"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/cmd/nonce_server/nonceapi"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	spectrumevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
	其他service回调DispatchService的入口
*/
type dispatchServiceBackend interface {
	getBtcNetworkParam() *chaincfg.Params
	getSelfPrivateKey() *ecdsa.PrivateKey
	getSelfNotaryInfo() *models.NotaryInfo
	getChainByName(chainName string) (c chain.Chain, err error)
	getLockinInfo(scTokenAddress common.Address, secretHash common.Hash) (lockinInfo *models.LockinInfo, err error)
	getLockInInfoBySCPrepareLockInRequest(req *userapi.SCPrepareLockinRequest) (lockinInfo *models.LockinInfo, err error)
	getLockoutInfo(scTokenAddress common.Address, secretHash common.Hash) (lockoutInfo *models.LockoutInfo, err error)
	getNotaryService() *NotaryService
	getSCTokenMetaInfoBySCTokenAddress(scTokenAddress common.Address) (scToken *models.SideChainTokenMetaInfo)
	/*
		发送http请求给nonce-sever,调用/api/1/apply-nonce接口申请可用某个账号的可用nonce,合约调用之前使用
	*/
	applyNonceFromNonceServer(chainName string, priveKeyID common.Hash, reason string, amount *big.Int) (nonce uint64, err error)
	/*
		让PBFT协商分配合适的UTXO用于prepareLockout
	*/
	applyUTXO(chainName string, priveKeyID common.Hash, reason string, amount *big.Int) (utxos string, err error)

	/*
		notaryService在部署合约之后调用,原则上除此和启动时,其余地方不能调用
	*/
	registerNewSCToken(scTokenMetaInfo *models.SideChainTokenMetaInfo) (err error)

	/*
		notaryService在协商调用合约之后,更新lockinInfo中的NotaryIDInCharge字段,其余地方不应该调用
	*/
	updateLockinInfoNotaryIDInChargeID(scTokenAddress common.Address, secretHash common.Hash, notaryID int) (err error)
	/*
		notaryService在协商调用合约之后,更新lockinInfo中的NotaryIDInCharge字段,其余地方不应该调用
	*/
	updateLockoutInfoNotaryIDInChargeID(scTokenAddress common.Address, secretHash common.Hash, notaryID int) (err error)
}

func (ds *DispatchService) getBtcNetworkParam() *chaincfg.Params {
	c, err := ds.getChainByName(bitcoin.ChainName)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	c2, _ := c.(*bitcoin.BTCService)
	return c2.GetNetParam()
}
func (ds *DispatchService) getSelfPrivateKey() *ecdsa.PrivateKey {
	return ds.notaryService.privateKey
}

func (ds *DispatchService) getSelfNotaryInfo() *models.NotaryInfo {
	return ds.notaryService.self
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

func (ds *DispatchService) getNotaryService() *NotaryService {
	return ds.notaryService
}

func (ds *DispatchService) getLockInInfoBySCPrepareLockInRequest(req *userapi.SCPrepareLockinRequest) (lockinInfo *models.LockinInfo, err error) {
	ds.scToken2CrossChainServiceMapLock.Lock()
	defer ds.scToken2CrossChainServiceMapLock.Unlock()
	cs, ok := ds.scToken2CrossChainServiceMap[req.SCTokenAddress]
	if !ok {
		panic("never happen")
	}
	return cs.getLockInInfoBySCPrepareLockInRequest(req)
}

func (ds *DispatchService) getLockinInfo(scTokenAddress common.Address, secretHash common.Hash) (lockinInfo *models.LockinInfo, err error) {
	ds.scToken2CrossChainServiceMapLock.Lock()
	defer ds.scToken2CrossChainServiceMapLock.Unlock()
	cs, ok := ds.scToken2CrossChainServiceMap[scTokenAddress]
	if !ok {
		panic("never happen")
	}
	return cs.lockinHandler.getLockin(secretHash)
}

func (ds *DispatchService) getLockoutInfo(scTokenAddress common.Address, secretHash common.Hash) (lockoutInfo *models.LockoutInfo, err error) {
	ds.scToken2CrossChainServiceMapLock.Lock()
	defer ds.scToken2CrossChainServiceMapLock.Unlock()
	cs, ok := ds.scToken2CrossChainServiceMap[scTokenAddress]
	if !ok {
		panic("never happen")
	}
	return cs.lockoutHandler.getLockout(secretHash)
}

func (ds *DispatchService) getSCTokenMetaInfoBySCTokenAddress(scTokenAddress common.Address) (scToken *models.SideChainTokenMetaInfo) {
	ds.scToken2CrossChainServiceMapLock.Lock()
	cs, ok := ds.scToken2CrossChainServiceMap[scTokenAddress]
	if !ok {
		panic("never happen")
	}
	ds.scToken2CrossChainServiceMapLock.Unlock()
	scToken = cs.meta
	return
}

func (ds *DispatchService) registerNewSCToken(scTokenMetaInfo *models.SideChainTokenMetaInfo) (err error) {
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
	// 6. 构造CrossChainService开始提供服务
	ds.scToken2CrossChainServiceMapLock.Lock()
	ds.scToken2CrossChainServiceMap[scTokenMetaInfo.SCToken] = NewCrossChainService(ds.db, ds, scTokenMetaInfo)
	ds.scToken2CrossChainServiceMapLock.Unlock()
	return
}

func (ds *DispatchService) updateLockinInfoNotaryIDInChargeID(scTokenAddress common.Address, secretHash common.Hash, notaryID int) (err error) {
	ds.scToken2CrossChainServiceMapLock.Lock()
	lh := ds.scToken2CrossChainServiceMap[scTokenAddress].lockinHandler
	ds.scToken2CrossChainServiceMapLock.Unlock()
	lockinInfo, err := lh.getLockin(secretHash)
	if err != nil {
		return
	}
	lockinInfo.NotaryIDInCharge = notaryID
	return lh.updateLockin(lockinInfo)
}
func (ds *DispatchService) updateLockoutInfoNotaryIDInChargeID(scTokenAddress common.Address, secretHash common.Hash, notaryID int) (err error) {
	ds.scToken2CrossChainServiceMapLock.Lock()
	lh := ds.scToken2CrossChainServiceMap[scTokenAddress].lockoutHandler
	ds.scToken2CrossChainServiceMapLock.Unlock()
	lockinInfo, err := lh.getLockout(secretHash)
	if err != nil {
		return
	}
	lockinInfo.NotaryIDInCharge = notaryID
	return lh.updateLockout(lockinInfo)
}

var debugpbft = false

//纯粹试了测试需要,避开pbft, 走传统的nonce server
func (ds *DispatchService) applyNonceFromNonceServerFake(chainName string, privKeyID common.Hash, reason string) (nonce uint64, err error) {
	pk, err := ds.db.LoadPrivateKeyInfo(privKeyID)
	if err != nil {
		panic(err)
	}
	account := pk.ToAddress()
	url := ds.nonceServerHost + nonceapi.APIName2URLMap[nonceapi.APINameApplyNonce]
	req := nonceapi.NewApplyNonceReq(chainName, account, "http://"+ds.userAPI.IPPort+userapi.APIName2URLMap[userapi.APIAdminNameCancelNonce])
	payload, err := json.Marshal(req)
	if err != nil {
		return
	}
	/* #nosec */
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
	var buf [4096 * 1024]byte
	n := 0
	n, err = resp.Body.Read(buf[:])
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	var response api.BaseResponse
	var applyNonceResponse nonceapi.ApplyNonceResponse
	err = json.Unmarshal(buf[:n], &response)
	if err != nil {
		return
	}
	err = response.ParseData(&applyNonceResponse)
	if err != nil {
		return
	}
	nonce = applyNonceResponse.Nonce
	return
}
func (ds *DispatchService) applyNonceFromNonceServer(chainName string, privKeyID common.Hash, reason string, amount *big.Int) (nonce uint64, err error) {
	if debugpbft {
		return ds.applyNonceFromNonceServerFake(chainName, privKeyID, reason)
	}
	key := fmt.Sprintf("%s-%s", chainName, privKeyID.String())
	ps, err := ds.getPbftService(key)
	if err != nil {
		return
	}
	return ps.(*PBFTService).newNonce(fmt.Sprintf("%s-%s-%s", chainName, reason, amount))
}
func (ds *DispatchService) applyUTXO(chainName string, priveKeyID common.Hash, reason string, amount *big.Int) (utxos string, err error) {
	key := fmt.Sprintf("%s-%s", chainName, priveKeyID.String())
	ps, err := ds.getPbftService(key)
	if err != nil {
		return
	}
	return ps.(*btcPBFTService).newUTXO(fmt.Sprintf("%s-%s-%s", chainName, reason, amount))
}
func (ds *DispatchService) getPbftService(key string) (ps pbft.PBFTAuxiliary, err error) {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	ps, ok := ds.pbftServices[key]
	if ok {
		return
	}
	/*
		需要创建pbftservice的时候需要验证
		chain存在
		privkey存在
	*/
	ss := strings.Split(key, "-")
	if len(ss) != 2 {
		err = fmt.Errorf("key err =%s", key)
		return
	}
	chainName := ss[0]
	privKeyID := common.HexToHash(ss[1])
	_, chainExist := ds.chainMap[chainName]
	if !chainExist {
		err = fmt.Errorf("chain %s unkown", chainName)
		return
	}
	_, err = ds.db.LoadPrivateKeyInfo(privKeyID)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("applyNonceFromNonceServer new pbft Service chainName=%s,privatekeyID=%s",
		chainName, privKeyID.String(),
	))
	typ := chainName2PBFTType(chainName)
	ps2 := NewPBFTService(key, chainName, privKeyID.String(),
		ds.notaries,
		ds.notaryService.notaryClient,
		ds, ds.db)
	switch typ {
	case pbftTypeEthereum:
		ds.pbftServices[key] = ps2
		err = ps2.server.UpdateAS(ps2)
		if err != nil {
			panic(err)
		}
	case pbftTypeBTC:
		ps3 := NewBTCPBFTService(ps2)
		ds.pbftServices[key] = ps3
		err = ps3.server.UpdateAS(ps3)
		if err != nil {
			panic(err)
		}
	default:
		return nil, errors.New("unkown chain")
	}
	ps = ds.pbftServices[key]
	return
}

func chainName2PBFTType(chainName string) pbftType {
	switch chainName {
	case spectrumevents.ChainName:
		return pbftTypeEthereum
	case events.ChainName:
		return pbftTypeEthereum
	case bitcoin.ChainName:
		return pbftTypeBTC
	}
	return pbftTypeUnkown
}
