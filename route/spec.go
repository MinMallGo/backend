package routes

import (
	"github.com/gin-gonic/gin"
	spec "mall_backend/api/v1/spec"
	"mall_backend/middleware"
)

func SpecRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		v1.POST("/spec/create", spec.CreateSpec)
		v1.POST("/spec/delete", spec.DeleteSpec)
		v1.POST("/spec/update", spec.UpdateSpec)
		v1.POST("/spec/search", spec.SearchSpec)
	}
}
