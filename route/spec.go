package routes

import (
	"github.com/gin-gonic/gin"
	spec "mall_backend/api/v1/spec"
	"mall_backend/middleware"
)

func SpecRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		// spec_value的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/spec_value/create", spec.CreateSpecValue)
		v1.POST("/spec_value/delete", spec.DeleteSpecValue)
		v1.POST("/spec_value/update", spec.UpdateSpecValue)
		v1.POST("/spec_value/search", spec.SearchSpecValue)
		//spec_value的操作

		// spec_key的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/spec_key/create", spec.CreateSpecKey)
		v1.POST("/spec_key/delete", spec.DeleteSpecKey)
		v1.POST("/spec_key/update", spec.UpdateSpecKey)
		v1.POST("/spec_key/search", spec.SearchSpecKey)
		//spec_key的操作
	}
}
