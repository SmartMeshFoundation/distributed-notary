package utils

import (
	"io"

	rand2 "crypto/rand"

	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//EmptyHash all zero,invalid
var EmptyHash = common.Hash{}

//EmptyAddress all zero,invalid
var EmptyAddress = common.Address{}

//func Bytes2Hash(data []byte) common.Hash {
//	var h common.Hash
//	h.SetBytes(data)
//	return h
//}

//NewRandomHash generate random hash,for testonly
func NewRandomHash() common.Hash {
	return Sha3(Random(64))
}

func readFullOrPanic(r io.Reader, v []byte) int {
	n, err := io.ReadFull(r, v)
	if err != nil {
		panic(err)
	}
	return n
}

// Random takes a parameter (int) and returns random slice of byte
// ex: var randomstrbytes []byte; randomstrbytes = utils.Random(32)
func Random(n int) []byte {
	v := make([]byte, n)
	readFullOrPanic(rand2.Reader, v)
	return v
}

//Sha3 is short for Keccak256Hash
func Sha3(data ...[]byte) common.Hash {
	return crypto.Keccak256Hash(data...)
}

func ToJsonStringFormat(v interface{}) string {
	buf, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
