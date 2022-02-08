package service

import (
	"context"
	"strconv"
	"time"

	"github.com/cclehui/esync/dao"
	"github.com/cclehui/esync/esyncdefine"
	"github.com/cclehui/esync/esyncutil"
	"github.com/pkg/errors"
)

type EventService struct{}

// 事件数据结构
type EventData struct {
	EventType   string
	UniqKey     string
	EventData   string
	EventOption *dao.EventOption
}

func (svc *EventService) AddEvent(ctx context.Context, eventData *EventData) error {
	if eventData.EventType == "" {
		return errors.New("event_type is empty")
	}

	if _, ok := allHandlerSliceMap[eventData.EventType]; !ok {
		return errors.New("event_type is not registered")
	}

	eventData.EventOption = svc.getFormatedEventOption(eventData)

	eventDao, err := dao.NewEsyncEventDefaultDao(ctx, &dao.EsyncEventDefaultDao{}, false)
	if err != nil {
		return err
	}

	eventDao.EventDate, _ = strconv.Atoi(time.Now().Format("20060102"))
	eventDao.EventType = eventData.EventType
	eventDao.UniqKey = eventData.UniqKey
	eventDao.UniqKeyCRC32 = esyncutil.CRC32(eventData.UniqKey)
	eventDao.EStatus = esyncdefine.EventNew

	eventDao.SetEventOption(eventData.EventOption)
	eventDao.EventData = eventData.EventData

	handlerParams := &HandlerParams{
		EventDefaultDao: eventDao,
	}

	// 持久化
	if eventDao.EventOption.Persistent {
		err = eventDao.Save(ctx)
		if err != nil {
			return err
		}

		handlerParams.EventID = eventDao.ID
	}

	// 添加到时间轮上等待执行
	GetEventTimeWheel.AddTimer(eventDao.GetNextRetryDelayDuration(),
		eventDao.GetTimerKey(), handlerParams)

	return nil
}

func (svc *EventService) getFormatedEventOption(eventData *EventData) *dao.EventOption {
	result := eventData.EventOption

	if result == nil {
		result = &dao.EventOption{}
	}

	if result.DelaySeconds == nil || len(result.DelaySeconds) < 1 {
		result.DelaySeconds = esyncdefine.GetDefaultDelaySeconds()
	}

	if result.MaxRunNum < 1 {
		result.MaxRunNum = esyncdefine.GetDefaultMaxRunNum()
	}

	result.StartAt = time.Now().Add(time.Second * time.Duration(result.DelaySeconds[0])).Unix()

	return result
}
