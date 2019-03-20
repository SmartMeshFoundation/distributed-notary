package models

import "github.com/jinzhu/gorm"

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
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return n.Nonce, nil
}
