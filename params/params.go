package params

import (
	"path/filepath"

	"github.com/ethereum/go-ethereum/node"
)

// SCTokenNameSuffix 侧链Token名统一后缀
var SCTokenNameSuffix = "-AtmosphereToken"

// ForkConfirmNumber : 分叉确认块数量,BlockNumber < 最新块-ForkConfirmNumber的事件被认为无分叉的风险
var ForkConfirmNumber uint64 = 17

// MinLockinSCExpiration : lockin侧链最小超时时间
var MinLockinSCExpiration uint64 = 300

// MinLockinMCExpiration : lockin主链最小超时时间
var MinLockinMCExpiration = MinLockinSCExpiration + 5*ForkConfirmNumber + 1

// MinLockoutMCExpiration : lockout主链最小超时时间
var MinLockoutMCExpiration uint64 = 300

// MinLockoutSCExpiration : lockout侧链最小超时时间
var MinLockoutSCExpiration = MinLockoutMCExpiration + 5*ForkConfirmNumber + 1

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
