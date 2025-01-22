package dto

type SpecCreate struct {
	SpuId   int `json:"spu_id" binding:"required"`
	SkuId   int `json:"sku_id" binding:"required"`
	KeyId   int `json:"key_id" binding:"required"`
	ValueId int `json:"value_id" binding:"required"`
}

type SpecSearch struct {
	Id    int `json:"id" binding:"omitempty gt=0"`
	SkuId int `json:"sku_id" binding:"omitempty gt=0"`
	SpuId int `json:"spu_id" binding:"omitempty gt=0"`
	Paginate
}

type SpecDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}

type SpecUpdate struct {
	Id      int `json:"id" binding:"required"`
	SkuId   int `json:"sku_id" binding:"omitempty gt=0"`
	SpuId   int `json:"spu_id" binding:"omitempty gt=0"`
	KeyId   int `json:"key_id" binding:"omitempty  gt=0"`
	ValueId int `json:"value_id" binding:"omitempty  gt=0"`
}
