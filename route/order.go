package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/order"
	"mall_backend/middleware"
)

func OrderRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/order/create", order.Create)
		v1.POST("/order/pay", order.Pay)
		v1.POST("/order/cancel", order.CancelOrder)
		v1.POST("/order/expire", order.Expire)
	}
	r.Group("v1").GET("/order/pay", order.Pay2)
}
