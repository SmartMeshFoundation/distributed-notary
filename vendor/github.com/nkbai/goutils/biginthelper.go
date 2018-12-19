package utils

import (
	"encoding/hex"
	"io"

	"math/big"

	"crypto/sha256"
	"fmt"

	"github.com/nkbai/log"
)

//Pex short string stands for data
func Pex(data []byte) string {
	return hex.EncodeToString(data[:4])
}

//ShaSecret is short for sha256
func Sha256(data ...[]byte) []byte {
	//	return crypto.Keccak256Hash(data...)
	var d = sha256.New()
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

//BPex bytes to string
func BPex(data []byte) string {
	return hex.EncodeToString(data)
}

//BigIntTo32Bytes convert a big int to bytes
func BigIntTo32Bytes(i *big.Int) []byte {
	data := i.Bytes()
	buf := make([]byte, 32)
	for i := 0; i < 32-len(data); i++ {
		buf[i] = 0
	}
	for i := 32 - len(data); i < 32; i++ {
		buf[i] = data[i-32+len(data)]
	}
	return buf
}

//ReadBigInt read big.Int from buffer
func ReadBigInt(reader io.Reader) *big.Int {
	bi := new(big.Int)
	tmpbuf := make([]byte, 32)
	_, err := reader.Read(tmpbuf)
	if err != nil {
		log.Error(fmt.Sprintf("read BigInt error %s", err))
	}
	bi.SetBytes(tmpbuf)
	return bi
}
