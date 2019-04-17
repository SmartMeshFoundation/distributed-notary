package bitcoin

import (
	"path/filepath"
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

var TestRPCHost = "192.168.124.13:18556"
var TestRPCUser = "wuhan"
var TestRPCPass = "wuhan"
var TestCertFilePath = filepath.Join("/home/chuck/.btcd", "rpc.cert")

func TestNormalTransfer(t *testing.T) {
	_, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	assert.Empty(t, err)
}

func TestWallet(t *testing.T) {
	bs, err := NewBTCService(TestRPCHost, TestRPCUser, TestRPCPass, TestCertFilePath)
	assert.Empty(t, err)
	m, err := bs.c.GetBlockChainInfo()
	assert.Empty(t, err)
	fmt.Println(utils.ToJSONStringFormat(m))
}
