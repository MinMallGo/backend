package dto

import (
	"mall_backend/dao/model"
)

type SpuCreate struct {
	Name       string   `json:"name" binding:"required,min=1,max=32"`
	Desc       string   `json:"desc" binding:"omitempty,min=1,max=255"`
	Img        []string `json:"img" binding:"omitempty,dive"` // TODO 图片上传
	BrandId    int      `json:"brand_id" binding:"omitempty,gte=0"`
	CategoryId int      `json:"category_id" binding:"required,gt=0"`
	MerchantId int      `json:"merchant_id" binding:"omitempty,gte=0"`
}

type SpuDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}

type SpuUpdate struct {
	SpuDelete
	Name       string   `json:"name" binding:"omitempty,min=1,max=32"`
	Desc       string   `json:"desc" binding:"omitempty,min=1,max=255"`
	Img        []string `json:"img" binding:"omitempty,dive"` // TODO 图片上传
	BrandId    int      `json:"brand_id" binding:"omitempty,gt=0"`
	CategoryId int      `json:"category_id" binding:"omitempty,gt=0"`
	MerchantId int      `json:"merchant_id" binding:"omitempty,gt=0"`
}

type SpuSearch struct {
	Id         int    `json:"id" binding:"omitempty,gt=0"`
	Name       string `json:"name" binding:"omitempty,min=1,max=32"`
	BrandId    int    `json:"brand_id" binding:"omitempty,gt=0"`
	CategoryId int    `json:"category_id" binding:"omitempty,gt=0"`
	MerchantId int    `json:"merchant_id" binding:"omitempty,gt=0"`
	Page       int    `json:"page" binding:"gt=0"`
	Limit      int    `json:"limit"`
}

// SpuSearchResponse 商品查询返回数据定义
type SpuSearchResponse struct {
	model.MmSpu
	Category *model.MmCategory `json:"category"`
}
