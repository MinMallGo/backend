package dto

type SpecCreate struct {
	SkuID int    `json:"sku_id" binding:"required,gt=0"`
	Name  string `json:"name" binding:"required,min=1,max=32"`
	ReasonableValue
}

type SpecDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}

type SpecUpdate struct {
	SpecDelete
	SkuID int    `json:"sku_id" binding:"omitempty,gt=0"`
	Name  string `json:"name" binding:"omitempty,min=1,max=32"`
	ReasonableValue
}

type SpecSearch struct {
	Id    int    `json:"id" binding:"omitempty,gt=0"`
	Name  string `json:"name" binding:"omitempty,min=1,max=32"`
	SkuID int    `json:"sku_id" binding:"omitempty,gt=0"`
	Page  int    `json:"page" binding:"gt=0"`
	Limit int    `json:"limit"`
}

type ReasonableValue struct {
	Price int `json:"price" binding:"omitempty,gte=0"`
	Stock int `json:"stock" binding:"omitempty,gt=0"`
}
