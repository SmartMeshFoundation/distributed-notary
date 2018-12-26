package models

import "crypto/ecdsa"

//NotaryInfo 公证人的基本信息
type NotaryInfo struct {
	ID        int // 公证人编号, 预先定死
	Name      string
	Host      string //how to contact with this notary
	PublicKey *ecdsa.PublicKey
}

// GetNotaryInfoMap :
func (db *DB) GetNotaryInfoMap() map[int]*NotaryInfo {
	// TODO
	return nil
}
