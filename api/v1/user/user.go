package user

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func Login(c *gin.Context) {
	param := &dto.UserLogin{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}

	service.Login(c, param)
}

func Logout(c *gin.Context) {

}

func Register(c *gin.Context) {

}
