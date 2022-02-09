package esyncsvr

import (
	"fmt"
	"sync"

	daoongorm "github.com/cclehui/dao-on-gorm"
	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr/esyncdao"
	"github.com/cclehui/esync/esyncsvr/esyncsvc"
	"github.com/cclehui/esync/esyncutil"
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
		esyncdao.InitStorage() // 初始化 mysql 和redis

		svr := &Server{}

		for _, optFunc := range options {
			optFunc(svr)
		}

		defaultServer = svr
	})

	return defaultServer
}

func (svr *Server) Start() {
	daoongorm.SetLogger(esyncutil.GetLogger())

	esyncdao.InitDao()
	esyncsvc.InitTimeWheel()

	fmt.Println("sssssssssssss:") // cclehui_test
	select {}
}
