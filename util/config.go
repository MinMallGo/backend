package util

import (
	"errors"
	"gopkg.in/ini.v1"
	"log"
	"mall_backend/structure"
	"sync"
)

var config *structure.Config

var once sync.Once

// ConfigInstance 懂了。可以在ini的时候初始化
func ConfigInstance(filename string) {
	once.Do(func() {
		if len(filename) == 0 {
			filename = "config.ini"
		}
		log.Println(filename)
		config = &structure.Config{}
		err := ini.MapTo(config, filename)
		if err != nil {
			config = nil
			log.Println("parse config.ini error:" + err.Error())
			return
		}
	})
}

// GetRedis 返回redisConfig的数据指针
func GetRedis() (*structure.RedisConfig, error) {
	if config == nil {
		return nil, errors.New("parse redis config failed")
	}
	return &config.Redis, nil
}

// GetDatabase 返回databaseConfig的数据指针
func GetDatabase() (*structure.DataBaseConfig, error) {
	if config == nil {
		log.Println("config.Redis is nil")
		return nil, errors.New("parse database config failed")
	}
	return &config.DataBase, nil
}
