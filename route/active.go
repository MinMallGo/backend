package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/active"
	"mall_backend/middleware"
)

func ActiveRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/active/create", active.Create)
		v1.POST("/active/update", active.Update)
	}
}
