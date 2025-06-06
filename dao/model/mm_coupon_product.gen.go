// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmCouponProduct = "mm_coupon_product"

// MmCouponProduct 优惠券使用范围字典表
type MmCouponProduct struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true;comment:使用范围ID" json:"id"` // 使用范围ID
	CouponID   int32     `gorm:"column:coupon_id;comment:优惠券id" json:"coupon_id"`                  // 优惠券id
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	CategoryID int32     `gorm:"column:category_id" json:"category_id"`
	SpuID      int32     `gorm:"column:spu_id;comment:适用商品" json:"spu_id"` // 适用商品
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName MmCouponProduct's table name
func (*MmCouponProduct) TableName() string {
	return TableNameMmCouponProduct
}
