package esyncdao

import (
	"encoding/json"
	"time"

	"github.com/cclehui/esync/esyncdefine"
	"github.com/pkg/errors"
)

// 事件的option
type EventOption struct {
	DelaySeconds    []int `json:"delay_seconds"`     // 延迟执行的秒数 [5, 15, 45]:第5秒开始执行,如失败15,45秒重试
	MaxRunNum       int   `json:"max_run_num"`       // 最多执行的次数
	StartAt         int64 `json:"start_at"`          // 能开始执行的时间戳
	MaxAliveSeconds int   `json:"max_alive_seconds"` // 最大存活的时间
	Persistent      bool  `json:"persistent"`        // 是否持久化， true false
}

func (myDao *EsyncEventDefaultDao) SetEventOption(option *EventOption) {
	optionBytes, _ := json.Marshal(option)

	myDao.EventOption = string(optionBytes)
}

func (myDao *EsyncEventDefaultDao) GetEventOption() (*EventOption, error) {
	result := &EventOption{}

	err := json.Unmarshal([]byte(myDao.EventOption), result)
	if err != nil {
		return result, errors.WithStack(err)
	}

	return result, nil
}

// 获取下一次重试的delay 时间
func (myDao *EsyncEventDefaultDao) GetNextRetryDelayDuration() time.Duration {
	eventOption, err := myDao.GetEventOption()
	if err != nil {
		return esyncdefine.GetDefaultRetryDelay()
	}

	if eventOption.DelaySeconds == nil ||
		len(eventOption.DelaySeconds) < 1 {
		return esyncdefine.GetDefaultRetryDelay()
	}

	handlerInfo, err := myDao.GetHandlerInfo()
	if err != nil || handlerInfo == nil ||
		len(handlerInfo.RunTS) < 1 {
		// 返回第一个延迟时间
		return time.Second * time.Duration(eventOption.DelaySeconds[0])
	}

	if len(handlerInfo.RunTS) >= len(eventOption.DelaySeconds) {
		return esyncdefine.GetDefaultRetryDelay()
	}

	lastIndex := len(handlerInfo.RunTS) - 1
	curIndex := lastIndex + 1

	delaySeconds := eventOption.DelaySeconds[curIndex] - eventOption.DelaySeconds[lastIndex]

	if delaySeconds < 0 {
		delaySeconds = 0
	}

	return time.Second * time.Duration(delaySeconds)
}

// 事件当前时间能否执行
func (myDao *EsyncEventDefaultDao) IsCanRunNow() bool {
	eventOption, err := myDao.GetEventOption()
	if err == nil && eventOption != nil &&
		eventOption.StartAt > time.Now().Unix() { // 未到执行时间
		return false
	}

	return true
}

// 添加到时间轮上task 的key
func (myDao *EsyncEventDefaultDao) GetTimerKey() interface{} {
	if myDao.ID > 0 {
		return myDao.ID
	}

	return nil
}

// handler处理的结果
type HandlerInfo struct {
	FailCount       int      `json:"fail_count"`
	RunTS           []int64  `json:"run_ts"`           // 执行的时间戳
	SucceedHandlers []string `json:"succeed_handlers"` // 已成功的handlers
}

func (hi *HandlerInfo) GetSucceedHandlerMap() map[string]bool {
	result := make(map[string]bool)

	for _, handlerID := range hi.SucceedHandlers {
		result[handlerID] = true
	}

	return result
}

func (myDao *EsyncEventDefaultDao) GetSucceedHandlerMap() map[string]bool {
	handlerInfo, _ := myDao.GetHandlerInfo()
	result := make(map[string]bool)

	if handlerInfo != nil {
		result = handlerInfo.GetSucceedHandlerMap()
	}

	return result
}

func (myDao *EsyncEventDefaultDao) GetHandlerInfo() (*HandlerInfo, error) {
	result := &HandlerInfo{}

	if myDao.HandlerInfo == "" {
		return result, nil
	}

	err := json.Unmarshal([]byte(myDao.HandlerInfo), result)

	return result, err
}

func (myDao *EsyncEventDefaultDao) SetHandlerInfo(handlerInfo *HandlerInfo) {
	handlerInfoBytes, _ := json.Marshal(handlerInfo)

	myDao.HandlerInfo = string(handlerInfoBytes)
}
