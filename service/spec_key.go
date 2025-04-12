package service

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
)

func SpecKeyCreate(c *gin.Context, create *dto.SpecKeyCreate) {
	res, err := dao.NewSpecKeyDao().Create(create)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, res)
	return
}

func SpecKeyDelete(c *gin.Context, delete *dto.SpecKeyDelete) {
	err := dao.NewSpecKeyDao().Delete(delete)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, []string{})
	return
}

func SpecKeyUpdate(c *gin.Context, update *dto.SpecKeyUpdate) {
	err := dao.NewSpecKeyDao().Update(update)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, []string{})
	return
}

func SpecKeySearch(c *gin.Context, search *dto.SpecKeySearch) {
	res, err := dao.NewSpecKeyDao().More(search)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.PaginateSuccess(c, res)
	return
}
