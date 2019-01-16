package service

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

type lockinHandler struct {
	db             *models.DB
	dealingMap     map[common.Hash]*models.LockinInfo
	dealingMapLock sync.Mutex
}

func newLockinHandler(db *models.DB, scTokenAddress common.Address) *lockinHandler {
	lockinInfoList, err := db.GetAllLockinInfoBySCToken(scTokenAddress)
	if err != nil {
		panic(err)
	}
	h := &lockinHandler{
		db:         db,
		dealingMap: make(map[common.Hash]*models.LockinInfo),
	}
	for _, lockinInfo := range lockinInfoList {
		h.dealingMap[lockinInfo.SecretHash] = lockinInfo
	}
	return h
}

func (lh *lockinHandler) registerLockin(lockinInfo *models.LockinInfo) (err error) {
	// 0. 状态校验
	if lockinInfo == nil || lockinInfo.MCLockStatus != models.LockStatusLock ||
		lockinInfo.SCLockStatus != models.LockStatusNone || lockinInfo.SecretHash == utils.EmptyHash {
		panic(fmt.Sprintf("call registerLockin with wrong lockinInfo :\n%s", utils.ToJSONStringFormat(lockinInfo)))
	}
	// 1. 写入db
	err = lh.db.NewLockinInfo(lockinInfo)
	if err != nil {
		err = fmt.Errorf("db.NewLockinInfo err = %s", err.Error())
		return
	}
	// 2. 写入dealingMap
	lh.dealingMapLock.Lock()
	lh.dealingMap[lockinInfo.SecretHash] = lockinInfo
	lh.dealingMapLock.Unlock()
	return
}

func (lh *lockinHandler) getLockin(secretHash common.Hash) (lockinInfo *models.LockinInfo, err error) {
	var ok bool
	// 1.优先查内存
	lh.dealingMapLock.Lock()
	if lockinInfo, ok = lh.dealingMap[secretHash]; ok {
		lh.dealingMapLock.Unlock()
		return
	}
	lh.dealingMapLock.Unlock()
	// 2. 查db
	return lh.db.GetLockinInfo(secretHash)
}

func (lh *lockinHandler) updateLockin(lockinInfo *models.LockinInfo) (err error) {
	// 仅允许更新dealMap中的数据,否则就是实现有问题
	lh.dealingMapLock.Lock()
	if _, ok := lh.dealingMap[lockinInfo.SecretHash]; !ok {
		panic("wrong code")
	}
	// 1. 写入db
	err = lh.db.UpdateLockinInfo(lockinInfo)
	if err != nil {
		lh.dealingMapLock.Lock()
		err = fmt.Errorf("db.UpdateLockinInfo err = %s", err.Error())
		return
	}
	// 2. 如果lockinInfo为完成状态,移除内存中数据,否则更新内存中数据
	if lockinInfo.IsEnd() {
		delete(lh.dealingMap, lockinInfo.SecretHash)
	} else {
		lh.dealingMap[lockinInfo.SecretHash] = lockinInfo
	}
	lh.dealingMapLock.Unlock()
	return
}

/*
处理lockin主链锁过期
主链锁定的是用户的钱,过期了之后公证人仅做记录,不去调用Cancel
*/
func (lh *lockinHandler) onMCNewBlockEvent(mcNewBlockNumber uint64) (err error) {
	lh.dealingMapLock.Lock()
	for _, lockin := range lh.dealingMap {
		if lockin.MCLockStatus == models.LockStatusLock && mcNewBlockNumber >= lockin.MCExpiration {
			//只处理锁住的
			lockin.MCLockStatus = models.LockStatusExpiration
			err = lh.db.UpdateLockinInfo(lockin)
			if err != nil {
				lh.dealingMapLock.Unlock()
				err = fmt.Errorf("db.UpdateLockinInfo err = %s", err.Error())
				return
			}
		}
	}
	lh.dealingMapLock.Unlock()
	return
}

/*
处理lockin侧链链锁过期
侧链锁定的是公证人的钱,如果过期,应尽快去cancel
*/
func (lh *lockinHandler) onSCNewBlockEvent(scNewBlockNumber uint64) (lockinListNeedCancel []*models.LockinInfo, err error) {
	lh.dealingMapLock.Lock()
	for _, lockin := range lh.dealingMap {
		if lockin.SCLockStatus == models.LockStatusLock && scNewBlockNumber >= lockin.SCExpiration {
			// 只处理锁住的
			lockin.SCLockStatus = models.LockStatusExpiration
			err = lh.db.UpdateLockinInfo(lockin)
			if err != nil {
				lh.dealingMapLock.Unlock()
				err = fmt.Errorf("db.UpdateLockinInfo err = %s", err.Error())
				return
			}
			// 记录需要cancel的部分,递交给上层处理
			lockinListNeedCancel = append(lockinListNeedCancel, lockin)
		}
	}
	lh.dealingMapLock.Unlock()
	return
}
