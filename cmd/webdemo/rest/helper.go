package rest

import (
	"context"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/utils"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/ethereum/go-ethereum/common"

	mcontracts "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
	scontracts "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ant0ine/go-json-rest/rest"
)

type Response struct {
	Error   string
	Message interface{}
}

func errString(w rest.ResponseWriter, err string) {
	e := w.WriteJson(Response{
		Error: err,
	})
	if e != nil {
		log.Error(fmt.Sprintf("write json err %s", e))
	}
	return
}
func errError(w rest.ResponseWriter, err error) {
	e := w.WriteJson(Response{
		Error: err.Error(),
	})
	if e != nil {
		log.Error(fmt.Sprintf("write json err %s", e))
	}
	return
}
func success(w rest.ResponseWriter, r interface{}) {
	e := w.WriteJson(Response{
		Error:   "",
		Message: r,
	})
	if e != nil {
		log.Error(fmt.Sprintf("write json err %s", e))
	}
	return
}
func pubkey2Address(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> pubkey2Address ,err=%v", err))
	}()
	pubkey := common.Hex2Bytes(r.PathParam("pubkey"))
	if len(pubkey) != 32*2+1 {
		errString(w, "arg error")
		return
	}
	if pubkey[0] != 0x04 { //未压缩Key
		errString(w, "arg error")
		return
	}
	addr := bytes2Address(pubkey)
	log.Trace(fmt.Sprintf("addr=%s", addr.String()))
	success(w, addr)
}

type secretStruct struct {
	Secret     common.Hash
	SecretHash common.Hash
}

func generateSecret(w rest.ResponseWriter, r *rest.Request) {
	ss := secretStruct{}
	ss.Secret = utils.NewRandomHash()
	ss.SecretHash = utils.ShaSecret(ss.Secret[:])
	success(w, ss)
}
func bytes2Address(pubkey []byte) common.Address {
	addr := common.BytesToAddress(crypto.Keccak256(pubkey[1:])[12:])
	return addr
}

type statusReq struct {
	MainChainContract common.Address
	SideChainContract common.Address
	Account           common.Address
	LockSecretHash    common.Hash
}
type lockStruct struct {
	SecretHash common.Hash
	Expiration int64
	Value      *big.Int
}
type statusResponse struct {
	MainChainBlockNumber     int64
	SideChainBlockNumber     int64
	MainChainContractBalance *big.Int //主链合约锁定多少Eth
	SideChainContractBalance *big.Int //侧链token总供应量
	MainChainBalance         *big.Int //账户主链Eth多少
	SideChainBalance         *big.Int //账户侧链Smt多少
	SideChainTokenBalance    *big.Int //侧链账户有多少EthToken
	MainChainLockIn          *lockStruct
	MainChainLockout         *lockStruct
	SideChainLockin          *lockStruct
	SideChainLockout         *lockStruct
}

func queryStatus(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var value *big.Int
	sr := statusResponse{}
	defer func() {
		if err != nil {
			errError(w, err)
		} else {
			success(w, sr)
		}
		log.Trace(fmt.Sprintf("Restful Api Call ----> pubkey2Address ,err=%v", err))
	}()
	var req statusReq
	err = r.DecodeJsonPayload(&req)
	if err != nil {
		return
	}
	mclient, err := ethclient.Dial(MainChainEndpoint)
	if err != nil {
		return
	}
	sclient, err := ethclient.Dial(SideChainEndpoint)
	if err != nil {
		return
	}

	ctx := context.Background()
	h, err := mclient.HeaderByNumber(ctx, nil)
	if err != nil {
		return
	}
	sr.MainChainBlockNumber = h.Number.Int64()
	h, err = sclient.HeaderByNumber(ctx, nil)
	if err != nil {
		return
	}
	sr.SideChainBlockNumber = h.Number.Int64()
	if req.MainChainContract != utils.EmptyAddress {
		value, err = mclient.BalanceAt(ctx, req.MainChainContract, nil)
		if err != nil {
			return
		}
		sr.MainChainContractBalance = new(big.Int).Set(value)
	}
	if req.SideChainContract == utils.EmptyAddress {
		return
	}
	mc, err := mcontracts.NewLockedEthereum(req.MainChainContract, mclient)
	if err != nil {
		return
	}
	sc, err := scontracts.NewAtmosphereToken(req.SideChainContract, sclient)
	if err != nil {
		return
	}
	value, err = sc.TotalSupply(nil)
	if err != nil {
		return
	}
	sr.SideChainContractBalance = new(big.Int).Set(value)

	if req.Account == utils.EmptyAddress {
		return

	}
	value, err = mclient.BalanceAt(ctx, req.Account, nil)
	if err != nil {
		return
	}
	sr.MainChainBalance = new(big.Int).Set(value)
	value, err = sclient.BalanceAt(ctx, req.Account, nil)
	if err != nil {
		return
	}
	sr.SideChainBalance = new(big.Int).Set(value)

	value, err = sc.BalanceOf(nil, req.Account)
	if err != nil {
		return
	}
	sr.SideChainTokenBalance = new(big.Int).Set(value)

	l, e, a, err := sc.QueryLockin(nil, req.Account)
	if err != nil {
		return
	}
	sr.SideChainLockin = &lockStruct{
		SecretHash: l,
		Expiration: e.Int64(),
		Value:      a,
	}
	l, e, a, err = sc.QueryLockout(nil, req.Account)
	if err != nil {
		return
	}
	sr.SideChainLockout = &lockStruct{
		SecretHash: l,
		Expiration: e.Int64(),
		Value:      a,
	}
	l, e, a, err = mc.QueryLockin(nil, req.Account)
	if err != nil {
		return
	}
	sr.MainChainLockIn = &lockStruct{
		SecretHash: l,
		Expiration: e.Int64(),
		Value:      a,
	}
	l, e, a, err = mc.QueryLockout(nil, req.Account)
	if err != nil {
		return
	}
	sr.MainChainLockout = &lockStruct{
		SecretHash: l,
		Expiration: e.Int64(),
		Value:      a,
	}
}
