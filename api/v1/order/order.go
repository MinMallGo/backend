package order

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func Create(c *gin.Context) {
	param := &dto.OrderCreate{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.OrderCreate(c, param)

}

func Update(c *gin.Context) {}

func Delete(c *gin.Context) {}

func Search(c *gin.Context) {}
