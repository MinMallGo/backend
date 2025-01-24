package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/sku"
	"mall_backend/middleware"
)

func SkuRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		// sku的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/sku/create", sku.CreateSku)
		v1.POST("/sku/delete", sku.DeleteSku)
		v1.POST("/sku/update", sku.UpdateSku)
		v1.POST("/sku/search", sku.SearchSku)
		// sku的操作
	}
}
