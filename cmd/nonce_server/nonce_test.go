package main

import (
	"testing"

	"fmt"

	"context"

	"crypto/ecdsa"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

func eth2Wei(ethAmount int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(params.Ether)), big.NewInt(ethAmount))
}

func TestTT(t *testing.T) {
	c, _ := ethclient.Dial("http://127.0.0.1:17888")
	account := common.HexToAddress("0x53D9591A9033c72caDa8eA3D8D51b858B850CCD5")
	nonce, _ := c.NonceAt(context.Background(), account, nil)
	fmt.Println(nonce)
}

func TestNonce(t *testing.T) {
	c, _ := ethclient.Dial("http://127.0.0.1:17888")
	key, _ := crypto.HexToECDSA("36234555bc087435cf52371f9a0139cb98a4267ba62b722e3f46b90d35f31678")
	account := crypto.PubkeyToAddress(key.PublicKey)
	fmt.Println("test account=", account.String())

	nonce, _ := c.NonceAt(context.Background(), account, nil)
	fmt.Println("nonce :", nonce)
	pendingNonce, _ := c.PendingNonceAt(context.Background(), account)
	fmt.Println("pendingNonce :", pendingNonce)
	b, _ := c.BalanceAt(context.Background(), account, nil)
	fmt.Println("=====", b)
	fmt.Println("transfer1")
	err := transfer(key, account, big.NewInt(1), c)
	if err != nil {
		fmt.Println(err)
	}

	nonce, _ = c.NonceAt(context.Background(), account, nil)
	fmt.Println("nonce :", nonce)
	pendingNonce, _ = c.PendingNonceAt(context.Background(), account)
	fmt.Println("pendingNonce :", pendingNonce)

	fmt.Println("transfer2")
	err = transferWithNonce(key, utils.NewRandomAddress(), big.NewInt(1), c, 501)
	if err != nil {
		fmt.Println(err)
	}

	nonce, _ = c.NonceAt(context.Background(), account, nil)
	fmt.Println("nonce :", nonce)
	pendingNonce, _ = c.PendingNonceAt(context.Background(), account)
	fmt.Println("pendingNonce :", pendingNonce)
}

// Transfer10ToAccount : impl chain.Chain
func transfer(key *ecdsa.PrivateKey, accountTo common.Address, amount *big.Int, c *ethclient.Client) (err error) {
	if amount == nil || amount.Cmp(big.NewInt(0)) == 0 {
		return
	}
	conn := c
	ctx := context.Background()
	auth := bind.NewKeyedTransactor(key)
	fromAddr := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := conn.NonceAt(ctx, fromAddr, nil)
	if err != nil {
		return err
	}
	msg := ethereum.CallMsg{From: fromAddr, To: &accountTo, Value: amount, Data: nil}
	gasLimit, err := conn.EstimateGas(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to estimate gas needed: %v", err)
	}
	gasPrice, err := conn.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %v", err)
	}
	chainID, err := conn.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get networkID : %v", err)
	}
	rawTx := types.NewTransaction(nonce, accountTo, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(types.NewEIP155Signer(chainID), auth.From, rawTx)
	if err != nil {
		return err
	}
	if err = conn.SendTransaction(ctx, signedTx); err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, conn, signedTx)
	return
}

// Transfer10ToAccount : impl chain.Chain
func transferWithNonce(key *ecdsa.PrivateKey, accountTo common.Address, amount *big.Int, c *ethclient.Client, nonce uint64) (err error) {
	if amount == nil || amount.Cmp(big.NewInt(0)) == 0 {
		return
	}
	conn := c
	ctx := context.Background()
	auth := bind.NewKeyedTransactor(key)
	fromAddr := crypto.PubkeyToAddress(key.PublicKey)
	msg := ethereum.CallMsg{From: fromAddr, To: &accountTo, Value: amount, Data: nil}
	gasLimit, err := conn.EstimateGas(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to estimate gas needed: %v", err)
	}
	gasPrice, err := conn.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %v", err)
	}
	chainID, err := conn.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get networkID : %v", err)
	}
	rawTx := types.NewTransaction(nonce, accountTo, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(types.NewEIP155Signer(chainID), auth.From, rawTx)
	if err != nil {
		return err
	}
	if err = conn.SendTransaction(ctx, signedTx); err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, conn, signedTx)
	return
}
