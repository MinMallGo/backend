package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/active/seckill"
	"mall_backend/middleware"
)

func ActiveRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		// sec_kill
		v1.POST("/sec_kill/create", seckill.CreateSecKill)
		// sec_kill
	}
}
