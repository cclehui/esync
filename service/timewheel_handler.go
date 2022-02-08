package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cclehui/esync/dao"
	"github.com/cclehui/esync/esyncdefine"
	"github.com/cclehui/esync/esyncsvr"
	"github.com/cclehui/esync/esyncutil"
	"github.com/go-redsync/redsync"
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
		esyncutil.GetLogger().Errorf(ctx, "timeWheelHandleEvent, error:%s, %+v", "事件参数类型不正确", handlerParamsIter)

		return
	}

	needRetry := handleOneEvent(ctx, handlerParams)

	// 添加到时间轮上等待再次执行 重试
	if needRetry {
		retryDelayDuration := handlerParams.EventDefaultDao.GetNextRetryDelayDuration()

		GetEventTimeWheel().AddTimer(retryDelayDuration,
			handlerParams.EventDefaultDao.GetTimerKey(), handlerParams)
	}
}

// 事件处理
// 有两个入口， 一个是从timewheel上触发， 另一个是从定时的 cron_monitor上触发
func handleOneEvent(ctx context.Context, handlerParams *HandlerParams) (needRetry bool) {
	needRetry = true
	var err error

	if handlerParams.EventDefaultDao == nil && handlerParams.EventID < 1 {
		esyncutil.GetLogger().Errorf(ctx, "handleEvent, error:%s, %+v", "事件参数不正确", handlerParams)
		needRetry = false

		return
	}

	defer func() {
		if err2 := recover(); err2 != nil {
			err = errors.New(fmt.Sprintf("recover error:%+v", err))
		}

		if err != nil {
			esyncutil.GetLogger().Errorf(ctx, "handleEvent, error:%+v, %s", err, handlerParams.LogIDStr())
		}
	}()

	// 持久化event 需要加锁 防止并发
	if handlerParams.EventID > 0 {
		lockOption := []redsync.Option{redsync.SetExpiry(time.Second * 30),
			redsync.SetTries(1)}

		redisPool := esyncsvr.GetServer().GetRedisPool()

		redisLock := redsync.New([]redsync.Pool{redisPool}).
			NewMutex(fmt.Sprintf("esync:220228_event_handle:%d", handlerParams.EventID), lockOption...)

		if err = redisLock.Lock(); err != nil {
			needRetry = true
			return
		}

		defer func() {
			_, _ = redisLock.Unlock() // 释放锁
		}()

		// 开始事件处理 每次执行需要重新load 防止重复执行 (因为状态可能已经发生了变化)
		tempDao, err2 := dao.NewEsyncEventDefaultDao(ctx,
			&dao.EsyncEventDefaultDao{ID: handlerParams.EventID}, false)
		if err2 != nil {
			err = err2
			needRetry = true

			return
		}

		if tempDao.GetDaoBase().IsNewRow() {
			err = errors.New("事件id不存在")
			needRetry = false

			return
		}

		handlerParams.EventDefaultDao = tempDao
	}

	if handlerParams.EventDefaultDao == nil {
		esyncutil.GetLogger().Errorf(ctx, "handleEvent, error, EventDefaultDao is nil")
		needRetry = false

		return
	}

	logSuffix := fmt.Sprintf("esync, handleEvent:%s", handlerParams.LogIDStr())

	if handlerParams.EventDefaultDao.EStatus != esyncdefine.EventNew {
		esyncutil.GetLogger().Infof(ctx, "事件已无需处理, %s", logSuffix)

		needRetry = false

		return
	}

	if !handlerParams.EventDefaultDao.IsCanRunNow() {
		esyncutil.GetLogger().Infof(ctx, "事件未到执行时间, %s", logSuffix)

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
			esyncutil.GetLogger().Errorf(ctx, "handler:%s, 处理异常, %+v, %s", handler.GetHandlerID(), err2, logSuffix)
			failHandlerList = append(failHandlerList, handler)
		} else {
			oldHandlerInfo.SucceedHandlers = append(oldHandlerInfo.SucceedHandlers, handler.GetHandlerID())
			esyncutil.GetLogger().Infof(ctx, "handler:%s, 处理成功, %s", handler.GetHandlerID(), logSuffix)
		}
	}

	// 运行处理的时间记录
	oldHandlerInfo.RunTs = append(oldHandlerInfo.RunTs, time.Now().Unix())

	// 事件状态处理

	if len(failHandlerList) < 1 { // 成功
		esyncutil.GetLogger().Infof(ctx, "全部处理成功, %s", logSuffix)

		needRetry = false
		handlerParams.EventDefaultDao.EStatus = esyncdefine.EventSuccess
	} else {
		err = errors.New(fmt.Sprintf("%d个handler处理失败", len(failHandlerList)))
		oldHandlerInfo.FailCount += 1

		eventOption, _ := handlerParams.EventDefaultDao.GetEventOption()

		maxAliveTimeDiff := time.Second * time.Duration(eventOption.MaxAliveSeconds)

		startTime := time.Unix(eventOption.StartAt, 0)
		failTimeDiff := time.Now().Sub(startTime)

		if oldHandlerInfo.FailCount >= eventOption.MaxRunNum ||
			(eventOption.MaxAliveSeconds > 0 && failTimeDiff >= maxAliveTimeDiff) {
			handlerParams.EventDefaultDao.EStatus = esyncdefine.EventFail
			needRetry = false // 不需要再重试了 失败超过阈值

			esyncutil.GetLogger().Errorf(ctx, "失败次数超过阈值,event处理失败:%s", logSuffix)
		} else {
			needRetry = true
		}
	}

	handlerParams.EventDefaultDao.SetHandlerInfo(oldHandlerInfo)

	if handlerParams.EventID > 0 { // 更新持久化数据状态
		_ = handlerParams.EventDefaultDao.GetDaoBase().Save(ctx)
	}

	return needRetry
}
