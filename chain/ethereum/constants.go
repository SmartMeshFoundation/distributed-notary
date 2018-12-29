package ethereum

import "time"

// 事件/交易确认块数
var confirmBlockNumber uint64 = 17

// 轮询间隔
//var pollPeriod = 7500 * time.Millisecond
var pollPeriod = 500 * time.Millisecond

// NewBlockNumber日志间隔
var logPeriod = uint64(150)

//defaultPollTimeout  request wait time
const defaultPollTimeout = 180 * time.Second

// spectrumRPCTimeout :
var spectrumRPCTimeout = 3 * time.Second
