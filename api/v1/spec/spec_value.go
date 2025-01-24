package Spec

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func CreateSpecValue(c *gin.Context) {
	create := &dto.SpecValueCreate{}
	if err := c.ShouldBindJSON(create); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecValueCreate(c, create)
}

func DeleteSpecValue(c *gin.Context) {
	del := &dto.SpecValueDelete{}
	if err := c.ShouldBindJSON(del); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecValueDelete(c, del)
}

func UpdateSpecValue(c *gin.Context) {
	update := &dto.SpecValueUpdate{}
	if err := c.ShouldBindJSON(update); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecValueUpdate(c, update)
}

func SearchSpecValue(c *gin.Context) {
	search := &dto.SpecValueSearch{}
	if err := c.ShouldBindJSON(search); err != nil {
		response.Error(c, err)
		return
	}

	service.SpecValueSearch(c, search)
}
