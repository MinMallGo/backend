package dto

type OrderCreate struct {
	// 来自购物车
	//CartID    string         `json:"cart_id" binding:"required"`
	AddressID int            `json:"address_id" binding:"required"` // TODO 这里下单的时候不能选多个地址，而是在下单完成之后，可以修改成多个地址
	Coupons   []int          `json:"coupons" binding:"required,omitempty"`
	Product   []ShoppingItem `json:"product" binding:"required,dive"`
	PayWay    string         `json:"pay_way" binding:"required,oneof=alipay wechatpay"`
	Source    string         `json:"source" binding:"required,oneof=H5 min-program"`
}

type ShoppingItem struct {
	SpuID     int `json:"spu_id" binding:"required"`
	SkuID     int `json:"sku_id" binding:"required"`
	Num       int `json:"num" binding:"required"`
	IsSecKill int `json:"is_sec_kill" binding:"required,oneof=1 2"` // 如果是秒杀则不允许合并创建订单
}

// 创建订单 / 下单之后，仅能修改收货地址

type UpdateOrder struct {
	OrderID   string `json:"order_id" binding:"required"`
	AddressId string `json:"address_id" binding:"required"`
}
