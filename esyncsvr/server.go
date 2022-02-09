package esyncsvr

import (
	"fmt"
	"sync"

	daoongorm "github.com/cclehui/dao-on-gorm"
	"github.com/cclehui/esync/config"
	"github.com/cclehui/esync/esyncsvr/dao"
	"github.com/cclehui/esync/esyncsvr/service"
	redisutil "github.com/cclehui/redis-util"
	"github.com/gomodule/redigo/redis"
)

var defaultServer *Server
var defaultServerOnce = sync.Once{}

type Server struct {
	mysqlClient *daoongorm.DBClient
	redisPool   *redis.Pool
	redisUtil   *redisutil.RedisUtil
}

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
	service.InitTimeWheel()

	fmt.Println("sssssssssssss:") // cclehui_test
	select {}
}
