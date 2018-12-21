package params

import (
	"crypto/ecdsa"
	"path/filepath"

	"github.com/ethereum/go-ethereum/node"
)

//ListenIP my listen ip
var ListenIP = "127.0.0.1"

//ListenPort my listen port
var ListenPort = "18000"

//NotaryShareArg share info between all notatories.
type NotaryShareArg struct {
	Index int //我的编号
}

/*
要求2/3以上的人都同意才能生成有效签名.
*/
var ThresholdCount = 4
var ShareCount = 7

//NotatoryShareInfo share info between all notatories.
//var NotatoryShareInfo NotaryShareArg

//NotatoryInfo 公证人的基本信息
type NotatoryInfo struct {
	Name string
	Addr string //how to contact with this notary
	Key  *ecdsa.PublicKey
}

//DefaultKeyStoreDir keystore path of ethereum
func DefaultKeyStoreDir() string {
	return filepath.Join(node.DefaultDataDir(), "keystore")
}
