package userapi

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	addr   = crypto.PubkeyToAddress(key.PublicKey)
)

func TestSCPrepareLockin2(t *testing.T) {
	req := &SCPrepareLockinRequest2{}
	req.Name = "User-SCPrepareLockin"
	req.RequestID = "0x03bc1c3ddeb1e428b61959b2eefe04e19cc679d1ec86bdb81a4c3d848bcdaa00"
	req.SCTokenAddress = common.HexToAddress("0x326dee230e67e5c124e9c36eae2126c2158bf361")
	req.SecretHash = common.HexToHash("0x1305497a65c3fdb9f3676548feea03e068662da2ad7f6e2cc23152082e365f4f")
	req.MCUserAddress = addr
	req.MCExpiration = big.NewInt(8644816)
	pubkey := crypto.CompressPubkey(&key.PublicKey) //[]byte类型
	req.Signer = hex.EncodeToString(pubkey)
	data, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}
	//{"name":"User-SCPrepareLockin","request_id":"0x03bc1c3ddeb1e428b61959b2eefe04e19cc679d1ec86bdb81a4c3d848bcdaa00","sc_token_address":"0x326dee230e67e5c124e9c36eae2126c2158bf361","signer":"03ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd3138","secret_hash":"0x1305497a65c3fdb9f3676548feea03e068662da2ad7f6e2cc23152082e365f4f","mc_user_address":"0x71562b71999873db5b286df957af199ec94617f7","mc_expiration":8644816}
	t.Logf("%s", data)
	digest := crypto.Keccak256([]byte(data))
	//digest=73aa7baa7ee1416b76aa69c3605c104e1995e4762a83eecfeea61862cb08d616
	t.Logf("digest=%s", hex.EncodeToString(digest))
	signaturee, err := crypto.Sign(digest, key)
	if err != nil {
		t.Error(err)
		return
	}
	req.Signature = hex.EncodeToString(signaturee)
	data, err = json.Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}
	//{"name":"User-SCPrepareLockin","request_id":"0x03bc1c3ddeb1e428b61959b2eefe04e19cc679d1ec86bdb81a4c3d848bcdaa00","sc_token_address":"0x326dee230e67e5c124e9c36eae2126c2158bf361","signer":"03ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd3138","signature":"16a7a9ccf3dfd347011cc4bcd2dffd48720712e7f837c0750973ebdb4a2ede1f3cced6fb280dc8103289db339120fdc3ce5f38f19bc741ec1e9a7b55865ecb5301","secret_hash":"0x1305497a65c3fdb9f3676548feea03e068662da2ad7f6e2cc23152082e365f4f","mc_user_address":"0x71562b71999873db5b286df957af199ec94617f7","mc_expiration":8644816}
	t.Logf("%s", data)
	if !req.VerifySign() {
		t.Error("signature verify error")
	}
}

func TestMCPrepareLockout2(t *testing.T) {
	req := &MCPrepareLockoutRequest2{}
	req.Name = "User-MCPrepareLockout"
	req.RequestID = "0x17c504eda9d55487783cda501a978a2ecc7beb0723906da494f9af795d06be4e"
	req.SCTokenAddress = common.HexToAddress("0x326dee230e67e5c124e9c36eae2126c2158bf361")
	req.SecretHash = common.HexToHash("0x1305497a65c3fdb9f3676548feea03e068662da2ad7f6e2cc23152082e365f4f")
	req.SCUserAddress = addr
	pubkey := crypto.CompressPubkey(&key.PublicKey) //[]byte类型
	req.Signer = hex.EncodeToString(pubkey)
	data, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}
	//{"name":"User-MCPrepareLockout","request_id":"0x17c504eda9d55487783cda501a978a2ecc7beb0723906da494f9af795d06be4e","sc_token_address":"0x326dee230e67e5c124e9c36eae2126c2158bf361","signer":"03ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd3138","secret_hash":"0x1305497a65c3fdb9f3676548feea03e068662da2ad7f6e2cc23152082e365f4f","sc_user_address":"0x71562b71999873db5b286df957af199ec94617f7"}
	t.Logf("%s", data)
	digest := crypto.Keccak256([]byte(data))
	//digest=f6b80c169ad021cbea0a8d225872bd56b8d41e15daec5173630734ea431edb48
	t.Logf("digest=%s", hex.EncodeToString(digest))
	signaturee, err := crypto.Sign(digest, key)
	if err != nil {
		t.Error(err)
		return
	}
	req.Signature = hex.EncodeToString(signaturee)
	data, err = json.Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}
	// {"name":"User-MCPrepareLockout","request_id":"0x17c504eda9d55487783cda501a978a2ecc7beb0723906da494f9af795d06be4e","sc_token_address":"0x326dee230e67e5c124e9c36eae2126c2158bf361","signer":"03ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd3138","signature":"446c8afe966be6960611fa8086e50d51f81eaeb2d88544ce542cbce62eb3616b2a40b98f88d7787eed6d6e6800251b42c51f744935266bb498dfeb70adaf31a601","secret_hash":"0x1305497a65c3fdb9f3676548feea03e068662da2ad7f6e2cc23152082e365f4f","sc_user_address":"0x71562b71999873db5b286df957af199ec94617f7"}
	t.Logf("%s", data)
	if !req.VerifySign() {
		t.Error("signature verify error")
	}
}
