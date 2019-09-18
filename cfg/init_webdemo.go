// 测试环境使用,连接测试链
// +build REG

package cfg

import "time"

// Env 环境名
const Env = "REG"

func init() {
	/*
		链
	*/
	SMC = &ChainCfg{
		Name:                  "spectrum",
		BlockPeriod:           time.Second,
		ConfirmBlockNumber:    0,
		BlockNumberPollPeriod: 14 * time.Millisecond,
		BlockNumberLogPeriod:  20,
		RPCTimeout:            3 * time.Second,
		CrossFeeRate:          10000,
	}

	ETH = &ChainCfg{
		Name:                  "ethereum",
		BlockPeriod:           time.Second,
		ConfirmBlockNumber:    0,
		BlockNumberPollPeriod: 14 * time.Millisecond,
		BlockNumberLogPeriod:  20,
		RPCTimeout:            3 * time.Second,
		CrossFeeRate:          10000,
	}

	BTC = &ChainCfg{
		Name:               "bitcoin",
		BlockPeriod:        time.Second,
		ConfirmBlockNumber: 0,
		//BlockNumberPollPeriod: 500 * time.Millisecond, // BTC不使用该参数
		BlockNumberLogPeriod: 20,
		//RPCTimeout:           3 * time.Second, // BTC不使用该参数
		CrossFeeRate: 10000,
	}
	/*
		跨链
	*/
	minExpiration4User := time.Minute * 10 // 开发环境取10分钟
	Cross = &CrossCfg{
		MinExpirationTime4User:   minExpiration4User,
		MinExpirationTime4Notary: minExpiration4User / 2,
	}
}
