package service

import (
	"context"
	"sync"
)

var allHandlerSliceMap = make(map[string][]HandlerBase)
var handlerRegisterMutex = sync.Mutex{}

// 注册handler
func RegisterHandler(eventType string, handlerSlice []HandlerBase) {
	handlerRegisterMutex.Lock()
	allHandlerSliceMap[eventType] = handlerSlice
	handlerRegisterMutex.Unlock()
}

type HandlerBase interface {
	GetHandlerID() string // 建议日期加三位数字 比如 210916002
	Do(ctx context.Context, params *HandlerParams) error
}

func GetHandlerList(handlerParams *HandlerParams) []HandlerBase {
	result := make([]HandlerBase, 0)

	eventType := handlerParams.EventDefaultDao.EventType
	if handlerList, ok := allHandlerSliceMap[eventType]; ok {
		result = append(result, handlerList...)
	}

	return result
}
