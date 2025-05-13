package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"mall_backend/response"
	"mall_backend/util"
	"time"
)

const (
	IPLimiterKey   string = "Limiter:IP:"
	PeerIPQPSLimit        = 11 // 为什么要设置成11？因为一开始incr就会增加1
)

var luaScript = redis.NewScript(`
	local current = redis.call('incr', KEYS[1]); 
	if tonumber(current) == 1 
		then redis.call('expire', KEYS[1], 1) 
	end; 
	return current
`)

// IPLimiter 创建一个IP限流器
func IPLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}
		ip = IPLimiterKey + ip

		ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
		defer cancel()

		counter, err := luaScript.Run(ctx, util.CacheClient(), []string{ip}).Result()
		cnt, ok := counter.(int64)

		if err != nil || !ok || cnt > PeerIPQPSLimit {
			response.Error(c, errors.New("请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}
