package dao

import (
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
)

func CouponGetByIds(ids []int) (*[]model.MmCoupon, error) {
	res := &[]model.MmCoupon{}
	err := util.DBClient().Model(&model.MmCoupon{}).Where("status = ?", constants.NormalStatus).Find(res, ids)
	if err != nil {
		return nil, err.Error
	}
	return res, nil
}

// CouponCreate TODO 封装
func CouponCreate(create *dto.CouponCreate) error {

	return nil
}
