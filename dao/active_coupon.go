package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/util"
	"time"
)

type ActiveCouponDao struct {
	db *gorm.DB
}

func NewActiveCouponDao(db *gorm.DB) *ActiveCouponDao {
	return &ActiveCouponDao{
		db: db,
	}
}

func (d *ActiveCouponDao) Create(activeID int, coupons ...int) error {
	if len(coupons) == 0 {
		return nil
	}
	
	param := make([]model.MmActiveCoupon, 0, len(coupons))
	for _, coupon := range coupons {
		param = append(param, model.MmActiveCoupon{
			ActiveID:   int32(activeID),
			CouponID:   int32(coupon),
			Status:     true,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		})
	}
	tx := d.db.Model(&model.MmActiveCoupon{}).Create(param)
	if tx.Error != nil || tx.RowsAffected != int64(len(coupons)) {
		return errors.New("创建活动优惠券失败")
	}
	return nil
}
