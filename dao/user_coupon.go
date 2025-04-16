package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
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
