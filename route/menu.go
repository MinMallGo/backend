package routes

import (
	"github.com/gin-gonic/gin"
	"mall_backend/api/v1/menu"
)

func MenuInit(r *gin.Engine) {
	oauth := r.Group("v1")
	{
		oauth.GET("/menu/list", menu.Search)
	}

}
