package dao

import (
	daoongorm "github.com/cclehui/dao-on-gorm"
	"github.com/gomodule/redigo/redis"
)

var configFile = "./config_demo.yaml"

var configData *ConfigDemo
var dbClient *daoongorm.DBClient
var redisUtil *CacheUtilDemo

// db client
func GetDBClient() *daoongorm.DBClient {
	if dbClient == nil {
		initBase()
	}

	return dbClient
}

// 缓存组件
func GetRedisUtil() *CacheUtilDemo {
	if redisUtil == nil {
		initBase()
	}

	return redisUtil
}

func initBase() {
	initConfig()
	initCacheUtil()
	initDBClient()
}

func initConfig() {
	configDataTmp := &ConfigDemo{}

	_, err := configDataTmp.DecodeFromFile(configFile)
	if err != nil {
		panic(err)
	}

	configData = configDataTmp
}

func initDBClient() {
	dbClientTmp, err := daoongorm.NewDBClient(configData.Mysql.Test)
	if err != nil {
		panic(err)
	}

	dbClient = dbClientTmp
}

func initCacheUtil() {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConfig.Server)
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("AUTH", redisConfig.Password); err != nil {
				c.Close()
				return nil, err
			}

			return c, nil
		},
	}

	redisUtilTmp := NewCacheUtilDemo(configData.Redis.Default)

	redisUtil = redisUtilTmp
}
