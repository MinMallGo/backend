package service

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
)

func AddressCreate(c *gin.Context, create *dto.AddressCreate) {
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}
	err = dao.NewAddrDao(util.DBClient()).Create(create, int(user.ID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, []string{})
	return
}
