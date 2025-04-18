// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmSeckillProduct = "mm_seckill_product"

// MmSeckillProduct 拆分秒杀活动商品
type MmSeckillProduct struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SeckillID  int64     `gorm:"column:seckill_id;not null;comment:秒杀id" json:"seckill_id"` // 秒杀id
	SpuID      int64     `gorm:"column:spu_id;not null;comment:商品id" json:"spu_id"`         // 商品id
	SkuID      int64     `gorm:"column:sku_id;not null;comment:商品规格id" json:"sku_id"`       // 商品规格id
	Price      int64     `gorm:"column:price;not null;comment:价格" json:"price"`             // 价格
	Stock      int32     `gorm:"column:stock;not null;comment:库存" json:"stock"`             // 库存
	Sold       int32     `gorm:"column:sold;comment:卖出" json:"sold"`                        // 卖出
	Version    int32     `gorm:"column:version;comment:算了，先使用悲观锁" json:"version"`           // 算了，先使用悲观锁
	Status     bool      `gorm:"column:status;default:1" json:"status"`
	StartTime  time.Time `gorm:"column:start_time;comment:开始时间" json:"start_time"` // 开始时间
	EndTime    time.Time `gorm:"column:end_time;comment:结束时间" json:"end_time"`     // 结束时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// TableName MmSeckillProduct's table name
func (*MmSeckillProduct) TableName() string {
	return TableNameMmSeckillProduct
}
