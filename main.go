package main

import (
	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr"
)

func main() {
	svr := esyncsvr.GetServer(config.InitConfigFromFile("./config/config.sample.yaml"))

	svr.Start()

}
