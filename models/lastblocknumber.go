package models

// LastBlockNumberInfo :
type LastBlockNumberInfo struct {
	ChainName       string `gorm:"primary_key"`
	LastBlockNumber uint64
}

// GetLastBlockNumber :
func (db *DB) GetLastBlockNumber(chainName string) uint64 {
	var lb LastBlockNumberInfo
	err := db.Where(&LastBlockNumberInfo{
		ChainName: chainName,
	}).First(&lb).Error
	if err != nil {
		return 0
	}
	return lb.LastBlockNumber
}

// SaveLastBlockNumber :
func (db *DB) SaveLastBlockNumber(chainName string, blockNumber uint64) (err error) {
	return db.Save(&LastBlockNumberInfo{
		ChainName:       chainName,
		LastBlockNumber: blockNumber,
	}).Error
}
