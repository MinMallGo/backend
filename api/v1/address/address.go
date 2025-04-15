package address

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func Create(c *gin.Context) {
	param := &dto.AddressCreate{}
	if err := c.ShouldBind(&param); err != nil {
		response.Error(c, err)
		return
	}
	service.AddressCreate(c, param)
}
