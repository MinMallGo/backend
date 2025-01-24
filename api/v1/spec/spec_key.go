package Spec

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func CreateSpecKey(c *gin.Context) {
	create := &dto.SpecKeyCreate{}
	if err := c.ShouldBindJSON(create); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeyCreate(c, create)
}

func DeleteSpecKey(c *gin.Context) {
	del := &dto.SpecKeyDelete{}
	if err := c.ShouldBindJSON(del); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeyDelete(c, del)
}

func UpdateSpecKey(c *gin.Context) {
	update := &dto.SpecKeyUpdate{}
	if err := c.ShouldBindJSON(update); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeyUpdate(c, update)
}

func SearchSpecKey(c *gin.Context) {
	search := &dto.SpecKeySearch{}
	if err := c.ShouldBindJSON(search); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeySearch(c, search)
}
