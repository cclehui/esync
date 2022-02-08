package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/dao"
	"github.com/cclehui/esync/esyncsvr"
	"github.com/cclehui/esync/service"
)

func main() {
	svr := esyncsvr.NewServer(config.InitConfigFromFile("./config/config.sample.yaml"))

	// daoongorm.SetGlobalCacheUtil(svr.GetRedisUtil())
	// 启动 time_wheel
	service.InitTimeWheel()

	service.RegisterHandler("test_nop", &service.HandlerNop{})

	go func() {
		time.Sleep(time.Second * 3)

		eventData := &service.EventData{
			EventType: "test_nop",
			EventData: "xxxxxxxxxxxx",
			EventOption: &dao.EventOption{
				DelaySeconds: []int{1, 3},
				Persistent:   false,
			},
		}

		ctx := context.Background()

		eventSvc := &service.EventService{}
		err := eventSvc.AddEvent(ctx, eventData)

		fmt.Println("mmmmmmmmmmmmmm:", err) // cclehui_test

	}()

	svr.Start()

}
