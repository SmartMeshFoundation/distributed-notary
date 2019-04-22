package models

import (
	"github.com/asdine/storm"
	"github.com/btcsuite/btcutil"
	"github.com/jinzhu/gorm"
)

// BTCOutpointStatus :
type BTCOutpointStatus int

const (
	// BTCOutpointStatusUsable 初始状态,即可用状态
	BTCOutpointStatusUsable = iota
	// BTCOutpointStatusUsed 已使用状态
	BTCOutpointStatusUsed
)

// BTCOutpoint 保存公证人分布式私钥对应地址上可用的普通utxo
type BTCOutpoint struct {
	PublicKeyHashStr string            `json:"public_key_hash_str"`
	TxHashStr        string            `json:"tx_hash" gorm:"primary_key"` // utxo所在的txHash
	Index            int               `json:"index"`                      // utxo在tx中的index
	Amount           btcutil.Amount    `json:"amount"`                     // 金额
	Status           BTCOutpointStatus `json:"status"`                     // 0-可用 1-已使用
	CreateTime       int64             `json:"create_time"`                // 创建时间
	UseTime          int64             `json:"use_time"`                   // 使用的时间
}

//NewBTCOutpoint :
func (db *DB) NewBTCOutpoint(outpoint *BTCOutpoint) error {
	return db.Create(outpoint).Error
}

// GetBTCOutpointList 条件查询
func (db *DB) GetBTCOutpointList(status BTCOutpointStatus) (list []*BTCOutpoint) {
	if status == -1 {
		err := db.Find(&list).Error
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	err := db.Where(&BTCOutpoint{
		Status: status,
	}).Find(&list).Error
	if err == storm.ErrNotFound {
		err = nil
	}
	return
}
