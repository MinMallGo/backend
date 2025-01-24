package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/spu"
	"mall_backend/middleware"
)

func CategoryRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		// spu_category的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/category/create", spu.CreateCategory)
		v1.POST("/category/delete", spu.DeleteCategory)
		v1.POST("/category/update", spu.UpdateCategory)
		v1.POST("/category/search", spu.SearchCategory)
		// spu_category的操作
	}
}
