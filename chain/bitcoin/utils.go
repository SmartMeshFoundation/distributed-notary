package bitcoin

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// PrivateKeyBytes2AddressPublicKeyHash :
func PrivateKeyBytes2AddressPublicKeyHash(k []byte, net *chaincfg.Params) *btcutil.AddressPubKeyHash {
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), k)
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, net)
	if err != nil {
		panic(err)
	}
	return addr
}

// PrivateKeyBytes2PrivateKey :
func PrivateKeyBytes2PrivateKey(k []byte) *btcec.PrivateKey {
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), k)
	return privateKey
}
