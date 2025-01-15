package structure

import "github.com/redis/go-redis/v9"

type RedisConfig struct {
	Host         string `ini:"host"`
	Port         int    `ini:"port"`
	DB           int    `ini:"name"`
	Password     string `ini:"password"`
	MaxIdleConns int    `ini:"max_idle_conns"`
	MaxOpenConns int    `ini:"max_open_conns"`
}

type RedisClient struct {
	Client *redis.Client
}
