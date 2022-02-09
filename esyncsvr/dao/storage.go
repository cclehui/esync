package dao

import (
	daoongorm "github.com/cclehui/dao-on-gorm"
	"github.com/cclehui/esync/config"
	"github.com/gomodule/redigo/redis"
)

var storage *Storage

func GetStorage() *Storage {
	return storage
}

func InitStorage() {
	storage = &Storage{}

	storage.initMysqlClient()
	storage.initRedis()
}

type Storage struct {
	mysqlClient *daoongorm.DBClient
	redisPool   *redis.Pool
}

func (s *Storage) GetMysqlClient() *daoongorm.DBClient {
	return s.mysqlClient
}

func (s *Storage) GetRedisPool() *redis.Pool {
	return s.redisPool
}

func (s *Storage) initMysqlClient() {
	dbClientTmp, err := daoongorm.NewDBClient(config.Conf.Mysql)
	if err != nil {
		panic(err)
	}

	s.mysqlClient = dbClientTmp
}

func (s *Storage) initRedis() {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Conf.Redis.Server)
			if err != nil {
				return nil, err
			}

			if config.Conf.Redis.Password != "" {
				if _, err := c.Do("AUTH", config.Conf.Redis.Password); err != nil {
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

	s.redisPool = redisPool
}
