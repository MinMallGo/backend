// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmCoupon = "mm_coupon"

// MmCoupon 优惠券： 包括 金额，名称，编码，开始结束时间以及固定的字段
type MmCoupon struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       string    `gorm:"column:name;not null;comment:优惠券名称" json:"name"` // 优惠券名称
	Type       string    `gorm:"column:type;comment:满减，折扣" json:"type"`          // 满减，折扣
	MinSpend   string    `gorm:"column:min_spend" json:"min_spend"`
	Price      int32     `gorm:"column:price;not null;comment:优惠券金额 * 100  或者 百分之多少的折扣" json:"price"` // 优惠券金额 * 100  或者 百分之多少的折扣
	Code       string    `gorm:"column:code;comment:优惠券唯一表示，可用于核销等" json:"code"`                      // 优惠券唯一表示，可用于核销等
	StartTime  time.Time `gorm:"column:start_time;not null;comment:优惠券开始时间" json:"start_time"`        // 优惠券开始时间
	EndTime    time.Time `gorm:"column:end_time;comment:优惠券结束时间" json:"end_time"`                     // 优惠券结束时间
	Status     bool      `gorm:"column:status;default:1" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
}

// TableName MmCoupon's table name
func (*MmCoupon) TableName() string {
	return TableNameMmCoupon
}
