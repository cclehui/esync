package service

import (
	"time"

	"git2.qingtingfm.com/podcaster/papi-go/dao"
)

// 事件数据结构
type EventData struct {
	EventType   string
	UniqKey     string
	EventData   string
	EventOption *dao.EventOption
}

type HandlerParams struct {
	EventID         int64
	EventDefaultDao *dao.EventDefaultDao
	EventData       *EventData
	Durations       []time.Duration // 添加到时间轮上执行的周期 递增
}

func (hp *HandlerParams) GetFirstAndNextDurations() (time.Duration, []time.Duration) {
	resSlice := make([]time.Duration, 0)

	if len(hp.Durations) < 1 {
		return time.Duration(0), resSlice
	}

	firstDuration := hp.Durations[0]

	for _, curDuration := range hp.Durations[1:] {
		tempDuration := curDuration - firstDuration
		resSlice = append(resSlice, tempDuration)
	}

	return firstDuration, resSlice
}
