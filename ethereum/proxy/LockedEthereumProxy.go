package proxy

import (
	"context"
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/ethereum/client"
	"github.com/SmartMeshFoundation/distributed-notary/ethereum/contracts"
	"github.com/ethereum/go-ethereum/common"
)

// LockedEthereumProxy :
type LockedEthereumProxy struct {
	Contract *contracts.LockedEthereum
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
		}
		return
	}
	err = fmt.Errorf("no code at %s", contractAddress.String())
	return
}
