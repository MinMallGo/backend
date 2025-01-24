package service

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
)

func SpecValueCreate(c *gin.Context, create *dto.SpecValueCreate) {
	err := dao.SpecValueCreate(create)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, []string{})
	return
}

func SpecValueDelete(c *gin.Context, delete *dto.SpecValueDelete) {
	err := dao.SpecValueDelete(delete)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, []string{})
	return
}

func SpecValueUpdate(c *gin.Context, update *dto.SpecValueUpdate) {
	err := dao.SpecValueUpdate(update)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, []string{})
	return
}

func SpecValueSearch(c *gin.Context, search *dto.SpecValueSearch) {
	res, err := dao.SpecValueSearch(search)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, res)
	return
}
