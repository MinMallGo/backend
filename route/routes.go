package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/user"
)

func Register(r *gin.Engine) {
	v1 := r.Group("v1")
	{
		v1.POST("/user/login", user.Login)
	}
}
