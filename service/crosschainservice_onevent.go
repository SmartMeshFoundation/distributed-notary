package service

import (
	"errors"

	"fmt"

	"time"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chain/bitcoin"
	ethevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/nkbai/log"
)

// OnEvent 链上事件逻辑处理
func (cs *CrossChainService) OnEvent(e chain.Event) {
	var err error
	switch event := e.(type) {
	/*
		events about block number
	*/
	case ethevents.NewBlockEvent:
		err = cs.onMCNewBlockEvent(event.BlockNumber)
		if err != nil {
			log.Error(SCTokenLogMsg(cs.meta, "%s event deal err =%s", e.GetEventName(), err.Error()))
		}
		return
	case bitcoin.NewBlockEvent:
		err = cs.onMCNewBlockEvent(event.BlockNumber)
		if err != nil {
			log.Error(SCTokenLogMsg(cs.meta, "%s event deal err =%s", e.GetEventName(), err.Error()))
		}
		return
	case smcevents.NewBlockEvent:
		err = cs.onSCNewBlockEvent(event)
		if err != nil {
			log.Error(SCTokenLogMsg(cs.meta, "%s event deal err =%s", e.GetEventName(), err.Error()))
		}
		return
	/*
		events about lockin
	*/
	case ethevents.PrepareLockinEvent: // MCPLI
		log.Info(SCTokenLogMsg(cs.meta, "Receive MC PrepareLockinEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onMCPrepareLockin4Ethereum(event)
	case smcevents.PrepareLockinEvent: // SCPLI
		log.Info(SCTokenLogMsg(cs.meta, "Receive SC PrepareLockinEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onSCPrepareLockin(event)
	case smcevents.LockinSecretEvent: //  SCLIS
		log.Info(SCTokenLogMsg(cs.meta, "Receive SC LockinSecretEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onSCLockinSecret(event)
	case ethevents.LockinEvent: // MCLI
		log.Info(SCTokenLogMsg(cs.meta, "Receive MC LockinEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onMCLockin4Ethereum(event)
	case bitcoin.LockinEvent:
		log.Info(SCTokenLogMsg(cs.meta, "Receive MC LockinEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onMCLockin4Bitcoin(event)
	case ethevents.CancelLockinEvent: // MCCancelLI
		log.Info(SCTokenLogMsg(cs.meta, "Receive MC CancelLockinEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onMCCancelLockin4Ethereum(event)
	case smcevents.CancelLockinEvent: // SCCancelLI
		log.Info(SCTokenLogMsg(cs.meta, "Receive SC CancelLockinEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onSCCancelLockin(event)
	/*
		events about lockout
	*/
	case smcevents.PrepareLockoutEvent: // SCPLO
		log.Info(SCTokenLogMsg(cs.meta, "Receive SC PrepareLockoutEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onSCPrepareLockout(event)
	case ethevents.PrepareLockoutEvent: // MCPLO
		log.Info(SCTokenLogMsg(cs.meta, "Receive MC PrepareLockoutEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onMCPrepareLockout4Ethereum(event)
	case ethevents.LockoutSecretEvent: // MCLOS
		log.Info(SCTokenLogMsg(cs.meta, "Receive MC LockoutSecretEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onMCLockoutSecret4Ethereum(event)
	case smcevents.LockoutEvent: // SCLO
		log.Info(SCTokenLogMsg(cs.meta, "Receive SC LockoutEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onSCLockout(event)
	case ethevents.CancelLockoutEvent: // MCCancelLO
		log.Info(SCTokenLogMsg(cs.meta, "Receive MC CancelLockoutEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onMCCancelLockout4Ethereum(event)
	case smcevents.CancelLockoutEvent: // SCCancelLO
		log.Info(SCTokenLogMsg(cs.meta, "Receive SC CancelLockoutEvent :\n%s", utils.ToJSONStringFormat(event)))
		err = cs.onSCCancelLockout(event)
	default:
		err = errors.New("unknow event")
	}
	if err != nil {
		log.Error(SCTokenLogMsg(cs.meta, "%s event deal err =%s", e.GetEventName(), err.Error()))
	} else {
		log.Info(SCTokenLogMsg(cs.meta, "Event deal SUCCESS"))
	}
	return
}

func (cs *CrossChainService) onMCNewBlockEvent(blockNumber uint64) (err error) {
	cs.mcLastedBlockNumber = blockNumber
	// lockin
	err = cs.lockinHandler.onMCNewBlockEvent(blockNumber)
	if err != nil {
		return
	}
	// lockout
	var lockoutListNeedCancel []*models.LockoutInfo
	lockoutListNeedCancel, err = cs.lockoutHandler.onMCNewBlockEvent(blockNumber)
	if err != nil {
		return
	}
	if len(lockoutListNeedCancel) > 0 {
		for _, lockout := range lockoutListNeedCancel {
			if lockout.NotaryIDInCharge == cs.selfNotaryID {
				// 如果我是负责人,尽快cancel,这里如果调用合约出错,也继续去尝试cancel下一个
				err = cs.callMCCancelLockout(lockout.MCUserAddressHex)
				if err != nil {
					log.Error("callMCCancelLockout err = %s", err.Error())
				}
			}
		}
	}
	return
}

func (cs *CrossChainService) onSCNewBlockEvent(event smcevents.NewBlockEvent) (err error) {
	cs.scLastedBlockNumber = event.BlockNumber
	// lockin
	var lockinListNeedCancel []*models.LockinInfo
	lockinListNeedCancel, err = cs.lockinHandler.onSCNewBlockEvent(event.BlockNumber)
	if err != nil {
		return
	}
	if len(lockinListNeedCancel) > 0 {
		for _, lockin := range lockinListNeedCancel {
			if lockin.NotaryIDInCharge == cs.selfNotaryID {
				// 如果我是负责人,尽快cancel,这里如果调用合约出错,也继续去尝试cancel下一个
				err = cs.callSCCancelLockin(lockin.SCUserAddress.String())
				if err != nil {
					log.Error("callSCCancelLockin err = %s", err.Error())
				}
			}
		}
	}
	// lockout
	err = cs.lockoutHandler.onSCNewBlockEvent(event.BlockNumber)
	return
}

/*
主链PrepareLockin(MCPLI)
该事件为一个LockinInfo的生命周期开端,发生于用户调用时
事件为已确认事件,直接构造LockinInfo并保存,等待后续调用
*/
func (cs *CrossChainService) onMCPrepareLockin4Ethereum(event ethevents.PrepareLockinEvent) (err error) {
	// 1. 查询
	secretHash, mcExpiration, amount, err := cs.mcProxy.QueryLockin(event.Account.String())
	if err != nil {
		err = fmt.Errorf("mcProxy.QueryLockin err = %s", err.Error())
		return
	}
	// 1.5 校验mcExpiration
	if mcExpiration-event.BlockNumber <= params.MinLockinMCExpiration {
		err = fmt.Errorf("mcExpiration must bigger than %d", params.MinLockinMCExpiration)
		return
	}
	// 1.6 计算scExpiration
	scExpiration := cs.scLastedBlockNumber + (mcExpiration - event.BlockNumber - 5*params.ForkConfirmNumber - 1)

	// 2. 构造LockinInfo
	lockinInfo := &models.LockinInfo{
		MCChainName:      ethevents.ChainName,
		SecretHash:       secretHash,
		Secret:           utils.EmptyHash,
		MCUserAddressHex: event.Account.String(),
		SCUserAddress:    utils.EmptyAddress,
		SCTokenAddress:   cs.meta.SCToken,
		Amount:           amount,
		MCExpiration:     mcExpiration,
		SCExpiration:     scExpiration,
		MCLockStatus:     models.LockStatusLock,
		SCLockStatus:     models.LockStatusNone,
		//Data:               data,
		NotaryIDInCharge:   models.UnknownNotaryIDInCharge,
		StartTime:          event.Time.Unix(),
		StartMCBlockNumber: event.BlockNumber,
	}
	// 3. 调用handler处理
	err = cs.lockinHandler.registerLockin(lockinInfo)
	if err != nil {
		err = fmt.Errorf("lockinHandler.registerLockin err = %s", err.Error())
		return
	}
	return
}

/*
侧链PrepareLockin(SCPLI)
该事件的发起方为公证人,可能为自己
事件为已确认事件,修改LockinInfo状态
*/
func (cs *CrossChainService) onSCPrepareLockin(event smcevents.PrepareLockinEvent) (err error) {
	// 1. 查询
	secretHash, scExpiration, amount, err := cs.scTokenProxy.QueryLockin(event.Account.String())
	if err != nil {
		err = fmt.Errorf("scTokenProxy.QueryLockin err = %s", err.Error())
		return
	}
	// 2. 获取本地LockinInfo信息
	lockinInfo, err := cs.lockinHandler.getLockin(secretHash)
	if err != nil {
		err = fmt.Errorf("lockinHandler.getLockin err = %s", err.Error())
		return
	}
	// 3.　校验
	// 如果我参与了签名,那么关键数据的校验我在签名时就已经校验过了
	// 如果我没有参与签名,而合约不出错,如果收到与之前主链事件数据不匹配的SCPLI事件,此时事件已经发生,也没法挽回,所以还是记录数据
	if lockinInfo.MCLockStatus != models.LockStatusLock || lockinInfo.SCLockStatus != models.LockStatusNone || lockinInfo.Secret != utils.EmptyHash {
		log.Error("local lockinInfo status does't right,something must wrong, local lockinInfo:\n%s", utils.ToJSONStringFormat(lockinInfo))
	}
	if lockinInfo.Amount.Cmp(amount) != 0 {
		log.Error("amount does't match")
	}

	// 4. 修改状态,等待后续调用
	lockinInfo.SCExpiration = scExpiration // 存在因为各节点区块高度细微差距导致的自己之前计算的SCExpiration不对,这里取合约里面的真实值
	lockinInfo.SCLockStatus = models.LockStatusLock
	lockinInfo.SCUserAddress = event.Account
	err = cs.lockinHandler.updateLockin(lockinInfo)
	if err != nil {
		err = fmt.Errorf("lockinHandler.UpdateLockinInfo err = %s", err.Error())
		return
	}
	return
}

/*
侧链LockinSecret(SCLIS)
该事件由用户发起,已确认
*/
func (cs *CrossChainService) onSCLockinSecret(event smcevents.LockinSecretEvent) (err error) {
	// 1.根据密码hash查询LockinInfo
	secretHash := utils.ShaSecret(event.Secret[:])
	lockinInfo, err := cs.lockinHandler.getLockin(secretHash)
	if err != nil {
		err = fmt.Errorf("lockinHandler.getLockin err = %s", err.Error())
		return
	}
	// 2. 重复校验 TODO
	if lockinInfo.Secret != utils.EmptyHash {
		err = fmt.Errorf("receive repeat SCLockinSecret, ignore")
		return
	}
	// 3. 校验状态,好像没啥用,用户都拿走钱了,就算状态不对,也需要继续操作,让负责人尝试去主链lockin
	if lockinInfo.MCLockStatus != models.LockStatusLock || lockinInfo.SCLockStatus != models.LockStatusLock {
		log.Error("local lockinInfo status does't right,something must wrong, local lockinInfo:\n%s", utils.ToJSONStringFormat(lockinInfo))
	}
	// 3. 更新状态
	lockinInfo.Secret = event.Secret
	lockinInfo.SCLockStatus = models.LockStatusUnlock
	err = cs.lockinHandler.updateLockin(lockinInfo)
	if err != nil {
		err = fmt.Errorf("lockinHandler.updateLockin err = %s", err.Error())
		return
	}
	// 4. 如果自己是负责人,发起主链Lockin
	if lockinInfo.NotaryIDInCharge == cs.selfNotaryID {
		err = cs.callMCLockin(lockinInfo)
		if err != nil {
			err = fmt.Errorf("callMCLockin err = %s", err.Error())
			return
		}
	}
	return
}

/*
主链Lockin
收到该事件,说明一次Lockin完整结束,合约上已经清楚该UserAccount的lockin信息
*/
func (cs *CrossChainService) onMCLockin4Ethereum(event ethevents.LockinEvent) (err error) {
	// 1. 获取本地LockinInfo信息
	lockinInfo, err := cs.lockinHandler.getLockin(event.SecretHash)
	if err != nil {
		err = fmt.Errorf("lockinHandler.getLockin err = %s", err.Error())
		return
	}
	// 2. 校验 TODO
	if lockinInfo.MCUserAddressHex != event.Account.String() {
		err = fmt.Errorf("MCUserAddressHex does't match")
		return
	}
	// 3. 更新本地信息
	lockinInfo.MCLockStatus = models.LockStatusUnlock
	lockinInfo.EndTime = event.Time.Unix()
	lockinInfo.EndMCBlockNumber = event.BlockNumber
	err = cs.lockinHandler.updateLockin(lockinInfo)
	if err != nil {
		err = fmt.Errorf("lockinHandler.UpdateLockinInfo err = %s", err.Error())
		return
	}
	return
}

/*
收到该事件,说明一次Lockin完整结束,合约上已经清楚该UserAccount的lockin信息
*/
func (cs *CrossChainService) onMCLockin4Bitcoin(event bitcoin.LockinEvent) (err error) {
	// 1. 获取本地LockinInfo信息
	lockinInfo, err := cs.lockinHandler.getLockin(event.SecretHash)
	if err != nil {
		err = fmt.Errorf("lockinHandler.getLockin err = %s", err.Error())
		return
	}
	// 2. 更新本地信息
	lockinInfo.MCLockStatus = models.LockStatusUnlock
	lockinInfo.EndTime = event.Time.Unix()
	lockinInfo.EndMCBlockNumber = event.BlockNumber
	err = cs.lockinHandler.updateLockin(lockinInfo)
	if err != nil {
		err = fmt.Errorf("lockinHandler.UpdateLockinInfo err = %s", err.Error())
		return
	}
	// 3. 注册新的普通outpoint
	var txHash chainhash.Hash
	err = chainhash.Decode(&txHash, event.TxHashStr)
	if err != nil {
		err = fmt.Errorf("chainhash.Decode err = %s", err.Error())
		return
	}
	for index, outpoint := range event.TxOuts {
		log.Trace("receive new utxo on BTC :[txHash=%s index=%d amount=%d]", event.TxHashStr, index, outpoint.Value)
		err = cs.lockinHandler.db.NewBTCOutpoint(&models.BTCOutpoint{
			PublicKeyHashStr: cs.meta.MCLockedPublicKeyHashStr,
			TxHashStr:        event.TxHashStr,
			Index:            index,
			Amount:           btcutil.Amount(outpoint.Value),
			Status:           models.BTCOutpointStatusUsable,
			CreateTime:       time.Now().Unix(),
		})
		if err != nil {
			// 这里应当如何处理比较好???
			log.Error(err.Error())
		}
		// 0. 获取bs,注册outpoint监听
		c, err := cs.dispatchService.getChainByName(bitcoin.ChainName)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		bs := c.(*bitcoin.BTCService)
		err = bs.RegisterP2PKHOutpoint(wire.OutPoint{
			Hash:  txHash,
			Index: uint32(index),
		}, cs.meta.MCLockedPublicKeyHashStr)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return
}

/*
主链取消
*/
func (cs *CrossChainService) onMCCancelLockin4Ethereum(event ethevents.CancelLockinEvent) (err error) {
	// 1. 获取本地LockinInfo信息
	lockinInfo, err := cs.lockinHandler.getLockin(event.SecretHash)
	if err != nil {
		err = fmt.Errorf("lockinHandler.getLockin err = %s", err.Error())
		return
	}
	// 2. 校验 TODO
	if lockinInfo.MCUserAddressHex != event.Account.String() {
		err = fmt.Errorf("MCUserAddressHex does't match")
		return
	}
	// 3. 更新本地信息,endTime哪个链取消在后面取哪个
	lockinInfo.MCLockStatus = models.LockStatusCancel
	lockinInfo.EndTime = event.Time.Unix()
	lockinInfo.EndMCBlockNumber = event.BlockNumber
	err = cs.lockinHandler.updateLockin(lockinInfo)
	if err != nil {
		err = fmt.Errorf("lockinHandler.updateLockin err = %s", err.Error())
		return
	}
	return
}

/*
侧链取消
*/
func (cs *CrossChainService) onSCCancelLockin(event smcevents.CancelLockinEvent) (err error) {
	// 1. 获取本地LockinInfo信息
	lockinInfo, err := cs.lockinHandler.getLockin(event.SecretHash)
	if err != nil {
		err = fmt.Errorf("lockinHandler.getLockin err = %s", err.Error())
		return
	}
	// 2. 校验 TODO
	if lockinInfo.SCUserAddress != event.Account {
		err = fmt.Errorf("SCUserAddress does't match")
		return
	}
	// 3. 更新本地信息,endTime哪个链取消在后面取哪个
	lockinInfo.SCLockStatus = models.LockStatusCancel
	lockinInfo.EndTime = event.Time.Unix()
	err = cs.lockinHandler.updateLockin(lockinInfo)
	if err != nil {
		err = fmt.Errorf("lockinHandler.updateLockin err = %s", err.Error())
		return
	}
	return
}

/*
	lockout 相关事件处理
*/

func (cs *CrossChainService) onSCPrepareLockout(event smcevents.PrepareLockoutEvent) (err error) {
	// 1. 查询
	secretHash, scExpiration, amount, err := cs.scTokenProxy.QueryLockout(event.Account.String())
	if err != nil {
		err = fmt.Errorf("scTokenProxy.QueryLockout err = %s", err.Error())
		return
	}
	// 1.5 校验scExpiration
	if scExpiration-event.BlockNumber <= params.MinLockoutSCExpiration {
		err = fmt.Errorf("scExpiration must bigger than %d", params.MinLockoutSCExpiration)
		return
	}
	// 1.6 计算mcExpiration
	mcExpiration := cs.mcLastedBlockNumber + (scExpiration - event.BlockNumber - 5*params.ForkConfirmNumber - 1)

	// 2. 构造LockoutInfo
	lockoutInfo := &models.LockoutInfo{
		SecretHash:       secretHash,
		Secret:           utils.EmptyHash,
		MCUserAddressHex: "",
		SCUserAddress:    event.Account,
		SCTokenAddress:   cs.meta.SCToken,
		Amount:           amount,
		MCExpiration:     mcExpiration,
		SCExpiration:     scExpiration,
		MCLockStatus:     models.LockStatusNone,
		SCLockStatus:     models.LockStatusLock,
		//Data:               data,
		NotaryIDInCharge:   models.UnknownNotaryIDInCharge,
		StartTime:          event.Time.Unix(),
		StartSCBlockNumber: event.BlockNumber,
	}
	// 3. 调用handler处理
	err = cs.lockoutHandler.registerLockout(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.registerLockout err = %s", err.Error())
		return
	}
	return
}

/*
主链PrepareLockout(MCPLO)
该事件的发起方为公证人,可能为自己
事件为已确认事件,修改LockoutInfo状态
*/
func (cs *CrossChainService) onMCPrepareLockout4Ethereum(event ethevents.PrepareLockoutEvent) (err error) {
	// 1. 查询
	secretHash, mcExpiration, amount, err := cs.mcProxy.QueryLockout(event.Account.String())
	if err != nil {
		err = fmt.Errorf("mcProxy.QueryLockout err = %s", err.Error())
		return
	}
	// 2. 获取本地LockoutInfo信息
	lockoutInfo, err := cs.lockoutHandler.getLockout(secretHash)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.getLockout err = %s", err.Error())
		return
	}
	// 3.　校验 TODO
	if lockoutInfo.MCLockStatus != models.LockStatusNone || lockoutInfo.SCLockStatus != models.LockStatusLock || lockoutInfo.Secret != utils.EmptyHash {
		err = fmt.Errorf("local lockoutInfo status does't right,something must wrong, local lockoutInfo:\n%s", utils.ToJSONStringFormat(lockoutInfo))
		return
	}
	if lockoutInfo.Amount.Cmp(amount) != 0 {
		err = fmt.Errorf("amount does't match")
		return
	}

	// 4. 修改状态,等待后续调用
	lockoutInfo.MCExpiration = mcExpiration // 这里由于自己本地块号和该笔MCPrepareLockout交易发起公证人的块号有些微差距,取合约上的值
	lockoutInfo.MCLockStatus = models.LockStatusLock
	lockoutInfo.MCUserAddressHex = event.Account.String()
	err = cs.lockoutHandler.updateLockout(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.updateLockout err = %s", err.Error())
		return
	}
	return
}

/*
主链LockoutSecret(MCLOS)
该事件由用户发起,已确认
*/
func (cs *CrossChainService) onMCLockoutSecret4Ethereum(event ethevents.LockoutSecretEvent) (err error) {
	// 1.根据密码hash查询LockoutInfo
	secretHash := utils.ShaSecret(event.Secret[:])
	lockoutInfo, err := cs.lockoutHandler.getLockout(secretHash)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.getLockout err = %s", err.Error())
		return
	}
	// 2. 重复校验 TODO
	if lockoutInfo.Secret != utils.EmptyHash {
		err = fmt.Errorf("receive repeat MCLockoutSecret, ignore")
		return
	}
	// 3. 校验状态,好像没啥用,用户都拿走钱了,就算状态不对,也需要继续操作,让负责人尝试去侧链lockout
	if lockoutInfo.MCLockStatus != models.LockStatusLock || lockoutInfo.SCLockStatus != models.LockStatusLock {
		err = fmt.Errorf("local lockoutInfo status does't right,something must wrong, local lockoutInfo:\n%s", utils.ToJSONStringFormat(lockoutInfo))
	}
	// 3. 更新状态
	lockoutInfo.Secret = event.Secret
	lockoutInfo.MCLockStatus = models.LockStatusUnlock
	err = cs.lockoutHandler.updateLockout(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.updateLockout err = %s", err.Error())
		return
	}
	// 4. 如果自己是负责人,发起侧链Lockout
	if lockoutInfo.NotaryIDInCharge == cs.selfNotaryID {
		err = cs.callSCLockout(lockoutInfo.SCUserAddress.String(), lockoutInfo.Secret)
		if err != nil {
			err = fmt.Errorf("callSCLockout err = %s", err.Error())
			return
		}
	}
	return
}

/*
侧链链Lockout
收到该事件,说明一次Lockout完整结束,合约上已经清楚该UserAccount的lockout信息
*/
func (cs *CrossChainService) onSCLockout(event smcevents.LockoutEvent) (err error) {
	// 1. 获取本地LockoutInfo信息
	lockoutInfo, err := cs.lockoutHandler.getLockout(event.SecretHash)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.getLockout err = %s", err.Error())
		return
	}
	// 2. 校验 TODO
	if lockoutInfo.SCUserAddress != event.Account {
		err = fmt.Errorf("SCUserAddress does't match")
		return
	}
	// 3. 更新本地信息
	lockoutInfo.SCLockStatus = models.LockStatusUnlock
	lockoutInfo.EndTime = event.Time.Unix()
	lockoutInfo.EndSCBlockNumber = event.BlockNumber
	err = cs.lockoutHandler.updateLockout(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.updateLockout err = %s", err.Error())
		return
	}
	return
}

/*
主链取消
*/
func (cs *CrossChainService) onMCCancelLockout4Ethereum(event ethevents.CancelLockoutEvent) (err error) {
	// 1. 获取本地LockoutInfo信息
	lockoutInfo, err := cs.lockoutHandler.getLockout(event.SecretHash)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.getLockout err = %s", err.Error())
		return
	}
	// 2. 校验 TODO
	if lockoutInfo.MCUserAddressHex != event.Account.String() {
		err = fmt.Errorf("MCUserAddressHex does't match")
		return
	}
	// 3. 更新本地信息,endTime哪个链取消在后面取哪个
	lockoutInfo.MCLockStatus = models.LockStatusCancel
	lockoutInfo.EndTime = event.Time.Unix()
	err = cs.lockoutHandler.updateLockout(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.updateLockout err = %s", err.Error())
		return
	}
	return
}

/*
侧链取消
*/
func (cs *CrossChainService) onSCCancelLockout(event smcevents.CancelLockoutEvent) (err error) {
	// 1. 获取本地LockoutInfo信息
	lockoutInfo, err := cs.lockoutHandler.getLockout(event.SecretHash)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.getLockout err = %s", err.Error())
		return
	}
	// 2. 校验 TODO
	if lockoutInfo.SCUserAddress != event.Account {
		err = fmt.Errorf("SCUserAddress does't match")
		return
	}
	// 3. 更新本地信息,endTime哪个链取消在后面取哪个
	lockoutInfo.SCLockStatus = models.LockStatusCancel
	lockoutInfo.EndTime = event.Time.Unix()
	lockoutInfo.EndSCBlockNumber = event.BlockNumber
	err = cs.lockoutHandler.updateLockout(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("lockoutHandler.updateLockout err = %s", err.Error())
		return
	}
	return
}
