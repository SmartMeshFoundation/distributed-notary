package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
)

/*
SideChainTokenMetaInfo :
定义一个侧链token相关的所有元数据
*/
type SideChainTokenMetaInfo struct {
	SCToken                  common.Address `json:"sc_token"`                               // 侧链token地址
	SCTokenName              string         `json:"sc_token_name"`                          // 侧链Token名
	SCTokenOwnerKey          common.Hash    `json:"sc_token_owner_key"`                     // 侧链token合约owner的key
	MCLockedContractAddress  common.Address `json:"mc_locked_contract_address"`             // 对应主链锁定合约地址
	MCName                   string         `json:"mc_name"`                                // 对应主链名
	MCLockedContractOwnerKey common.Hash    `json:"mc_locked_contract_owner_key,omitempty"` // 对应主链锁定合约owner的key
}

type sideChainTokenMetaInfoModel struct {
	SCToken                  []byte `gorm:"primary_key"`
	SCTokenName              string // 侧链Token名
	SCTokenOwnerKey          []byte // 侧链token合约owner的key
	MCLockedContractAddress  []byte // 对应主链锁定合约地址
	MCName                   string // 对应主链名
	MCLockedContractOwnerKey []byte // 对应主链锁定合约owner的key
}

func (m *sideChainTokenMetaInfoModel) toSideChainTokenMetaInfo() *SideChainTokenMetaInfo {
	return &SideChainTokenMetaInfo{
		SCToken:                  common.BytesToAddress(m.SCToken),
		SCTokenName:              m.SCTokenName,
		SCTokenOwnerKey:          common.BytesToHash(m.SCTokenOwnerKey),
		MCLockedContractAddress:  common.BytesToAddress(m.MCLockedContractAddress),
		MCName:                   m.MCName,
		MCLockedContractOwnerKey: common.BytesToHash(m.MCLockedContractOwnerKey),
	}
}

func (m *sideChainTokenMetaInfoModel) fromSideChainTokenMetaInfo(sc *SideChainTokenMetaInfo) *sideChainTokenMetaInfoModel {
	m.SCToken = sc.SCToken[:]
	m.SCTokenName = sc.SCTokenName
	m.SCTokenOwnerKey = sc.SCTokenOwnerKey[:]
	m.MCLockedContractAddress = sc.MCLockedContractAddress[:]
	m.MCName = sc.MCName
	m.MCLockedContractOwnerKey = sc.MCLockedContractOwnerKey[:]
	return m
}

// GetSCTokenMetaInfoList :
func (db *DB) GetSCTokenMetaInfoList() (list []*SideChainTokenMetaInfo, err error) {
	var l []sideChainTokenMetaInfoModel
	err = db.Find(&l).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	for _, m := range l {
		list = append(list, m.toSideChainTokenMetaInfo())
	}
	return
}

// NewSCTokenMetaInfo :
func (db *DB) NewSCTokenMetaInfo(scToken *SideChainTokenMetaInfo) (err error) {
	var m sideChainTokenMetaInfoModel
	return db.Create(m.fromSideChainTokenMetaInfo(scToken)).Error
}
