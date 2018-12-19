package utils

import "github.com/ethereum/go-ethereum/common"
import butils "github.com/nkbai/goutils"

//func Bytes2Hash(data []byte) common.Hash {
//	var h common.Hash
//	h.SetBytes(data)
//	return h
//}
//NewRandomHash 随机生成 hash
func NewRandomHash() common.Hash {
	return common.BytesToHash(butils.Random(32))
}
