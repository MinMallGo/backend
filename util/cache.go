package util

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"mall_backend/structure"
	"sync"
	"time"
)

// 写一个redis的单例模式
var (
	redisX    *structure.RedisClient
	redisOnce sync.Once
)

func RedisInstance() (*structure.RedisClient, error) {
	redisOnce.Do(func() {
		redisX = conn()
	})

	if redisX == nil {
		return nil, fmt.Errorf("failed to initialize Redis client")
	}

	return redisX, nil
}

func conn() *structure.RedisClient {
	config, err := GetRedis()
	if err != nil {
		log.Fatalln(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Println("Failed to connect to Redis:", err)
		return nil
	}
	log.Println("Redis connection initialized")
	return &structure.RedisClient{Client: client}
}

// CacheClient 直接返回缓存实例
func CacheClient() *redis.Client {
	return redisX.Client
}
