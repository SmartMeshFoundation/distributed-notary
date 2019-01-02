package models

import "github.com/ethereum/go-ethereum/common"

/*
SideChainTokenMetaInfo :
定义一个侧链token相关的所有元数据
*/
type SideChainTokenMetaInfo struct {
	MCName                   string         // 对应主链名
	SCToken                  common.Address // 侧链token地址
	SCTokenOwnerKey          common.Hash    // 侧链token合约owner的key
	MCLockedContractAddress  common.Address // 对应主链锁定合约地址
	MCLockedContractOwnerKey common.Hash    // 对应主链锁定合约owner的key
}

// GetSCTokenMetaInfoList :
func (db *DB) GetSCTokenMetaInfoList() (list []*SideChainTokenMetaInfo, err error) {
	return
}
