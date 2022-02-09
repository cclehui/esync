package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr"
	"github.com/cclehui/esync/esyncsvr/dao"
	"github.com/cclehui/esync/esyncsvr/service"
)

func main() {
	service.RegisterHandler("test_nop", []service.HandlerBase{&service.HandlerNop{}})

	go func() {
		svr := esyncsvr.NewServer(config.InitConfigFromFile("./config/config.sample.yaml"))
		svr.Start()
	}()

	time.Sleep(time.Second * 3)

	eventData := &service.EventData{
		EventType: "test_nop",
		EventData: "xxxxxxxxxxxx",
		EventOption: &dao.EventOption{
			DelaySeconds: []int{1, 3},
			Persistent:   true,
		},
	}

	ctx := context.Background()

	eventSvc := &service.EventService{}
	err := eventSvc.AddEvent(ctx, eventData)

	fmt.Println("mmmmmmmmmmmmmm:", err) // cclehui_test

}
