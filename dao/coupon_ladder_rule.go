package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

type CouponLadderRuleDao struct {
	db *gorm.DB
}

func NewCouponLadderRuleDao(db *gorm.DB) *CouponLadderRuleDao {
	return &CouponLadderRuleDao{db: db}
}

func (d *CouponLadderRuleDao) Create(couponID int, create *dto.CouponCreate) (err error) {
	param := make([]model.MmCouponLadderRule, 0, len(create.LadderRule))
	for _, ladder := range create.LadderRule {
		param = append(param, model.MmCouponLadderRule{
			CouponID:       int32(couponID),
			ThresholdType:  int32(create.ThresholdType),
			ThresholdValue: int32(ladder.ThresholdValue),
			DiscountType:   int32(create.DiscountType),
			DiscountValue:  int32(ladder.DiscountValue),
			Status:         true,
			CreateTime:     time.Now(),
			UpdateTime:     util.MinDateTime(),
		})
	}

	tx := d.db.Model(&model.MmCouponLadderRule{}).Create(param)
	if tx.Error != nil || tx.RowsAffected != int64(len(param)) {
		return tx.Error
	}
	return nil
}

func (d *CouponLadderRuleDao) Delete(couponID int) error {
	return d.db.Model(&model.MmCouponLadderRule{}).Where("coupon_id = ?", couponID).Update("status", false).Error
}
