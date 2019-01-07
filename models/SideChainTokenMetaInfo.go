package models

import "github.com/ethereum/go-ethereum/common"

/*
SideChainTokenMetaInfo :
定义一个侧链token相关的所有元数据
*/
type SideChainTokenMetaInfo struct {
	SCTokenName              string         `json:"sc_token_name"`                // 侧链Token名
	MCName                   string         `json:"mc_name"`                      // 对应主链名
	SCToken                  common.Address `json:"sc_token"`                     // 侧链token地址
	SCTokenOwnerKey          common.Hash    `json:"sc_token_owner_key"`           // 侧链token合约owner的key
	MCLockedContractAddress  common.Address `json:"mc_locked_contract_address"`   // 对应主链锁定合约地址
	MCLockedContractOwnerKey common.Hash    `json:"mc_locked_contract_owner_key"` // 对应主链锁定合约owner的key
}

// GetSCTokenMetaInfoList :
func (db *DB) GetSCTokenMetaInfoList() (list []*SideChainTokenMetaInfo, err error) {
	// TODO
	return
}

// NewSCTokenMetaInfo :
func (db *DB) NewSCTokenMetaInfo(scToken *SideChainTokenMetaInfo) (err error) {
	// TODO
	return
}
