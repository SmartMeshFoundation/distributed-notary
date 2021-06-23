package proxy

import (
	"context"

	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain/heco/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/heco/contracts"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

// SideChainErc20TokenProxy :
type SideChainErc20TokenProxy struct {
	Contract *contracts.HecoToken
	conn     *client.SafeEthClient
}

// NewSideChainErc20TokenProxy :
func NewSideChainErc20TokenProxy(conn *client.SafeEthClient, tokenAddress common.Address) (p *SideChainErc20TokenProxy, err error) {
	code, err := conn.CodeAt(context.Background(), tokenAddress, nil)
	if err == nil && len(code) > 0 {
		c, err2 := contracts.NewHecoToken(tokenAddress, conn)
		if err = err2; err != nil {
			return
		}
		p = &SideChainErc20TokenProxy{
			Contract: c,
			conn:     conn,
		}
		return
	}
	err = fmt.Errorf("no code at %s", tokenAddress.String())
	return
}

// QueryLockin impl chain.ContractProxy
func (p *SideChainErc20TokenProxy) QueryLockin(accountHex string) (secretHash common.Hash, expiration uint64, amount *big.Int, err error) {
	account := common.HexToAddress(accountHex)
	var cExpiration *big.Int
	secretHash, cExpiration, amount, err = p.Contract.QueryLockin(nil, account)
	if err != nil {
		return
	}
	expiration = cExpiration.Uint64()
	return
}

// QueryLockout impl chain.ContractProxy
func (p *SideChainErc20TokenProxy) QueryLockout(accountHex string) (secretHash common.Hash, expiration uint64, amount *big.Int, err error) {
	account := common.HexToAddress(accountHex)
	var cExpiration *big.Int
	secretHash, cExpiration, amount, err = p.Contract.QueryLockout(nil, account)
	if err != nil {
		return
	}
	expiration = cExpiration.Uint64()
	return
}

// PrepareLockin impl chain.ContractProxy
func (p *SideChainErc20TokenProxy) PrepareLockin(opts *bind.TransactOpts, accountHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (err error) {
	account := common.HexToAddress(accountHex)
	expiration2 := new(big.Int).SetUint64(expiration)
	var tx *types.Transaction
	tx, err = p.Contract.PrepareLockin(opts, account, secretHash, expiration2, amount)
	if err != nil {
		return
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.conn, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract PrepareLockin success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}

// Lockin impl chain.ContractProxy
func (p *SideChainErc20TokenProxy) Lockin(opts *bind.TransactOpts, accountHex string, secret common.Hash) (err error) {
	account := common.HexToAddress(accountHex)
	var tx *types.Transaction
	tx, err = p.Contract.Lockin(opts, account, secret)
	if err != nil {
		return
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.conn, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract Lockin success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}

// CancelLockin impl chain.ContractProxy
func (p *SideChainErc20TokenProxy) CancelLockin(opts *bind.TransactOpts, accountHex string) (err error) {
	account := common.HexToAddress(accountHex)
	var tx *types.Transaction
	tx, err = p.Contract.CancelLockin(opts, account)
	if err != nil {
		return
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.conn, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract CancelLockin success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}

// PrepareLockout impl chain.ContractProxy
// 侧链的PrepareLockout由用户发起,不需要使用accountHex参数,传""即可
func (p *SideChainErc20TokenProxy) PrepareLockout(opts *bind.TransactOpts, accountHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (err error) {
	expiration2 := new(big.Int).SetUint64(expiration)
	var tx *types.Transaction
	tx, err = p.Contract.PrepareLockout(opts, secretHash, expiration2, amount)
	if err != nil {
		return
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.conn, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract PrepareLockout success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}

// Lockout impl chain.ContractProxy
func (p *SideChainErc20TokenProxy) Lockout(opts *bind.TransactOpts, accountHex string, secret common.Hash) (err error) {
	account := common.HexToAddress(accountHex)
	var tx *types.Transaction
	tx, err = p.Contract.Lockout(opts, account, secret)
	if err != nil {
		return
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.conn, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract Lockout success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}

// CancelLockout impl chain.ContractProxy
func (p *SideChainErc20TokenProxy) CancelLockout(opts *bind.TransactOpts, accountHex string) (err error) {
	account := common.HexToAddress(accountHex)
	var tx *types.Transaction
	tx, err = p.Contract.CancelLockOut(opts, account)
	if err != nil {
		return
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.conn, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract CancelLockout success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}
