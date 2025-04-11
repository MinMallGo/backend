package Spec

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func CreateSpec(c *gin.Context) {
	create := &dto.SpecKeyCreate{}
	if err := c.ShouldBindJSON(create); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeyCreate(c, create)
}

func DeleteSpec(c *gin.Context) {
	del := &dto.SpecKeyDelete{}
	if err := c.ShouldBindJSON(del); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeyDelete(c, del)
}

func UpdateSpec(c *gin.Context) {
	update := &dto.SpecKeyUpdate{}
	if err := c.ShouldBindJSON(update); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeyUpdate(c, update)
}

func SearchSpec(c *gin.Context) {
	search := &dto.SpecKeySearch{}
	if err := c.ShouldBindJSON(search); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecKeySearch(c, search)
}
