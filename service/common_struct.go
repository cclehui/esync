package service

import (
	"fmt"

	"github.com/cclehui/esync/dao"
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
	EventDefaultDao *dao.EsyncEventDefaultDao
}

func (hp *HandlerParams) LogIDStr() string {
	if hp.EventID > 0 {
		return fmt.Sprintf("event_id:%d", hp.EventID)
	}

	return fmt.Sprintf("event_type:%s, uniqkey:%s",
		hp.EventDefaultDao.EventType, hp.EventDefaultDao.UniqKey)
}
