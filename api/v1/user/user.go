package user

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
	"mall_backend/util"
)

func Login(c *gin.Context) {
	param := &dto.UserLogin{}
	if err := c.ShouldBind(param); err != nil {
		response.Failure(c, util.HandleValidationError(err))
		return
	}

	service.Login(c, param)
}

func Logout(c *gin.Context) {
	service.Logout(c)
}

func Register(c *gin.Context) {
	register := &dto.UserRegister{}
	if err := c.ShouldBind(register); err != nil {
		response.Failure(c, util.HandleValidationError(err))
		return
	}
	service.Register(c, register)
}
