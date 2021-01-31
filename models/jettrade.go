package models

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	utils "github.com/nkbai/goutils"
	"github.com/nkbai/log"
)

type JettradeEventInfo struct {
	ID               int
	ChainName        string
	FromAddress      common.Address
	BlockNumber      uint64
	EventName        string
	From             common.Address
	To               common.Address
	TokenID          *big.Int
	NotaryIDInCharge int
	TxHash           common.Hash
}

func NewJettradeEventInfo(chainName, eventName string, fromContractAddress, from, to common.Address, blockNumber uint64, tokenID *big.Int, txHash common.Hash) *JettradeEventInfo {
	return &JettradeEventInfo{
		ChainName:        chainName,
		FromAddress:      fromContractAddress,
		BlockNumber:      blockNumber,
		EventName:        eventName,
		From:             from,
		To:               to,
		TokenID:          tokenID,
		NotaryIDInCharge: -1,
		TxHash:           txHash,
	}
}

type jettradeEventInfoModel struct {
	ID               int `gorm:"primary_key"`
	ChainName        string
	FromAddress      string
	BlockNumber      uint64
	EventName        string
	From             string
	To               string
	TokenID          []byte
	NotaryIDInCharge int32
	TxHash           string
}

func (m *jettradeEventInfoModel) toJettradeEventInfo() *JettradeEventInfo {
	tokenID := new(big.Int)
	tokenID.SetBytes(m.TokenID)
	return &JettradeEventInfo{
		ID:               m.ID,
		ChainName:        m.ChainName,
		FromAddress:      common.HexToAddress(m.FromAddress),
		BlockNumber:      m.BlockNumber,
		EventName:        m.EventName,
		From:             common.HexToAddress(m.From),
		To:               common.HexToAddress(m.To),
		TokenID:          tokenID,
		NotaryIDInCharge: int(m.NotaryIDInCharge),
		TxHash:           common.HexToHash(m.TxHash),
	}
}
func (m *jettradeEventInfoModel) fromJettradeEventInfo(j *JettradeEventInfo) *jettradeEventInfoModel {
	m.ID = j.ID
	m.ChainName = j.ChainName
	m.FromAddress = j.FromAddress.String()
	m.BlockNumber = j.BlockNumber
	m.EventName = j.EventName
	m.From = j.From.String()
	m.To = j.To.String()
	m.TokenID = j.TokenID.Bytes()
	m.TxHash = j.TxHash.String()
	m.NotaryIDInCharge = int32(j.NotaryIDInCharge)
	return m
}
func (db *DB) NewJettradeEventInfo(j *JettradeEventInfo) (err error) {
	var t jettradeEventInfoModel
	t.fromJettradeEventInfo(j)
	log.Info("NewJettradeEventInfo %s", utils.StringInterface(t, 3))
	return db.Create(&t).Error
}

// GetAllJettradeEventInfo :
func (db *DB) GetAllJettradeEventInfo() (list []*JettradeEventInfo, err error) {
	var t []jettradeEventInfoModel
	err = db.Find(&t).Error
	if err != nil {
		return
	}
	for _, l := range t {
		list = append(list, l.toJettradeEventInfo())
	}
	return
}

// GetJettradeEventInfo : 确保只会出现一次
func (db *DB) GetJettradeEventInfo(chainName, eventName string, fromContract common.Address, tokeID *big.Int) (j *JettradeEventInfo, err error) {
	var lim jettradeEventInfoModel
	err = db.Where(&jettradeEventInfoModel{
		ChainName:   chainName,
		EventName:   eventName,
		TokenID:     tokeID.Bytes(),
		FromAddress: fromContract.String(),
	}).First(&lim).Error
	if err != nil {
		return
	}
	j = lim.toJettradeEventInfo()
	return
}

func (db *DB) UpdateJettradeEventInfo(j *JettradeEventInfo) (err error) {
	var j2 *JettradeEventInfo
	j2, err = db.GetJettradeEventInfo(j.ChainName, j.EventName, j.FromAddress, j.TokenID)
	if j2 == nil {
		err = fmt.Errorf("can not update jettrade eventinfo chain=%s,event=%s, id=%s", j.ChainName, j.EventName, j.TokenID)
		return
	}
	var t jettradeEventInfoModel
	return db.Save(t.fromJettradeEventInfo(j)).Error
}
