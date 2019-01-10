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

/*
LockinInfo :
该结构体记录一次Lockin的所有数据及状态
*/
type LockinInfo struct {
	UserAddress    common.Address `json:"user_address"`     // 提出lockin请求用户的地址,唯一ID
	SCTokenAddress common.Address `json:"sc_token_address"` // SCToken
	PrivateKeyID   common.Hash    `json:"private_key_id"`   // 该次Lockin使用的keyID
	Amount         *big.Int       `json:"amount"`           // 金额
	MCExpiration   int64          `json:"mc_expiration"`    // 主链过期块
	SCExpiration   int64          `json:"sc_expiration"`    // 侧链过期块
	MCLockStatus   LockStatus     `json:"mc_lock_status"`   // 主链锁状态
	SCLockStatus   LockStatus     `json:"sc_lock_status"`   // 侧链锁状态
}
