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
	ID               int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name             string    `gorm:"column:name;not null;comment:优惠券名称" json:"name"`                               // 优惠券名称
	DiscountType     int32     `gorm:"column:discount_type;comment:1 = 满减，2 = 折扣" json:"discount_type"`              // 1 = 满减，2 = 折扣
	UseRange         int32     `gorm:"column:use_range;comment:使用范围：1 全平台，2商家 3 类别 4 具体商品 5 某个规格" json:"use_range"`  // 使用范围：1 全平台，2商家 3 类别 4 具体商品 5 某个规格
	Threshold        int32     `gorm:"column:threshold;comment:最低消费多少才能使用" json:"threshold"`                         // 最低消费多少才能使用
	DiscountPrice    int32     `gorm:"column:discount_price;comment:优惠券金额 * 100  或者 百分之多少的折扣" json:"discount_price"` // 优惠券金额 * 100  或者 百分之多少的折扣
	DiscountRate     int32     `gorm:"column:discount_rate;comment:打折比例，1-100 .这是里放大了100倍的" json:"discount_rate"`    // 打折比例，1-100 .这是里放大了100倍的
	Code             string    `gorm:"column:code;comment:优惠券唯一表示，可用于核销等" json:"code"`                               // 优惠券唯一表示，可用于核销等
	ReceiveCount     int32     `gorm:"column:receive_count;comment:可以领取的数量" json:"receive_count"`                    // 可以领取的数量
	ReceiveStartTime time.Time `gorm:"column:receive_start_time;comment:开始领取时间" json:"receive_start_time"`           // 开始领取时间
	ReceiveEndTime   time.Time `gorm:"column:receive_end_time;comment:领取结束时间" json:"receive_end_time"`               // 领取结束时间
	UseStartTime     time.Time `gorm:"column:use_start_time;comment:优惠券开始时间" json:"use_start_time"`                  // 优惠券开始时间
	UseEndTime       time.Time `gorm:"column:use_end_time;comment:优惠券结束时间" json:"use_end_time"`                      // 优惠券结束时间
	TotalCount       int32     `gorm:"column:total_count;comment:优惠券总数" json:"total_count"`                          // 优惠券总数
	AssginCount      int32     `gorm:"column:assgin_count;comment:已经领取的数量" json:"assgin_count"`                      // 已经领取的数量
	UseCount         int32     `gorm:"column:use_count;comment:已经使用的优惠券数量" json:"use_count"`                         // 已经使用的优惠券数量
	Status           bool      `gorm:"column:status;default:1" json:"status"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime       time.Time `gorm:"column:delete_time" json:"delete_time"`
}

// TableName MmCoupon's table name
func (*MmCoupon) TableName() string {
	return TableNameMmCoupon
}
