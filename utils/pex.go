package utils

import "github.com/ethereum/go-ethereum/common"

//Pex short string stands for data
func Pex(data []byte) string {
	return common.Bytes2Hex(data[:4])
}

//HPex pex for hash
func HPex(data common.Hash) string {
	return common.Bytes2Hex(data[:2])
}

//BPex bytes to string
func BPex(data []byte) string {
	return common.Bytes2Hex(data)
}

//APex pex for address
func APex(data common.Address) string {
	return common.Bytes2Hex(data[:4])
}

//APex2 shorter than APex
func APex2(data common.Address) string {
	return common.Bytes2Hex(data[:2])
}
