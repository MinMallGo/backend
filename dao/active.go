package dao

import (
	"errors"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

func ActiveGetByIds(ids []int) (*model.MmCoupon, error) {
	res := &model.MmCoupon{}
	err := util.DBClient().Model(&model.MmCoupon{}).Where("status = ?", constants.NormalStatus).Find(res, ids)
	if err != nil {
		return nil, err.Error
	}
	return res, nil
}

// ActiveCreate 创建活动
func ActiveCreate(create *dto.ActiveCreate) (int, error) {
	startTime, err := time.Parse("2006-01-02 15:04:05", create.ActiveStartTime)
	if err != nil {
		return 0, err
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", create.ActiveEndTime)
	if err != nil {
		return 0, err
	}

	if len(create.Coupons) > 0 {
		// 先检查这些优惠券是不是全都在
		coupon, err := CouponGetByIds(create.Coupons)
		if err != nil {
			return 0, err
		}
		if coupon == nil || len(*coupon) < len(create.Coupons) {
			return 0, errors.New("请选择正确的优惠券")
		}
	}

	param := &model.MmActive{
		Name:       create.Name,
		Type:       int32(create.Type),
		Desc:       create.Desc,
		StartTime:  startTime,
		EndTime:    endTime,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	res := util.DBClient().Model(&model.MmActive{}).Create(param)
	if res.Error != nil {
		return 0, res.Error
	}

	return int(param.ID), nil
}
