package client

import (
	"context"
	"errors"
	"math/big"
	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ErrNotConnected = errors.New("eth not connected")

//SafeEthClient :
type SafeEthClient struct {
	*ethclient.Client
	lock sync.Mutex
}

//NewSafeClient create safeclient
func NewSafeClient(client *ethclient.Client) *SafeEthClient {
	c := &SafeEthClient{
		Client: client,
	}
	return c
}

//BlockByHash wrapper of BlockByHash
func (c *SafeEthClient) BlockByHash(ctx context.Context, hash common.Hash) (r1 *types.Block, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	r1, err = c.Client.BlockByHash(ctx, hash)
	return
}

//BlockByNumber wrapper of BlockByNumber
func (c *SafeEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.BlockByNumber(ctx, number)
}

// HeaderByHash returns the block header with the given hash.
func (c *SafeEthClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.HeaderByHash(ctx, hash)
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (c *SafeEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.HeaderByNumber(ctx, number)
}

//TransactionByHash wrapper of TransactionByHash
func (c *SafeEthClient) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, false, ErrNotConnected
	}
	return c.Client.TransactionByHash(ctx, hash)
}

//TransactionSender wrapper of TransactionSender
func (c *SafeEthClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return common.Address{}, ErrNotConnected
	}
	return c.Client.TransactionSender(ctx, tx, block, index)
}

// TransactionCount returns the total number of transactions in the given block.
func (c *SafeEthClient) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, ErrNotConnected
	}
	return c.Client.TransactionCount(ctx, blockHash)
}

//TransactionInBlock wrapper of TransactionInBlock
func (c *SafeEthClient) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.TransactionInBlock(ctx, blockHash, index)
}

//TransactionReceipt wrappper of TransactionReceipt
func (c *SafeEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.TransactionReceipt(ctx, txHash)
}

//SyncProgress wrapper of SyncProgress
func (c *SafeEthClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.SyncProgress(ctx)
}

//SubscribeNewHead wrapper of SubscribeNewHead
func (c *SafeEthClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.SubscribeNewHead(ctx, ch)
}

//NetworkID wrapper of NetworkID
func (c *SafeEthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.NetworkID(ctx)
}

//BalanceAt wrapper of BalanceAt
func (c *SafeEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.BalanceAt(ctx, account, blockNumber)
}

//StorageAt wrapper of StorageAt
func (c *SafeEthClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.StorageAt(ctx, account, key, blockNumber)
}

//CodeAt wrapper of CodeAt
func (c *SafeEthClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.CodeAt(ctx, account, blockNumber)
}

//NonceAt wrapper of NonceAt
func (c *SafeEthClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, ErrNotConnected
	}
	return c.Client.NonceAt(ctx, account, blockNumber)
}

//FilterLogs wrapper of FilterLogs
func (c *SafeEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.FilterLogs(ctx, q)
}

//SubscribeFilterLogs wrapper of SubscribeFilterLogs
func (c *SafeEthClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.SubscribeFilterLogs(ctx, q, ch)
}

//PendingBalanceAt wrapper of PendingBalanceAt
func (c *SafeEthClient) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.PendingBalanceAt(ctx, account)
}

//PendingStorageAt wrapper of PendingStorageAt
func (c *SafeEthClient) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.PendingStorageAt(ctx, account, key)
}

//PendingCodeAt wrapper of PendingCodeAt
func (c *SafeEthClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.PendingCodeAt(ctx, account)
}

//PendingNonceAt wrapper of PendingNonceAt
// 考虑到短时间内并发调用合约出现nonce相同导致调用失败的问题,在这里获取可用nonce的时候,加入了缓冲机制
func (c *SafeEthClient) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, ErrNotConnected
	}
	nonce, err = c.Client.PendingNonceAt(ctx, account)
	return
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (c *SafeEthClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, ErrNotConnected
	}
	return c.Client.PendingTransactionCount(ctx)
}

//CallContract wrapper of CallContract
func (c *SafeEthClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.CallContract(ctx, msg, blockNumber)
}

//PendingCallContract wrapper of PendingCallContract
func (c *SafeEthClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.PendingCallContract(ctx, msg)
}

//SuggestGasPrice wrapper of SuggestGasPrice
func (c *SafeEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, ErrNotConnected
	}
	return c.Client.SuggestGasPrice(ctx)
}

//EstimateGas wrapper of EstimateGas
func (c *SafeEthClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, ErrNotConnected
	}
	return c.Client.EstimateGas(ctx, msg)
}

//SendTransaction wrapper of SendTransaction
func (c *SafeEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return ErrNotConnected
	}
	return c.Client.SendTransaction(ctx, tx)
}

// GenesisBlockHash :
func (c *SafeEthClient) GenesisBlockHash(ctx context.Context) (genesisBlockHash common.Hash, err error) {

	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return utils.EmptyHash, ErrNotConnected
	}
	genesisBlockHead, err := c.Client.HeaderByNumber(ctx, big.NewInt(1))
	if err != nil {
		return
	}
	return genesisBlockHead.Hash(), nil
}
