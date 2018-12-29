package params

import (
	"path/filepath"

	"github.com/ethereum/go-ethereum/node"
)

//ListenIP my listen ip
var ListenIP = "127.0.0.1"

//ListenPort my listen port
var ListenPort = "18000"

/*
ThresholdCount 要求2/3以上的人都同意才能生成有效签名.
*/
var ThresholdCount = 1

// ShareCount :
var ShareCount = 2

//DefaultKeyStoreDir keystore path of ethereum
func DefaultKeyStoreDir() string {
	return filepath.Join(node.DefaultDataDir(), "keystore")
}
