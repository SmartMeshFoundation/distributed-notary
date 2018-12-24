package models

import (
	"github.com/ethereum/go-ethereum/common"
)

/*
NotaryInfo :
定义一个侧链token对应的公证人中的一位
*/
type NotaryInfo struct {
	Key            string         `gorm:"primary_key"` // 唯一ID,格式为 SCTokenAddress-Index
	Index          int            // 公证人序号,从1开始递增
	Address        common.Address // 公证人地址
	APIHost        string         // 公证人restful接口地址,ip:port
	SCTokenAddress common.Address // 该组公证人对应的侧链token地址
}
