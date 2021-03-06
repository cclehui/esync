package esyncsvc

import (
	"fmt"

	"github.com/cclehui/esync/esyncsvr/esyncdao"
)

type HandlerParams struct {
	EventID         int64
	EventDefaultDao *esyncdao.EsyncEventDefaultDao
}

func (hp *HandlerParams) LogIDStr() string {
	if hp.EventID > 0 {
		return fmt.Sprintf("event_id:%d", hp.EventID)
	}

	return fmt.Sprintf("event_type:%s, uniqkey:%s",
		hp.EventDefaultDao.EventType, hp.EventDefaultDao.UniqKey)
}
