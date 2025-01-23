package dto

/**
规格集合。包括 spec_id,spec_key,spec_value
*/

// Specs 创建Sku时的需要多个这种组合
type Specs struct {
	KeyID int `json:"key_id" binding:"required,gt=0"`
	ValID int `json:"value_id" binding:"required,gt=0"`
}

type ReasonableValue struct {
	Price int `json:"price" binding:"omitempty,gt=0"`
	Stock int `json:"stock" binding:"omitempty,gt=0"`
}

type SkuCreate struct {
	SpuID int     `json:"spu_id" binding:"required,gt=0"`
	Spec  []Specs `json:"spec" binding:"required"`
	ReasonableValue
}

type SkuDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}

type SkuUpdate struct {
	SkuDelete
	SpuID int     `json:"spu_id" binding:"required,gt=0"`
	SkuID int     `json:"sku_id" binding:"required,gt=0"`
	Spec  []Specs `json:"spec" binding:"required"`
	ReasonableValue
}

type SkuSearch struct {
	SpuID int `json:"spu_id" binding:"required,gt=0"`
}
