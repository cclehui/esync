package esyncsvr

import (
	"fmt"
	"sync"

	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr/dao"
	"github.com/cclehui/esync/esyncsvr/esyncsvc"
)

var defaultServer *Server
var defaultServerOnce = sync.Once{}

type Server struct{}

func GetServer() *Server {
	return defaultServer
}

// 初始化server
func NewServer(configData *config.Config, options ...OptionFunc) *Server {
	defaultServerOnce.Do(func() {
		config.Conf = configData
		dao.InitStorage() // 初始化 mysql 和redis

		svr := &Server{}

		for _, optFunc := range options {
			optFunc(svr)
		}

		defaultServer = svr
	})

	return defaultServer
}

func (svr *Server) Start() {
	dao.InitDao()
	esyncsvc.InitTimeWheel()

	fmt.Println("sssssssssssss:") // cclehui_test
	select {}
}
