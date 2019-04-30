package models

import (
	"fmt"
	"math/big"

	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

/*
LockoutInfo 状态变更:
以下所有事件均为确认后的事件,即会延后一定块处理 :
	1. 收到侧链PrepareLockout事件即SCPLO 		: MCLockStatus=LockStatusNone SCLockStatus=LockStatusLock
	2. 收到主链PrepareLockout事件即MCPLO 		: MCLockStatus=LockStatusLock SCLockStatus=LockStatusLock
	3. 收到主链LockoutSecret事件即MCLOS 		: MCLockStatus=LockStatusUnlock SCLockStatus=LockStatusLock
	4. 收到侧链Lockout事件即SCLO  		 	: MCLockStatus=LockStatusUnlock SCLockStatus=LockStatusUnlock 完结状态
	5. 收到主链CancelLockout事件即MCCancel	: MCLockStatus = LockStatusCancel
	6. 收到侧链CancelLockout事件即SCCancel	: SCLockStatus = LockStatusCancel
	7. 锁过期							 	: MCLockStatus/SCLockStatus = LockStatusExpiration
*/
/*
LockoutInfo :
该结构体记录一次lockout的所有数据及状态
*/
type LockoutInfo struct {
	MCChainName                string         `json:"mc_chain_name"`                   // 链名
	SecretHash                 common.Hash    `json:"secret_hash"`                     // 密码hash,db唯一ID
	Secret                     common.Hash    `json:"secret"`                          // 密码
	MCUserAddressHex           string         `json:"mc_user_address_hex"`             // 用户主链地址,格式根据链不同不同
	SCUserAddress              common.Address `json:"sc_user_address"`                 // 用户侧链地址,即在Spectrum上的地址
	SCTokenAddress             common.Address `json:"sc_token_address"`                // SCToken
	Amount                     *big.Int       `json:"amount"`                          // 金额
	MCExpiration               uint64         `json:"mc_expiration"`                   // 主链过期块
	SCExpiration               uint64         `json:"sc_expiration"`                   // 侧链过期块
	MCLockStatus               LockStatus     `json:"mc_lock_status"`                  // 主链锁状态
	SCLockStatus               LockStatus     `json:"sc_lock_status"`                  // 侧链锁状态
	Data                       []byte         `json:"data"`                            // 附加信息
	NotaryIDInCharge           int            `json:"notary_id_in_charge"`             // 负责公证人ID,如果没参与lockout签名的公证人,在整个LockoutInfo生命周期中,该值都为UnknownNotaryIdInCharge
	StartTime                  int64          `json:"start_time"`                      // 开始时间,即MCPrepareLockout事件发生的时间
	StartSCBlockNumber         uint64         `json:"start_sc_block_number"`           // 开始时侧链块数
	EndTime                    int64          `json:"end_time"`                        // 结束时间,即MCLockout事件发生的时间
	EndSCBlockNumber           uint64         `json:"end_sc_block_number"`             // 结束时侧链块数
	BTCPrepareLockoutTXHashHex string         `json:"btc_prepare_lockout_tx_hash_hex"` // 比特币锁定的utxo的txhash
	BTCPrepareLockoutVout      uint32         `json:"btc_prepare_lockout_vout"`        // 比特币锁定的utxo的index
	BTCLockScriptHex           string         `json:"btc_lock_script_hex"`             // 比特币锁定交易
	CrossFee                   *big.Int       `json:"cross_fee"`                       // 该字段记录本次lockin收取的手续费,即用户主链锁定的金额-lockin完成后侧链用户账户新增的金额
}

/*
IsEnd :
判断一次lockout流程是否已经完整结束
*/
func (l *LockoutInfo) IsEnd() bool {
	if l.MCLockStatus == LockStatusUnlock && l.SCLockStatus == LockStatusUnlock {
		//主侧链都已解锁
		return true
	}
	if l.MCLockStatus == LockStatusCancel && (l.SCLockStatus == LockStatusCancel || l.SCLockStatus == LockStatusNone) {
		//主侧链都已Cancel 或主链canncel,侧链没发生
		return true
	}
	return false
}

type lockoutInfoModel struct {
	MCChainName                string
	SecretHash                 []byte `gorm:"primary_key"`
	Secret                     []byte
	MCUserAddressHex           string
	SCUserAddress              []byte
	SCTokenAddress             []byte
	Amount                     []byte
	MCExpiration               uint64
	SCExpiration               uint64
	MCLockStatus               int
	SCLockStatus               int
	Data                       []byte
	NotaryIDInCharge           int
	StartTime                  int64
	StartMCBlockNumber         uint64
	EndTime                    int64
	EndMCBlockNumber           uint64
	BTCPrepareLockoutTXHashHex string // 比特币锁定的utxo的txhash
	BTCPrepareLockoutVout      uint32 // 比特币锁定的utxo的index
	BTCLockScriptHex           string //
	CrossFee                   []byte
}

func (lom *lockoutInfoModel) toLockoutInfo() *LockoutInfo {
	amount := new(big.Int)
	amount.SetBytes(lom.Amount)
	crossFee := new(big.Int)
	crossFee.SetBytes(lom.CrossFee)
	return &LockoutInfo{
		MCChainName:                lom.MCChainName,
		SecretHash:                 common.BytesToHash(lom.SecretHash),
		Secret:                     common.BytesToHash(lom.Secret),
		MCUserAddressHex:           lom.MCUserAddressHex,
		SCUserAddress:              common.BytesToAddress(lom.SCUserAddress),
		SCTokenAddress:             common.BytesToAddress(lom.SCTokenAddress),
		Amount:                     amount,
		MCExpiration:               lom.MCExpiration,
		SCExpiration:               lom.SCExpiration,
		MCLockStatus:               LockStatus(lom.MCLockStatus),
		SCLockStatus:               LockStatus(lom.SCLockStatus),
		Data:                       lom.Data,
		NotaryIDInCharge:           lom.NotaryIDInCharge,
		StartTime:                  lom.StartTime,
		StartSCBlockNumber:         lom.StartMCBlockNumber,
		EndTime:                    lom.EndTime,
		EndSCBlockNumber:           lom.EndMCBlockNumber,
		BTCPrepareLockoutTXHashHex: lom.BTCPrepareLockoutTXHashHex,
		BTCPrepareLockoutVout:      lom.BTCPrepareLockoutVout,
		BTCLockScriptHex:           lom.BTCLockScriptHex,
		CrossFee:                   crossFee,
	}
}

func (lom *lockoutInfoModel) fromLockoutInfo(l *LockoutInfo) *lockoutInfoModel {
	lom.MCChainName = l.MCChainName
	lom.SecretHash = l.SecretHash[:]
	lom.Secret = l.Secret[:]
	lom.MCUserAddressHex = l.MCUserAddressHex
	lom.SCUserAddress = l.SCUserAddress[:]
	lom.SCTokenAddress = l.SCTokenAddress[:]
	lom.Amount = l.Amount.Bytes()
	lom.MCExpiration = l.MCExpiration
	lom.SCExpiration = l.SCExpiration
	lom.MCLockStatus = int(l.MCLockStatus)
	lom.SCLockStatus = int(l.SCLockStatus)
	lom.Data = l.Data
	lom.NotaryIDInCharge = l.NotaryIDInCharge
	lom.StartTime = l.StartTime
	lom.StartMCBlockNumber = l.StartSCBlockNumber
	lom.EndTime = l.EndTime
	lom.EndMCBlockNumber = l.EndSCBlockNumber
	lom.BTCPrepareLockoutTXHashHex = l.BTCPrepareLockoutTXHashHex
	lom.BTCPrepareLockoutVout = l.BTCPrepareLockoutVout
	lom.BTCLockScriptHex = l.BTCLockScriptHex
	lom.CrossFee = l.CrossFee.Bytes()
	return lom
}

// NewLockoutInfo :
func (db *DB) NewLockoutInfo(lockoutInfo *LockoutInfo) (err error) {
	var t lockoutInfoModel
	return db.Create(t.fromLockoutInfo(lockoutInfo)).Error
}

// GetAllLockoutInfo :
func (db *DB) GetAllLockoutInfo() (list []*LockoutInfo, err error) {
	var t []lockoutInfoModel
	err = db.Find(&t).Error
	if err == storm.ErrNotFound {
		err = nil
		return
	}
	for _, l := range t {
		list = append(list, l.toLockoutInfo())
	}
	return
}

// GetAllLockoutInfoBySCToken :
func (db *DB) GetAllLockoutInfoBySCToken(scToken common.Address) (list []*LockoutInfo, err error) {
	var t []lockoutInfoModel
	err = db.Where(&lockoutInfoModel{
		SCTokenAddress: scToken[:],
	}).Find(&t).Error
	if err == storm.ErrNotFound {
		err = nil
		return
	}
	for _, l := range t {
		list = append(list, l.toLockoutInfo())
	}
	return
}

// GetLockoutInfo :
func (db *DB) GetLockoutInfo(secretHash common.Hash) (lockoutInfo *LockoutInfo, err error) {
	var lim lockoutInfoModel
	err = db.Where(&lockoutInfoModel{
		SecretHash: secretHash[:],
	}).First(&lim).Error
	if err != nil {
		return
	}
	lockoutInfo = lim.toLockoutInfo()
	return
}

// UpdateLockoutInfo :
func (db *DB) UpdateLockoutInfo(lockoutInfo *LockoutInfo) (err error) {
	var l *LockoutInfo
	l, err = db.GetLockoutInfo(lockoutInfo.SecretHash)
	if l == nil {
		err = fmt.Errorf("can not update LockoutInfo because key : secretHash=%s not found in db", lockoutInfo.SecretHash.String())
		return
	}
	var t lockoutInfoModel
	return db.Save(t.fromLockoutInfo(lockoutInfo)).Error
}
