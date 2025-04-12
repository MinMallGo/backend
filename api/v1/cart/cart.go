package cart

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func Search(c *gin.Context) {
	// 好像没有需要验证的。直接就嗯
	param := &dto.CartSearch{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.CartSearch(c, param)
}

func Create(c *gin.Context) {
	param := &dto.CartCreate{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.CartCreate(c, param)
}

func Update(c *gin.Context) {
	param := &dto.CartUpdate{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.CartUpdate(c, param)
}

func Delete(c *gin.Context) {
	param := &dto.CartDelete{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.CartDelete(c, param)
}
