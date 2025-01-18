package sku

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func CreateSku(c *gin.Context) {
	create := &dto.SkuCreate{}
	if err := c.ShouldBindJSON(create); err != nil {
		response.Error(c, err)
		return
	}
	service.SkuCreate(c, create)
}

func UpdateSku(c *gin.Context) {
	update := &dto.SkuUpdate{}
	if err := c.ShouldBindJSON(update); err != nil {
		response.Error(c, err)
		return
	}
	service.SkuUpdate(c, update)
}

func SearchSku(c *gin.Context) {
	search := &dto.SkuSearch{}
	if err := c.ShouldBindJSON(search); err != nil {
		response.Error(c, err)
		return
	}
	service.SkuSearch(c, search)
}

func DeleteSku(c *gin.Context) {
	del := &dto.SkuDelete{}
	if err := c.ShouldBindJSON(del); err != nil {
		response.Error(c, err)
		return
	}
	service.SkuDelete(c, del)
}
