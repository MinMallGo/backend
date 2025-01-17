package routes

import (
	"github.com/gin-gonic/gin"
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
		// spu分类的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/category/create", spu.CreateCategory)
		v1.POST("/category/delete", spu.DeleteCategory)
		v1.POST("/category/update", spu.UpdateCategory)
		v1.POST("/category/search", spu.SearchCategory)
		// spu分类的操作
	}
}
