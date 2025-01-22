package dao

import (
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/util"
	"time"
)

// ActiveCouponCreate 创建活动时使用优惠券
func ActiveCouponCreate(activeId int, coupon *[]int) error {
	param := &[]model.MmActiveCoupon{}
	for _, v := range *coupon {
		*param = append(*param, model.MmActiveCoupon{
			ActiveID:   int32(activeId),
			CouponID:   int32(v),
			Status:     constants.NormalStatus,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		})
	}

	return util.DBClient().Model(&model.MmActiveCoupon{}).Create(param).Error
}
