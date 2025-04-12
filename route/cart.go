package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/cart"
	"mall_backend/middleware"
)

func CartRegister(r *gin.Engine) {
	v := r.Group("v1").Use(middleware.AuthVerify())
	{
		v.POST("/cart/search", cart.Search)
		v.POST("/cart/create", cart.Create)
		v.POST("/cart/update", cart.Update)
		v.POST("/cart/delete", cart.Delete)
	}
}
