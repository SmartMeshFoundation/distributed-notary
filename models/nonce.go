package models

import (
	"fmt"

	"github.com/kataras/go-errors"

	"github.com/jinzhu/gorm"
)

var (
	//ErrDuplicateNnonce 在以太坊nonce分配中,为同一笔交易分配了不同的nonce
	ErrDuplicateNnonce = errors.New("duplicate nonce")
)

//Nonce for pbft 协商历史
type Nonce struct {
	Key   string `gorm:"primary_key"` //哪一条链上的哪个账户
	Nonce int    //nonce 编号
}

//UpdateNonce 更新最新协商的nonce
func (db *DB) UpdateNonce(key string, nonce int) (err error) {
	return db.Save(&Nonce{key, nonce}).Error
}

//GetNonce 获取指定账户的nonce
func (db *DB) GetNonce(key string) (nonce int, err error) {
	n := &Nonce{
		Key: key,
	}
	err = db.Where(n).First(n).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return -1, nil
		}
		return
	}
	return n.Nonce, nil
}

//记录 为每一个op分配的nonce,主要用于以太坊
type OpNonce struct {
	Key     string `gorm:"primary_key"`
	Nonce   int    //为op分配的Nonce
	Account string //在哪一个私钥上协商
	Op      string
}

//getOpNonce 获取指定账户的nonce
func (db *DB) getOpNonce(key string) (opNOnce *OpNonce, err error) {
	n := &OpNonce{
		Key: key,
	}
	err = db.Where(n).First(n).Error
	if err != nil {
		return
	}
	return n, nil
}

//NewNonceForOp 主节点为op分配了一个Nonce,确保没有数据库确保没有重复分配
//不同的chain会采用相同的op
func (db *DB) NewNonceForOp(view, nonce int, op, chain, account string) error {
	key := fmt.Sprintf("%s%s%d", op, chain, view)
	on, err := db.getOpNonce(key)
	//不存在,新创建就行了
	if err != nil {
		return db.Save(&OpNonce{
			Key:     key,
			Nonce:   nonce,
			Account: account,
			Op:      op,
		}).Error
	}
	//已存在记录,并且分配了不同的nonce
	if on.Nonce != nonce {
		return ErrDuplicateNnonce
	}
	return nil
}
