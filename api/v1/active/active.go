package active

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service/active"
)

func Create(c *gin.Context) {
	param := &dto.ActiveCreate{}
	if err := c.ShouldBind(param); err != nil {
		response.Failure(c, err.Error())
		return
	}
	active.Create(c, param)
}

func Update(c *gin.Context) {
	param := &dto.ActiveUpdate{}
	if err := c.ShouldBind(param); err != nil {
		response.Failure(c, err.Error())
		return
	}
	active.Update(c, param)
}
