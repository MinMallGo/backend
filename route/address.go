package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/address"
	"mall_backend/middleware"
)

func AddressRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/addr/create", address.Create)
	}
}
