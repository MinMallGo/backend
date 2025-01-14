package middleware

import (
	"github.com/gin-gonic/gin"
	"mall_backend/response"
	"mall_backend/util"
)

func AuthVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT验证登录
		auth := c.GetHeader("Authorization")
		decode, err := util.JWTDecode(auth)
		if err != nil || decode.UserID == 0 {
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
