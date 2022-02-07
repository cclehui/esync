package dao

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// 事件的option
type EventOption struct {
	DelaySeconds  []int `json:"delay_seconds"`   // 延迟执行的秒数
	StartAt       int64 `json:"start_at"`        // 能开始执行的时间戳
	MaxRetryCount int   `json:"max_retry_count"` // 最多retry执行的次数
}

func (myDao *EsyncEventDefaulDao) GetEventOption() (*EventOption, error) {
	result := &EventOption{}

	err := json.Unmarshal([]byte(myDao.EventOption), result)
	if err != nil {
		return result, errors.WithStack(err)
	}

	return result, nil
}

// 事件当前时间能否执行
func (myDao *EsyncEventDefaulDao) IsCanRunNow() bool {
	eventOption, err := myDao.GetEventOption()
	if err == nil && eventOption != nil &&
		eventOption.StartAt > time.Now().Unix() { // 未到执行时间
		return false
	}

	return true
}

// handler处理的结果
type HandlerInfo struct {
	FailCount       int      `json:"fail_count"`
	RunTs           []int64  `json:"run_ts"`           // 执行的时间戳
	SucceedHandlers []string `json:"succeed_handlers"` // 已成功的handlers
}

func (hi *HandlerInfo) GetSucceedHandlerMap() map[string]bool {
	result := make(map[string]bool)

	for _, handlerID := range hi.SucceedHandlers {
		result[handlerID] = true
	}

	return result
}

func (myDao *EsyncEventDefaulDao) GetSucceedHandlerMap() map[string]bool {
	handlerInfo, _ := myDao.GetHandlerInfo()
	result := make(map[string]bool)

	if handlerInfo != nil {
		result = handlerInfo.GetSucceedHandlerMap()
	}

	return result
}

func (myDao *EsyncEventDefaulDao) GetHandlerInfo() (*HandlerInfo, error) {
	result := &HandlerInfo{}

	if myDao.HandlerInfo == "" {
		return result, nil
	}

	err := json.Unmarshal([]byte(myDao.HandlerInfo), result)

	return result, err
}

func (myDao *EsyncEventDefaulDao) SetHandlerInfo(handlerInfo *HandlerInfo) {
	handlerInfoBytes, _ := json.Marshal(handlerInfo)

	myDao.HandlerInfo = string(handlerInfoBytes)
}
