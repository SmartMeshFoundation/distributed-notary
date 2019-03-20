package models

import (
	"fmt"

	"strings"

	"strconv"

	"sort"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/huamou/config"
	"github.com/jinzhu/gorm"
	"github.com/nkbai/log"
)

//NotaryInfo 公证人的基本信息
type NotaryInfo struct {
	Key        string `gorm:"primary_key" json:"-"` // 这里key就是ID的字符串格式
	ID         int    `json:"id"`                   // 公证人编号, 预先定死
	Name       string `json:"name"`
	Host       string `json:"host"` //how to contact with this notary
	AddressStr string `json:"address"`
}

type notaryInfoSlice []*NotaryInfo

func (c notaryInfoSlice) Len() int {
	return len(c)
}
func (c notaryInfoSlice) Less(i, j int) bool {
	return c[i].ID < c[j].ID
}
func (c notaryInfoSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

/*
SortNotaryInfoSlice :
要求有先后顺序,影响消息发送顺序
*/
func SortNotaryInfoSlice(chs []*NotaryInfo) {
	sort.Stable(notaryInfoSlice(chs))
}

// GetAddress :
func (ns *NotaryInfo) GetAddress() common.Address {
	return common.HexToAddress(ns.AddressStr)
}

// GetNotaryInfo :
func (db *DB) GetNotaryInfo() (notaries []*NotaryInfo, err error) {
	err = db.Order("id asc").Find(&notaries).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

// NewNotaryInfoFromConfFile :
func (db *DB) NewNotaryInfoFromConfFile(confFilePath string) (notaries []*NotaryInfo, err error) {
	if !utils.Exists(confFilePath) {
		err = fmt.Errorf("notary-conf-file does't exist")
		return
	}
	c, err := config.ReadDefault(confFilePath)
	if err != nil {
		err = fmt.Errorf("load notary-conf-file error: %s", err.Error())
		return
	}
	section := "NOTARY"
	options, err := c.Options(section)
	if err != nil {
		err = fmt.Errorf("load notary-conf-file error: %s", err.Error())
		return
	}
	for _, option := range options {
		s := strings.Split(c.RdString(section, option, ""), "-")
		id, err2 := strconv.Atoi(s[0])
		if err2 != nil {
			err = fmt.Errorf("load notary-conf-file error: %s", err2.Error())
			return
		}
		notaryInfo := NotaryInfo{
			Key:        s[0],
			ID:         id,
			Name:       "Notary-" + s[0],
			Host:       s[2],
			AddressStr: s[1],
		}
		err = db.Save(&notaryInfo).Error
		if err != nil {
			return
		}
		notaries = append(notaries, &notaryInfo)
	}
	log.Trace("read notary list from notary-conf-file %s :\n%s", confFilePath, utils.ToJSONStringFormat(notaries))
	return
}
