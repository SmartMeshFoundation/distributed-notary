package service

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

type lockoutHandler struct {
	db             *models.DB
	dealingMap     map[common.Hash]*models.LockoutInfo
	dealingMapLock sync.Mutex
}

func newLockoutHandler(db *models.DB, scTokenAddress common.Address) *lockoutHandler {
	lockoutInfoList, err := db.GetAllLockoutInfoBySCToken(scTokenAddress)
	if err != nil {
		panic(err)
	}
	h := &lockoutHandler{
		db:         db,
		dealingMap: make(map[common.Hash]*models.LockoutInfo),
	}
	for _, lockoutInfo := range lockoutInfoList {
		h.dealingMap[lockoutInfo.SecretHash] = lockoutInfo
	}
	return h
}

func (lh *lockoutHandler) registerLockout(lockoutInfo *models.LockoutInfo) (err error) {
	// 0. 状态校验
	if lockoutInfo == nil || lockoutInfo.MCLockStatus != models.LockStatusNone ||
		lockoutInfo.SCLockStatus != models.LockStatusLock || lockoutInfo.SecretHash == utils.EmptyHash {
		panic(fmt.Sprintf("call registerLockout with wrong lockoutInfo :\n%s", utils.ToJSONStringFormat(lockoutInfo)))
	}
	// 1. 写入db
	err = lh.db.NewLockoutInfo(lockoutInfo)
	if err != nil {
		err = fmt.Errorf("db.NewLockoutInfo err = %s", err.Error())
		return
	}
	// 2. 写入dealingMap
	lh.dealingMapLock.Lock()
	lh.dealingMap[lockoutInfo.SecretHash] = lockoutInfo
	lh.dealingMapLock.Unlock()
	return
}

func (lh *lockoutHandler) getLockout(secretHash common.Hash) (lockoutInfo *models.LockoutInfo, err error) {
	var ok bool
	// 1.优先查内存
	lh.dealingMapLock.Lock()
	if lockoutInfo, ok = lh.dealingMap[secretHash]; ok {
		lh.dealingMapLock.Unlock()
		return
	}
	lh.dealingMapLock.Unlock()
	// 2. 查db
	return lh.db.GetLockoutInfo(secretHash)
}

func (lh *lockoutHandler) updateLockout(lockoutInfo *models.LockoutInfo) (err error) {
	// 仅允许更新dealMap中的数据,否则就是实现有问题
	lh.dealingMapLock.Lock()
	if _, ok := lh.dealingMap[lockoutInfo.SecretHash]; !ok {
		panic("wrong code")
	}
	// 1. 写入db
	err = lh.db.UpdateLockoutInfo(lockoutInfo)
	if err != nil {
		lh.dealingMapLock.Lock()
		err = fmt.Errorf("db.UpdateLockoutInfo err = %s", err.Error())
		return
	}
	// 2. 如果lockoutInfo为完成状态,移除内存中数据,否则更新内存中数据
	if lockoutInfo.IsEnd() {
		delete(lh.dealingMap, lockoutInfo.SecretHash)
	} else {
		lh.dealingMap[lockoutInfo.SecretHash] = lockoutInfo
	}
	lh.dealingMapLock.Unlock()
	return
}

/*
处理lockout主链锁过期
主链锁定的是公证人的钱,如果过期,应尽快去cancel
*/
func (lh *lockoutHandler) onMCNewBlockEvent(mcNewBlockNumber uint64) (lockoutListNeedCancel []*models.LockoutInfo, err error) {
	lh.dealingMapLock.Lock()
	for _, lockout := range lh.dealingMap {
		if lockout.MCLockStatus == models.LockStatusLock && mcNewBlockNumber >= lockout.MCExpiration {
			//只处理锁住的
			lockout.MCLockStatus = models.LockStatusExpiration
			err = lh.db.UpdateLockoutInfo(lockout)
			if err != nil {
				lh.dealingMapLock.Unlock()
				err = fmt.Errorf("db.UpdateLockoutInfo err = %s", err.Error())
				return
			}
			// 记录需要cancel的部分,递交给上层处理
			lockoutListNeedCancel = append(lockoutListNeedCancel, lockout)
		}
	}
	lh.dealingMapLock.Unlock()
	return
}

/*
处理lockout侧链链锁过期
侧链锁定的是用户的钱,过期了之后公证人仅做记录,不去调用Cancel
*/
func (lh *lockoutHandler) onSCNewBlockEvent(scNewBlockNumber uint64) (err error) {
	lh.dealingMapLock.Lock()
	for _, lockout := range lh.dealingMap {
		if lockout.SCLockStatus == models.LockStatusLock && scNewBlockNumber >= lockout.SCExpiration {
			// 只处理锁住的
			lockout.SCLockStatus = models.LockStatusExpiration
			err = lh.db.UpdateLockoutInfo(lockout)
			if err != nil {
				lh.dealingMapLock.Unlock()
				err = fmt.Errorf("db.UpdateLockoutInfo err = %s", err.Error())
				return
			}
		}
	}
	lh.dealingMapLock.Unlock()
	return
}
