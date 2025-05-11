package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"time"
)

const (
	UnUsed = iota
	Used
	Expired
	Frozen
)

type UserCouponDao struct {
	db *gorm.DB
}

func NewUserCouponDao(db *gorm.DB) *UserCouponDao {
	return &UserCouponDao{
		db: db,
	}
}

func (d *UserCouponDao) Use(userId int, ids ...int) error {
	if len(ids) == 0 {
		return nil
	}
	
	tx := d.db.Model(&model.MmUserCoupon{}).
		Where("user_id = ?", userId).
		Where("coupon_id in ?", ids).
		Updates(map[string]interface{}{
			"is_used": 1,
		})
	if tx.Error != nil || tx.RowsAffected != int64(len(ids)) {
		return errors.New("用户优惠券扣减失败")
	}
	return nil
}

// 没问题，一张券只能使用一次。所以不存在使用了多次之后，取消券改到之前的

func (d *UserCouponDao) Cancel(userId int, couponIds ...int) error {
	tx := d.db.Model(&model.MmUserCoupon{}).
		Where("user_id = ? AND coupon_id in ?", userId, couponIds).
		Updates(map[string]interface{}{
			"is_used":     0,
			"update_time": time.Now(),
		})
	if tx.Error != nil || tx.RowsAffected < int64(len(couponIds)) {
		return errors.New("用户优惠券归还失败")
	}
	return nil
}

func (d *UserCouponDao) Exists(userId int, ids ...int) error {
	var count int64
	res := d.db.Model(&model.MmUserCoupon{}).
		Where("id in ?", ids).
		Where("user_id = ?", userId).
		Where("is_used = ?", false).
		Where("status = ?", 1).
		Count(&count)
	if res.Error != nil || (count) < int64(len(ids)) {
		return res.Error
	}
	return nil
}
