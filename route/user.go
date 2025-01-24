package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/user"
	"mall_backend/middleware"
)

func UserRegister(r *gin.Engine) {
	oauth := r.Group("v1")
	{
		oauth.POST("/user/login", user.Login)
		oauth.POST("/user/register", user.Register)
	}

	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/user/logout", user.Logout)
	}
}
