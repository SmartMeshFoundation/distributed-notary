package cfg

import (
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
)

// NotaryCfg 公证人相关参数
type NotaryCfg struct {
	ShareCount     int // 公证人总人数
	ThresholdCount int // 参与一次签名的最少人数
	Self           *models.NotaryInfo
	Others         []*models.NotaryInfo
}

// Notaries :
var Notaries *NotaryCfg

// ChainCfg 区块链相关参数
type ChainCfg struct {
	Name                  string        // 链名
	BlockPeriod           time.Duration // 出块间隔
	ConfirmBlockNumber    uint64        // 事件/交易确认块数
	BlockNumberPollPeriod time.Duration // 新块查询轮询间隔
	BlockNumberLogPeriod  uint64        // 块号日志打印间隔
	RPCTimeout            time.Duration // 区块链节点RPC接口超时时间
	CrossFeeRate          int64         // 跨链手续费费率,fee = amount/crossFeeRate
}

// SMC Spectrum链相关参数
var SMC *ChainCfg

// BTC Bitcoin链相关参数
var BTC *ChainCfg

// ETH Ethereum链相关参数
var ETH *ChainCfg

// CrossCfg 跨链相关参数
type CrossCfg struct {
	MinExpirationTime4User   time.Duration // 用户prepareLockin或者prepareLockout的最小过期时间,低于该值的请求/事件将被忽略
	MinExpirationTime4Notary time.Duration // 取MinExpirationTime4User的一半
}

// Cross :
var Cross *CrossCfg

// RegisterNotaries 注册公证人的基础信息,需在启动的时候调用,且只调用一次
func RegisterNotaries(notaries []*models.NotaryInfo, selfAddress common.Address) {
	if Notaries != nil {
		panic("wrong code")
	}
	Notaries = &NotaryCfg{}
	// 根据notaries数量初始化 ShareCount及ThresholdCount
	Notaries.ShareCount = len(notaries)
	Notaries.ThresholdCount = Notaries.ShareCount / 3 * 2
	if Notaries.ShareCount%3 > 1 {
		Notaries.ThresholdCount++
	}
	for _, notary := range notaries {
		if notary.GetAddress() == selfAddress {
			Notaries.Self = notary
		} else {
			Notaries.Others = append(Notaries.Others, notary)
		}
	}
}
