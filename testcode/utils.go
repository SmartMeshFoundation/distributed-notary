package testcode

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

// GetTestPrivateKey1 :
func GetTestPrivateKey1() *ecdsa.PrivateKey {
	key, err := hex.DecodeString("4359f525e2b373089be5fe8f9a4e8ffb6d30e2960918be426217921e1b2547f7")
	if err != nil {
		panic(err)
	}
	privatekey, err := crypto.ToECDSA(key)
	if err != nil {
		panic(err)
	}
	return privatekey
}
