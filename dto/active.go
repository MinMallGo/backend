package dto

import "time"

type ActiveCreate struct {
	Name      string          `json:"name" binding:"required,min=1,max=32"`
	Type      int             `json:"type" binding:"required,oneof=1 2"`
	Desc      string          `json:"desc" binding:"omitempty,max=255"`
	StartTime time.Time       `json:"start_time" binding:"required"`
	EndTime   time.Time       `json:"end_time" binding:"required,gtfield=StartTime"`
	Coupons   []int           `json:"coupons" binding:"omitempty,dive"` // 本场活动所包含的优惠券
	SecKills  []SecKillCreate `json:"sec_kills" binding:"omitempty,dive"`
	// ? 其他的再扩展吧
}

type SecKillCreate struct {
	SpuID     int       `json:"spu_id" binding:"required,gt=0"`
	SkuID     int       `json:"sku_id" binding:"required,gt=0"`
	Price     int       `json:"price" binding:"required,gt=0"`
	Stock     int       `json:"stock" binding:"required,gt=0"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
}

type ActiveUpdate struct {
	ID int `json:"id" binding:"required,gt=0"`
	ActiveCreate
}
