package dto

type CartSearch struct {
	Paginate
}

type CartCreate struct {
	SpuId int `json:"spu_id" binding:"required"`
	SkuId int `json:"sku_id" binding:"required"`
	Num   int `json:"num" binding:"required"`
}

// CartUpdate 只能修改sku_id 和 number
type CartUpdate struct {
	Id    int `json:"id" binding:"required"`
	SpuId int `json:"spu_id" binding:"required"`
	SkuId int `json:"sku_id" binding:"required"`
	Num   int `json:"num" binding:"required"`
}

type CartDelete struct {
	Id []int `json:"id" binding:"required"`
}
