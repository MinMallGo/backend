package menu

import (
	"github.com/gin-gonic/gin"
	"mall_backend/service"
)

func Search(c *gin.Context) {
	service.Menu()
}
