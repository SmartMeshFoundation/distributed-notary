package proxy

import (
	"context"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/client"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/contracts"
	"github.com/ethereum/go-ethereum/common"
)

// SideChainErc20TokenProxy :
type SideChainErc20TokenProxy struct {
	Contract *contracts.EthereumToken
}

// NewSideChainErc20TokenProxy :
func NewSideChainErc20TokenProxy(conn *client.SafeEthClient, tokenAddress common.Address) (p *SideChainErc20TokenProxy, err error) {

	code, err := conn.CodeAt(context.Background(), tokenAddress, nil)
	if err == nil && len(code) > 0 {
		c, err2 := contracts.NewEthereumToken(tokenAddress, conn)
		if err = err2; err != nil {
			return
		}
		p = &SideChainErc20TokenProxy{
			Contract: c,
		}
		return
	}
	err = fmt.Errorf("no code at %s", tokenAddress.String())
	return
}
