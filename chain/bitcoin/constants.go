package bitcoin

import "time"

// 事件/交易确认块数
var confirmBlockNumber int32 = 6

// 轮询间隔
var pollPeriod = 500 * time.Millisecond

//var pollPeriod = 7500 * time.Millisecond

// NewBlockNumber日志间隔
//var logPeriod = int32(10)
var logPeriod = int32(20)

// btcdRPCTimeout :
var btcdRPCTimeout = 3 * time.Second
