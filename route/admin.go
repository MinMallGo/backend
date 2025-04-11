package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/admin"
	"mall_backend/middleware"
)

func AdminRegister(r *gin.Engine) {
	oauth := r.Group("v1").Use(middleware.AdminLogin())
	{
		oauth.POST("/admin/login", admin.Login)
		oauth.POST("/admin/register", admin.Register)
	}

	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/admin/logout", admin.Logout)
	}
}
