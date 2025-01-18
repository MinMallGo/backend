package spu

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

// CRUD spu_category

func CreateCategory(c *gin.Context) {
	create := &dto.CategoryCreate{}
	if err := c.ShouldBindJSON(create); err != nil {
		response.Error(c, err)
		return
	}

	service.CategoryCreate(c, create)
}

func UpdateCategory(c *gin.Context) {
	update := &dto.CategoryUpdate{}
	if err := c.ShouldBindJSON(update); err != nil {
		response.Error(c, err)
		return
	}

	service.CategoryUpdate(c, update)
}

func SearchCategory(c *gin.Context) {
	search := &dto.CategorySearch{}
	if err := c.ShouldBindJSON(search); err != nil {
		response.Error(c, err)
		return
	}

	service.CategorySearch(c, search)
}

func DeleteCategory(c *gin.Context) {
	del := &dto.CategoryDelete{}
	if err := c.ShouldBindJSON(del); err != nil {
		response.Error(c, err)
		return
	}

	service.CategoryDelete(c, del)
}

// CURD 商品

func CreateSpu(c *gin.Context) {
	create := &dto.SpuCreate{}
	if err := c.ShouldBindJSON(create); err != nil {
		response.Error(c, err)
		return
	}
	service.SpuCreate(c, create)
}

func UpdateSpu(c *gin.Context) {
	update := &dto.SpuUpdate{}
	if err := c.ShouldBindJSON(update); err != nil {
		response.Error(c, err)
		return
	}
	service.SpuUpdate(c, update)
}

func SearchSpu(c *gin.Context) {
	search := &dto.SpuSearch{}
	if err := c.ShouldBindJSON(search); err != nil {
		response.Error(c, err)
		return
	}
	service.SpuSearch(c, search)
}

func DeleteSpu(c *gin.Context) {
	del := &dto.SpuDelete{}
	if err := c.ShouldBindJSON(del); err != nil {
		response.Error(c, err)
		return
	}
	service.SpuDelete(c, del)
}
