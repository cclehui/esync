package main

import (
	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr"
)

func main() {
	svr := esyncsvr.NewServer(config.InitConfigFromFile("./config/config.sample.yaml"))

	// daoongorm.SetGlobalCacheUtil(svr.GetRedisUtil())

	svr.Start()

}
