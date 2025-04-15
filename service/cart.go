package service

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
)

func CartSearch(c *gin.Context, param *dto.CartSearch) {
	//	 TODO 获取当前登录的用户信息
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	data, err := dao.NewCartDao(util.DBClient()).Search(param, int(user.ID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.PaginateSuccess(c, data)
	return
}

func CartCreate(c *gin.Context, param *dto.CartCreate) {
	//	 TODO 获取当前登录的用户信息
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	err = dao.NewCartDao(util.DBClient()).Create(param, int(user.ID))

	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, []struct{}{})
	return
}

func CartUpdate(c *gin.Context, param *dto.CartUpdate) {
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	err = dao.NewCartDao(util.DBClient()).Update(param, int(user.ID))

	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, []struct{}{})
	return
}

func CartDelete(c *gin.Context, param *dto.CartDelete) {
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	err = dao.NewCartDao(util.DBClient()).Delete(param, int(user.ID))

	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, []struct{}{})
	return
}
