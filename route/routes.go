package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/sku"
	spec "mall_backend/api/v1/spec"
	"mall_backend/api/v1/spu"
	"mall_backend/api/v1/user"
	"mall_backend/middleware"
)

func Register(r *gin.Engine) {
	oauth := r.Group("v1")
	{
		oauth.POST("/user/login", user.Login)
		oauth.POST("/user/register", user.Register)
	}

	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/user/logout", user.Logout)
		// spu_category的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/category/create", spu.CreateCategory)
		v1.POST("/category/delete", spu.DeleteCategory)
		v1.POST("/category/update", spu.UpdateCategory)
		v1.POST("/category/search", spu.SearchCategory)
		// spu_category的操作

		// spu的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/spu/create", spu.CreateSpu)
		v1.POST("/spu/delete", spu.DeleteSpu)
		v1.POST("/spu/update", spu.UpdateSpu)
		v1.POST("/spu/search", spu.SearchSpu)
		// spu的操作

		// sku的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/sku/create", sku.CreateSku)
		v1.POST("/sku/delete", sku.DeleteSku)
		v1.POST("/sku/update", sku.UpdateSku)
		v1.POST("/sku/search", sku.SearchSku)
		// sku的操作

		// spec的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/spec/create", spec.CreateSpec)
		v1.POST("/spec/delete", spec.DeleteSpec)
		v1.POST("/spec/update", spec.UpdateSpec)
		v1.POST("/spec/search", spec.SearchSpec)
		// spec的操作
	}
}
