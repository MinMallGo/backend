package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
)

type OrderCouponDao struct {
	db *gorm.DB
}

func NewOrderCouponDao(db *gorm.DB) *OrderCouponDao {
	return &OrderCouponDao{
		db: db,
	}
}

func (d *OrderCouponDao) Create(coupon *[]model.MmOrderCoupon) error {
	if len(*coupon) < 1 {
		return nil
	}
	if tx := d.db.Model(&model.MmOrderCoupon{}).Create(coupon); tx.Error != nil {
		return tx.Error
	}
	return nil
}
