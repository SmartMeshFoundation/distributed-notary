package models

import (
	"math/big"

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
	UserAddress        common.Address `json:"user_address"`        // 提出lockin请求用户的地址
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

// NewLockinInfo :
func (db *DB) NewLockinInfo(lockinInfo *LockinInfo) (err error) {
	// TODO
	return
}

// GetLockinInfo :
func (db *DB) GetLockinInfo(secretHash common.Hash) (lockinInfo *LockinInfo, err error) {
	// TODO
	return
}

// UpdateLockinInfo :
func (db *DB) UpdateLockinInfo(lockinInfo *LockinInfo) (err error) {
	// TODO
	return
}
