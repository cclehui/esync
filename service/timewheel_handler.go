package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git2.qingtingfm.com/infra/qt-boot/pkg/base/ctime"
	"git2.qingtingfm.com/infra/qt-boot/pkg/log"
	"git2.qingtingfm.com/infra/qt-boot/pkg/net/redlock"
	"git2.qingtingfm.com/infra/qt-boot/pkg/net/sentry"
	"git2.qingtingfm.com/podcaster/papi-go/dao"
	"git2.qingtingfm.com/podcaster/papi-go/define"
	"git2.qingtingfm.com/podcaster/papi-go/global"
	"github.com/pkg/errors"
)

var eventTimeWheel *TimeWheel
var eventTimeWheelSync sync.Once

// 初始化
func InitTimeWheel() {
	eventTimeWheelSync.Do(func() {
		eventTimeWheel = NewTimeWheel(time.Second*1, 3600, timeWheelHandleEvent)
		eventTimeWheel.Start()
	})
}

func GetEventTimeWheel() *TimeWheel {
	return eventTimeWheel
}

// 从时间轮上触发执行
func timeWheelHandleEvent(handlerParamsIter interface{}) {
	handlerParams, ok := handlerParamsIter.(*HandlerParams)
	ctx := context.Background()

	if !ok {
		log.Errorc(ctx, "timeWheelHandleEvent, error:%s, %+v", "事件参数类型不正确", handlerParamsIter)

		return
	}

	err, needRetry := handleOneEvent(ctx, handlerParams)

	if err == nil {
		return
	}

	// 添加到时间轮上等待再次执行 重试
	if len(handlerParams.Durations) > 0 && needRetry {
		firstDuration, nextDurationSlice := handlerParams.GetFirstAndNextDurations()
		handlerParams.Durations = nextDurationSlice
		handlerParams.EventDefaultDao = nil // 让dao重新load一次

		GetEventTimeWheel().AddTimer(firstDuration, handlerParams.EventID, handlerParams)
	}
}

// 事件处理
// 有两个入口， 一个是从timewheel上触发， 另一个是从定时的 cron_monitor上触发
func handleOneEvent(ctx context.Context, handlerParams *HandlerParams) (err error, needRetry bool) {
	needRetry = true

	if handlerParams.EventData == nil || handlerParams.EventID < 1 {
		log.Errorc(ctx, "handleEvent, error:%s, %+v", "事件参数不正确", handlerParams)

		return
	}

	defer func() {
		if err2 := recover(); err2 != nil {
			err = errors.New(fmt.Sprintf("recover error:%+v", err))
		}

		if err != nil {
			log.Errorc(ctx, "handleEvent, error:%+v, %+v", err, handlerParams)
		}
	}()

	// 锁住当前事件 防止并发
	redisPool := global.GetDao().GetRedisDefault()
	configLock := &redlock.Config{
		ExpiryTime: ctime.Duration(time.Second * 30),
		Tries:      1,
	}

	redisLock := redlock.New(configLock, redisPool).
		NewMutex(fmt.Sprintf("210916_event_handle:%d", handlerParams.EventID))
	if err = redisLock.Lock(ctx); err != nil {
		return
	}

	defer func() {
		_ = redisLock.Unlock(ctx) // 释放锁
	}()

	// 开始事件处理 每次执行需要重新load 防止重复执行 (因为状态可能已经发生了变化)
	tempDao, err2 := dao.NewEventDefaultDao(ctx, &dao.EventDefaultDao{ID: handlerParams.EventID}, false)
	if err2 != nil {
		err = err2

		return
	}

	if tempDao.IsNewRow() {
		err = errors.New("事件id不存在")
		needRetry = false

		return
	}

	handlerParams.EventDefaultDao = tempDao

	logSuffix := fmt.Sprintf("event:%d", handlerParams.EventID)

	if handlerParams.EventDefaultDao.EStatus != define.EventNew {
		log.Infoc(ctx, "事件已无需处理, %s", logSuffix)

		needRetry = false

		return
	}

	if !handlerParams.EventDefaultDao.CanEventHandleNow() {
		log.Infoc(ctx, "事件未到执行时间, %s", logSuffix)

		needRetry = true

		return
	}

	// 已成功处理过的
	oldHandlerInfo, _ := handlerParams.EventDefaultDao.GetHandlerInfo()
	succeedHandlerMap := oldHandlerInfo.GetSucceedHandlerMap()

	handlerList := GetHandlerList(handlerParams)

	failHandlerList := make([]HandlerBase, 0)

	for _, handler := range handlerList {
		if succeedHandlerMap[handler.GetHandlerID()] {
			continue // 已处理过
		}

		err2 := handler.Do(ctx, handlerParams)

		if err2 != nil {
			log.Errorc(ctx, "handler:%s, 处理异常, %+v, %s", handler.GetHandlerID(), err2, logSuffix)
			failHandlerList = append(failHandlerList, handler)
		} else {
			oldHandlerInfo.SucceedHandlers = append(oldHandlerInfo.SucceedHandlers, handler.GetHandlerID())
			log.Infoc(ctx, "handler:%s, 处理成功, %s", handler.GetHandlerID(), logSuffix)
		}
	}

	// 事件状态处理

	if len(failHandlerList) < 1 { // 成功
		log.Infoc(ctx, "全部处理成功, %s", logSuffix)

		needRetry = false
		handlerParams.EventDefaultDao.EStatus = define.EventSuccess
	} else {
		err = errors.New(fmt.Sprintf("%d个handler处理失败", len(failHandlerList)))
		oldHandlerInfo.FailCount += 1

		failTimeDiff := time.Now().Sub(handlerParams.EventDefaultDao.CreatedAt)

		if failTimeDiff > time.Second*3600*12 ||
			(oldHandlerInfo.FailCount >= 10 && failTimeDiff >= time.Second*3600*2) {
			handlerParams.EventDefaultDao.EStatus = define.EventFail
			needRetry = false // 不需要再重试了 失败超过阈值

			// 发sentry 报警
			sentry.CaptureWithTags(ctx, errors.New(fmt.Sprintf("失败次数超过阈值,event处理失败:%s", logSuffix)))
		}
	}

	handlerParams.EventDefaultDao.SetHandlerInfo(oldHandlerInfo)

	_ = handlerParams.EventDefaultDao.Save(ctx)

	return err, needRetry
}
