package utils

import (
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

// NewRandomBTCPKH :
func NewRandomBTCPKH(net *chaincfg.Params) *btcutil.AddressPubKeyHash {
	t, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		panic(err)
	}
	pubKeyHash := btcutil.Hash160(t.PubKey().SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, net)
	if err != nil {
		panic(err)
	}
	return addr
}

// PrintBTCBalanceOfAccount :
func PrintBTCBalanceOfAccount(walletConn *rpcclient.Client, account string) btcutil.Amount {
	if walletConn == nil {
		return 0
	}
	balance, err := walletConn.GetBalance(account)
	if err != nil {
		fmt.Println("GetBalance err : ", err)
		os.Exit(-1)
	}
	fmt.Printf("account %s balance : %s\n", account, balance.String())
	return balance
}
