package dto

type OrderCreate struct {
	OrderType int            `json:"order_type"  binding:"required,oneof=1 2"` // 1 普通订单 秒杀订单
	AddressID int            `json:"address_id" binding:"required"`            // TODO 这里下单的时候不能选多个地址，而是在下单完成之后，可以修改成多个地址
	Coupons   []int          `json:"coupons" binding:"required,omitempty"`
	Product   []ShoppingItem `json:"product" binding:"required,dive"`
	Source    int            `json:"source" binding:"required,oneof=1 2"`
}

type ShoppingItem struct {
	SpuID int `json:"spu_id" binding:"required"`
	SkuID int `json:"sku_id" binding:"required"`
	Num   int `json:"num" binding:"required"`
	//MerchantID int `json:"merchant_id" binding:"omitempty"`
}

// 创建订单 / 下单之后，仅能修改收货地址

type UpdateOrder struct {
	OrderID   string `json:"order_id" binding:"required"`
	AddressId string `json:"address_id" binding:"required"`
}

type PayOrder struct {
	BatchCode string `json:"batch_code" form:"batch_code" binding:"required"`
	OrderCode string `json:"order_code" form:"order_code" binding:"required"`
	PayWay    int    `json:"pay_way" form:"pay_way" binding:"required,oneof=1 2"`
}

type CancelOrder struct {
	ID        int    `json:"id" binding:"required"`
	OrderCode string `json:"order_code" binding:"required"`
}

type OrderExpire struct {
	ID        int    `json:"id" binding:"required"`
	OrderCode string `json:"order_code" binding:"required"`
}
