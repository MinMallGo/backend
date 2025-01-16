package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"mall_backend/response"
	"mall_backend/util"
)

func AuthVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		login, err := util.CacheClient().Get(context.Background(), token).Result()
		if len(login) < 2 || err != nil {
			// 返回请登录的提示
			response.Failure(c, "请登录后操作")
			c.Abort() // 终止请求
			return
		} else {
			// 调用该请求的剩余处理程序
			c.Next()
		}
	}
}
