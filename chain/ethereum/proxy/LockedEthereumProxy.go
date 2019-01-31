package proxy

import (
	"context"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/contracts"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

// LockedEthereumProxy :
type LockedEthereumProxy struct {
	Contract *contracts.LockedEthereum
	conn     *client.SafeEthClient
}

// NewLockedEthereumProxy :
func NewLockedEthereumProxy(conn *client.SafeEthClient, contractAddress common.Address) (p *LockedEthereumProxy, err error) {
	code, err := conn.CodeAt(context.Background(), contractAddress, nil)
	if err == nil && len(code) > 0 {
		c, err2 := contracts.NewLockedEthereum(contractAddress, conn)
		if err = err2; err != nil {
			return
		}
		p = &LockedEthereumProxy{
			Contract: c,
			conn:     conn,
		}
		return
	}
	err = fmt.Errorf("no code at %s", contractAddress.String())
	return
}

// QueryLockin impl chain.ContractProxy
func (p *LockedEthereumProxy) QueryLockin(accountHex string) (secretHash common.Hash, expiration uint64, amount *big.Int, err error) {
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
func (p *LockedEthereumProxy) QueryLockout(accountHex string) (secretHash common.Hash, expiration uint64, amount *big.Int, err error) {
	account := common.HexToAddress(accountHex)
	var cExpiration *big.Int
	secretHash, cExpiration, amount, err = p.Contract.QueryLockout(nil, account)
	if err != nil {
		return
	}
	expiration = cExpiration.Uint64()
	return
}

// PrepareLockin : impl chain.ContractProxy
// 主链的PrepareLockin由用户发起,不需要使用accountHex参数,传""即可
func (p *LockedEthereumProxy) PrepareLockin(opts *bind.TransactOpts, accountHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (err error) {
	opts.Value = amount
	expiration2 := new(big.Int).SetUint64(expiration)
	tx, err := p.Contract.PrepareLockin(opts, secretHash, expiration2)
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
func (p *LockedEthereumProxy) Lockin(opts *bind.TransactOpts, accountHex string, secret common.Hash) (err error) {
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
func (p *LockedEthereumProxy) CancelLockin(opts *bind.TransactOpts, accountHex string) (err error) {
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

// PrepareLockout : impl chain.ContractProxy
func (p *LockedEthereumProxy) PrepareLockout(opts *bind.TransactOpts, accountHex string, secretHash common.Hash, expiration uint64, amount *big.Int) (err error) {
	//opts.Value = amount
	account := common.HexToAddress(accountHex)
	expiration2 := new(big.Int).SetUint64(expiration)
	tx, err := p.Contract.PrepareLockoutHTLC(opts, account, secretHash, expiration2, amount)
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
func (p *LockedEthereumProxy) Lockout(opts *bind.TransactOpts, accountHex string, secret common.Hash) (err error) {
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
func (p *LockedEthereumProxy) CancelLockout(opts *bind.TransactOpts, accountHex string) (err error) {
	account := common.HexToAddress(accountHex)
	var tx *types.Transaction
	tx, err = p.Contract.CancleLockOut(opts, account)
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
