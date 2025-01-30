package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/user"
	"mall_backend/middleware"
)

func AdminRegister(r *gin.Engine) {
	oauth := r.Group("v1").Use(middleware.AdminLogin())
	{
		oauth.POST("/admin/login", user.Login)
		oauth.POST("/admin/register", user.Register)
	}

	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/admin/logout", user.Logout)
	}
}
