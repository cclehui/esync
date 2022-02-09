package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr"
	"github.com/cclehui/esync/esyncsvr/esyncdao"
	"github.com/cclehui/esync/esyncsvr/esyncsvc"
	"github.com/cclehui/esync/esyncsvr/handler"
)

func main() {
	esyncsvc.RegisterHandler("test_nop_handler", []esyncsvc.HandlerBase{&handler.NopHandler{}})

	go func() {
		svr := esyncsvr.NewServer(config.InitConfigFromFile("./config/config.sample.yaml"))
		svr.Start()
	}()

	time.Sleep(time.Second * 3)

	eventData := &esyncsvc.EventData{
		EventType: "test_nop_handler",
		EventData: "xxxxxxxxxxxx",
		EventOption: &esyncdao.EventOption{
			DelaySeconds: []int{1, 3},
			Persistent:   true,
		},
	}

	ctx := context.Background()

	eventSvc := &esyncsvc.EventService{}
	err := eventSvc.AddEvent(ctx, eventData)

	fmt.Println("mmmmmmmmmmmmmm:", err) // cclehui_test

	select {}

}
