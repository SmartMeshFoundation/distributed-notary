package models

import (
	"fmt"
	"math/big"

	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

/*
	注释缩写:
	MC = MainChain
	SC = SideChain
	PLI = PrepareLockin
	LIS = LockinSecret
	LOS = LockoutSecret
	MCPLI = MainChain PrepareLockin
	SCPLI = SideChain PrepareLockin
*/

// LockStatus :
type LockStatus int

/* #nosec */
const (
	LockStatusNone = iota
	LockStatusLock
	LockStatusUnlock
	LockStatusExpiration
	LockStatusCancel
)

/*
LockinInfo 状态变更:
以下所有事件均为确认后的事件,即会延后一定块处理 :
	1. 收到主链PrepareLockin事件即MCPLI 		: MCLockStatus=LockStatusLock SCLockStatus=LockStatusNone
	2. 收到侧链PrepareLockin事件即SCPLI 		: MCLockStatus=LockStatusLock SCLockStatus=LockStatusLock
	3. 收到侧链LockinSecretS事件即SCLIS 		: MCLockStatus=LockStatusLock SCLockStatus=LockStatusUnlock
	4. 收到主链Lockin事件即MCLI  		 		: MCLockStatus=LockStatusUnlock SCLockStatus=LockStatusUnlock 完结状态
	5. 收到主链CancelLockin事件即MCCancel		: MCLockStatus = LockStatusCancel
	6. 收到侧链CancelLockin事件即SCCancel		: SCLockStatus = LockStatusCancel
	7. 锁过期							 	: MCLockStatus/SCLockStatus = LockStatusExpiration
*/

// UnknownNotaryIDInCharge 为确定负责公证人之前的默认值
var UnknownNotaryIDInCharge = -1

/*
LockinInfo :
该结构体记录一次Lockin的所有数据及状态
*/
type LockinInfo struct {
	SecretHash         common.Hash    `json:"secret_hash"`         // 密码hash,db唯一ID
	Secret             common.Hash    `json:"secret"`              // 密码
	MCUserAddressHex   string         `json:"mc_user_address_hex"` // 用户主链地址,格式根据链不同不同
	SCUserAddress      common.Address `json:"sc_user_address"`     // 用户侧链地址,即在Spectrum上的地址
	SCTokenAddress     common.Address `json:"sc_token_address"`    // SCToken
	Amount             *big.Int       `json:"amount"`              // 金额
	MCExpiration       uint64         `json:"mc_expiration"`       // 主链过期块
	SCExpiration       uint64         `json:"sc_expiration"`       // 侧链过期块
	MCLockStatus       LockStatus     `json:"mc_lock_status"`      // 主链锁状态
	SCLockStatus       LockStatus     `json:"sc_lock_status"`      // 侧链锁状态
	Data               []byte         `json:"data"`                // 附加信息
	NotaryIDInCharge   int            `json:"notary_id_in_charge"` // 负责公证人ID,如果没参与lockin签名的公证人,在整个LockinInfo生命周期中,该值都为UnknownNotaryIdInCharge
	StartTime          int64          `json:"start_time"`          // 开始时间,即MCPrepareLockin事件发生的时间
	StartMCBlockNumber uint64         `json:"start_block_number"`  // 开始时主链块数
	EndTime            int64          `json:"end_time"`            // 结束时间,即MCLockin事件发生的时间
	EndMCBlockNumber   uint64         `json:"end_mc_block_number"` // 结束时主链块数
}

/*
IsEnd :
判断一次lockin流程是否已经完整结束
*/
func (l *LockinInfo) IsEnd() bool {
	if l.MCLockStatus == LockStatusUnlock && l.SCLockStatus == LockStatusUnlock {
		//主侧链都已解锁
		return true
	}
	if l.MCLockStatus == LockStatusCancel && l.SCLockStatus == LockStatusCancel {
		//主侧链都已Cancel
		return true
	}
	return false
}

type lockinInfoModel struct {
	SecretHash         []byte `gorm:"primary_key"`
	Secret             []byte
	MCUserAddressHex   string
	SCUserAddress      []byte
	SCTokenAddress     []byte
	Amount             []byte
	MCExpiration       uint64
	SCExpiration       uint64
	MCLockStatus       int
	SCLockStatus       int
	Data               []byte
	NotaryIDInCharge   int
	StartTime          int64
	StartMCBlockNumber uint64
	EndTime            int64
	EndMCBlockNumber   uint64
}

func (lim *lockinInfoModel) toLockinInfo() *LockinInfo {
	amount := new(big.Int)
	amount.SetBytes(lim.Amount)
	return &LockinInfo{
		SecretHash:         common.BytesToHash(lim.SecretHash),
		Secret:             common.BytesToHash(lim.Secret),
		MCUserAddressHex:   lim.MCUserAddressHex,
		SCUserAddress:      common.BytesToAddress(lim.SCUserAddress),
		SCTokenAddress:     common.BytesToAddress(lim.SCTokenAddress),
		Amount:             amount,
		MCExpiration:       lim.MCExpiration,
		SCExpiration:       lim.SCExpiration,
		MCLockStatus:       LockStatus(lim.MCLockStatus),
		SCLockStatus:       LockStatus(lim.SCLockStatus),
		Data:               lim.Data,
		NotaryIDInCharge:   lim.NotaryIDInCharge,
		StartTime:          lim.StartTime,
		StartMCBlockNumber: lim.StartMCBlockNumber,
		EndTime:            lim.EndTime,
		EndMCBlockNumber:   lim.EndMCBlockNumber,
	}
}
func (lim *lockinInfoModel) fromLockinInfo(l *LockinInfo) *lockinInfoModel {
	lim.SecretHash = l.SecretHash[:]
	lim.Secret = l.Secret[:]
	lim.MCUserAddressHex = l.MCUserAddressHex
	lim.SCUserAddress = l.SCUserAddress[:]
	lim.SCTokenAddress = l.SCTokenAddress[:]
	lim.Amount = l.Amount.Bytes()
	lim.MCExpiration = l.MCExpiration
	lim.SCExpiration = l.SCExpiration
	lim.MCLockStatus = int(l.MCLockStatus)
	lim.SCLockStatus = int(l.SCLockStatus)
	lim.Data = l.Data
	lim.NotaryIDInCharge = l.NotaryIDInCharge
	lim.StartTime = l.StartTime
	lim.StartMCBlockNumber = l.StartMCBlockNumber
	lim.EndTime = l.EndTime
	lim.EndMCBlockNumber = l.EndMCBlockNumber
	return lim
}

// NewLockinInfo :
func (db *DB) NewLockinInfo(lockinInfo *LockinInfo) (err error) {
	var t lockinInfoModel
	return db.Create(t.fromLockinInfo(lockinInfo)).Error
}

// GetAllLockinInfo :
func (db *DB) GetAllLockinInfo() (list []*LockinInfo, err error) {
	var t []lockinInfoModel
	err = db.Find(&t).Error
	if err == storm.ErrNotFound {
		err = nil
		return
	}
	for _, l := range t {
		list = append(list, l.toLockinInfo())
	}
	return
}

// GetLockinInfo :
func (db *DB) GetLockinInfo(secretHash common.Hash) (lockinInfo *LockinInfo, err error) {
	var lim lockinInfoModel
	err = db.Where(&lockinInfoModel{
		SecretHash: secretHash[:],
	}).First(&lim).Error
	if err != nil {
		return
	}
	lockinInfo = lim.toLockinInfo()
	return
}

// UpdateLockinInfo :
func (db *DB) UpdateLockinInfo(lockinInfo *LockinInfo) (err error) {
	var l *LockinInfo
	l, err = db.GetLockinInfo(lockinInfo.SecretHash)
	if l == nil {
		err = fmt.Errorf("can not update LockinInfo because key : secretHash=%s not found in db", lockinInfo.SecretHash.String())
		return
	}
	var t lockinInfoModel
	return db.Save(t.fromLockinInfo(lockinInfo)).Error
}
