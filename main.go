package main

import (
	"context"
	"time"

	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr"
	"github.com/cclehui/esync/esyncsvr/esyncdao"
	"github.com/cclehui/esync/esyncsvr/esyncsvc"
	"github.com/cclehui/esync/esyncsvr/handler"
)

const (
	EtypeTestNop  = "esync_test_nop"
	EtypeTestFail = "esync_test_fail"
)

func main() {
	esyncsvc.RegisterHandler(EtypeTestNop, []esyncsvc.HandlerBase{&handler.NopHandler{}})
	esyncsvc.RegisterHandler(EtypeTestFail,
		[]esyncsvc.HandlerBase{&handler.FailHandler{FailNum: 30}})

	go func() {
		svr := esyncsvr.NewServer(config.InitConfigFromFile("./config/config.sample.yaml"))
		svr.Start()
	}()

	time.Sleep(time.Second * 3)

	eventSvc := &esyncsvc.EventService{}
	ctx := context.Background()

	eventData := &esyncsvc.EventData{
		EventType: EtypeTestNop,
		EventData: "xxxxxxxxxxxx",
		EventOption: &esyncdao.EventOption{
			DelaySeconds: []int{1, 3},
			Persistent:   true,
		},
	}

	_ = eventSvc.AddEvent(ctx, eventData)

	eventData2 := &esyncsvc.EventData{
		EventType: EtypeTestFail,
		EventData: "ffffffffffff",
		EventOption: &esyncdao.EventOption{
			DelaySeconds: []int{0, 3, 5, 7, 9, 11},
			Persistent:   true,
		},
	}

	_ = eventSvc.AddEvent(ctx, eventData2)

	select {}
}
