package params

import (
	"path/filepath"

	"github.com/ethereum/go-ethereum/node"
)

//ListenIP my listen ip
var ListenIP = "127.0.0.1"

//ListenPort my listen port
var ListenPort = "18000"

// SCTokenNameSuffix 侧链Token名统一后缀
var SCTokenNameSuffix = "-AtmosphereToken"

/*
ThresholdCount 要求2/3以上的人都同意才能生成有效签名.
*/
var ThresholdCount = 4

// ShareCount :
var ShareCount = 7

//DefaultKeyStoreDir keystore path of ethereum
func DefaultKeyStoreDir() string {
	return filepath.Join(node.DefaultDataDir(), "keystore")
}
