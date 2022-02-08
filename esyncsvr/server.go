package esyncsvr

import (
	"fmt"
	"sync"

	daoongorm "github.com/cclehui/dao-on-gorm"
	"github.com/cclehui/esync/config"
	redisutil "github.com/cclehui/redis-util"
	"github.com/gomodule/redigo/redis"
)

var defaultServer *Server
var defaultServerOnce = sync.Once{}

type Server struct {
	mysqlClient *daoongorm.DBClient
	redisPool   *redis.Pool
	redisUtil   *redisutil.RedisUtil
	configData  *config.Config
}

func GetServer() *Server {
	return defaultServer
}

// 初始化server
func NewServer(configData *config.Config, options ...OptionFunc) *Server {
	defaultServerOnce.Do(func() {
		svr := &Server{
			configData: configData,
		}

		for _, optFunc := range options {
			optFunc(svr)
		}

		svr.initMysqlClient() // mysql
		svr.initRedisUtil()   // redis

		defaultServer = svr
	})

	return defaultServer
}

func (svr *Server) Start() {

	fmt.Println("sssssssssssss:") // cclehui_test
	select {}
}

func (svr *Server) GetMysqlClient() *daoongorm.DBClient {
	return svr.mysqlClient
}

func (svr *Server) GetRedisUtil() *redisutil.RedisUtil {
	return svr.redisUtil
}

func (svr *Server) GetRedisPool() *redis.Pool {
	return svr.redisPool
}

func (svr *Server) initMysqlClient() {
	dbClientTmp, err := daoongorm.NewDBClient(svr.configData.Mysql)
	if err != nil {
		panic(err)
	}

	svr.mysqlClient = dbClientTmp
}

func (svr *Server) initRedisUtil() {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", svr.configData.Redis.Server)
			if err != nil {
				return nil, err
			}

			if svr.configData.Redis.Password != "" {
				if _, err := c.Do("AUTH", svr.configData.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}

			}

			return c, nil
		},
	}

	_, err := redisPool.Dial()
	if err != nil {
		panic(err) // 配置异常panic 无法启动
	}

	svr.redisPool = redisPool
	svr.redisUtil = redisutil.NewRedisUtil(redisPool)
}
