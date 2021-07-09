// 开发环境使用,可自定义出块时间的私链
// +build !REG
// +build !PRO

package cfg

import "time"

// Env 环境名
const Env = "DEV"

func init() {
	/*
		链
	*/
	/*
		SMC = &ChainCfg{
			Name:                  "spectrum",
			BlockPeriod:           time.Second,
			ConfirmBlockNumber:    0,
			BlockNumberPollPeriod: 500 * time.Millisecond,
			BlockNumberLogPeriod:  20,
			RPCTimeout:            3 * time.Second,
			CrossFeeRate:          10000,
		}

		HECO = &ChainCfg{
			Name:                  "heco",
			BlockPeriod:           time.Second,
			ConfirmBlockNumber:    0,
			BlockNumberPollPeriod: 500 * time.Millisecond,
			BlockNumberLogPeriod:  20,
			RPCTimeout:            3 * time.Second,
			CrossFeeRate:          10000,
		}*/
	SMC = &ChainCfg{
		Name:                  "spectrum",
		BlockPeriod:           time.Second,
		ConfirmBlockNumber:    0,
		BlockNumberPollPeriod: 14 * time.Second,
		BlockNumberLogPeriod:  1,
		RPCTimeout:            10 * time.Second,
		CrossFeeRate:          10000,
	}

	HECO = &ChainCfg{
		Name:                  "heco",
		BlockPeriod:           time.Second,
		ConfirmBlockNumber:    0,
		BlockNumberPollPeriod: 3 * time.Second,
		BlockNumberLogPeriod:  1,
		RPCTimeout:            10 * time.Second,
		CrossFeeRate:          10000,
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
