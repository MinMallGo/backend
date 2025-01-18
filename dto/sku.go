package dto

import (
	"mall_backend/dao/model"
)

type SkuCreate struct {
	SpuID int    `json:"spu_id" binding:"required,gt=0"`
	Name  string `json:"name" binding:"required,min=1,max=32"`
	Desc  string `json:"desc" binding:"omitempty,min=1,max=255"`
}

type SkuDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}

type SkuUpdate struct {
	SkuDelete
	SpuID int    `json:"spu_id" binding:"omitempty,gt=0"`
	Name  string `json:"name" binding:"omitempty,min=1,max=32"`
	Desc  string `json:"desc" binding:"omitempty,min=1,max=255"`
}

type SkuSearch struct {
	Id    int    `json:"id" binding:"omitempty,gt=0"`
	SpuID int    `json:"spu_id" binding:"omitempty,gt=0"`
	Name  string `json:"name" binding:"omitempty,min=1,max=32"`
	Page  int    `json:"page" binding:"gt=0"`
	Limit int    `json:"limit"`
}

// SkuSearchResponse 商品查询返回数据定义
type SkuSearchResponse struct {
	model.MmSku
	Specs *[]model.MmSpec `json:"specs" gorm:"foreignkey:SkuID"`
}
