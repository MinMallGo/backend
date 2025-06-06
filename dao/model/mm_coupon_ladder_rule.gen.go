// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmCouponLadderRule = "mm_coupon_ladder_rule"

// MmCouponLadderRule 优惠券多梯度规则表
type MmCouponLadderRule struct {
	ID       int32 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CouponID int32 `gorm:"column:coupon_id;not null;comment:优惠券id" json:"coupon_id"` // 优惠券id
	/*
		1. 金额门槛
		2. 件数门槛
	*/
	ThresholdType  int32     `gorm:"column:threshold_type;comment:1. 金额门槛\n2. 件数门槛" json:"threshold_type"`
	ThresholdValue int32     `gorm:"column:threshold_value;comment:金额门槛或者数额门槛" json:"threshold_value"` // 金额门槛或者数额门槛
	DiscountType   int32     `gorm:"column:discount_type;comment:1. 满减 2. 打折" json:"discount_type"`    // 1. 满减 2. 打折
	DiscountValue  int32     `gorm:"column:discount_value;comment:折扣比例：9折 = 90" json:"discount_value"` // 折扣比例：9折 = 90
	Status         bool      `gorm:"column:status;default:1" json:"status"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime     time.Time `gorm:"column:update_time" json:"update_time"`
}

// TableName MmCouponLadderRule's table name
func (*MmCouponLadderRule) TableName() string {
	return TableNameMmCouponLadderRule
}
