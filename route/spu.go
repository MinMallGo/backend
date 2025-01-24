package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/spu"
	"mall_backend/middleware"
)

func SpuRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		// spu的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/spu/create", spu.CreateSpu)
		v1.POST("/spu/delete", spu.DeleteSpu)
		v1.POST("/spu/update", spu.UpdateSpu)
		v1.POST("/spu/search", spu.SearchSpu)
		// spu的操作
	}
}
