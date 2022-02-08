package esyncdefine

import "time"

// 默认间隔一秒
var defaultRetryDelay = time.Second

func SetDefaultRetryDelay(delay time.Duration) { defaultRetryDelay = delay }

func GetDefaultRetryDelay() time.Duration { return defaultRetryDelay }

// 默认最大run 的次数

var defaultMaxRunNum = 5

func SetDefaultMaxRunNum(value int) { defaultMaxRunNum = value }
func GetDefaultMaxRunNum() int      { return defaultMaxRunNum }

// 默认的运行和重试时间点
var defaultDelaySeconds = []int{0, 5, 15, 35, 55}

func SetDefaultDelaySeconds(value []int) { defaultDelaySeconds = value }
func GetDefaultDelaySeconds() []int      { return defaultDelaySeconds }
