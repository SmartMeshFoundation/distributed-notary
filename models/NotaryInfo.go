package models

import (
	"fmt"

	"strings"

	"strconv"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/huamou/config"
	"github.com/jinzhu/gorm"
)

//NotaryInfo 公证人的基本信息
type NotaryInfo struct {
	ID         int `gorm:"primary_key"` // 公证人编号, 预先定死
	Name       string
	Host       string //how to contact with this notary
	AddressStr string
}

// GetAddress :
func (ns *NotaryInfo) GetAddress() common.Address {
	return common.HexToAddress(ns.AddressStr)
}

// GetNotaryInfo :
func (db *DB) GetNotaryInfo() (notaries []NotaryInfo, err error) {
	err = db.Find(&notaries).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

// NewNotaryInfoFromConfFile :
func (db *DB) NewNotaryInfoFromConfFile(confFilePath string) (notaries []NotaryInfo, err error) {
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
			ID:         id,
			Name:       "Notary-" + s[0],
			Host:       s[2],
			AddressStr: s[1],
		}
		err = db.Save(&notaryInfo).Error
		if err != nil {
			return
		}
		notaries = append(notaries, notaryInfo)
	}
	return
}
