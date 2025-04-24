package coupon

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
)

/*
优惠券的 crud
*/

// Create 创建Coupon分类
func Create(c *gin.Context, create *dto.CouponCreate) {
	// 1 全平台，2商家 3 类别 4 具体商品 5 某个规格
	// TODO 如果creat.UseRange 不是 1 检查具体的东西
	// todo 这里需要检查，如果是多梯度，则ladder_rule不能为空
	// todo 这里需要用事务，因为还要写多梯度表

	// 1. 检查usage_scope是否存在
	if !dao.NewCouponUsageScopeDao(util.DBClient()).Exists(create.ScopeID) {
		response.Failure(c, "创建优惠券失败：请检查使用范围")
		return
	}

	// 还有当创建类型是折扣是，注意不能超过100或者小于1
	if create.DiscountType == dao.RateType && (create.DiscountValue > 100 || create.DiscountValue < 0) {
		response.Failure(c, "创建优惠券失败：请选择正确的折扣范围")
		return
	}

	if create.IsLadder == 1 && len(create.LadderRule) < 1 {
		response.Failure(c, "创建优惠券失败：多梯度优惠券创建失败")
		return
	}

	// todo 还需要检查，多梯度优惠券里面的，如果是打折内容，打折的比例也应该保持在 1-99范围内
	if create.IsLadder == 1 && create.DiscountType == dao.RateType {
		for _, ladder := range create.LadderRule {
			if ladder.DiscountValue > 100 || ladder.DiscountValue < 0 {
				response.Failure(c, "创建优惠券失败：请选择正确的多梯度优惠券折扣")
				return
			}
		}
	}

	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		couponID, err := dao.NewCouponDao(tx).Create(create)
		if err != nil {
			return errors.Join(errors.New("创建优惠券失败："), err)
		}
		err = dao.NewCouponLadderRuleDao(tx).Create(couponID, create)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, []string{})
	return
}

// Delete 删除Coupon分类
func Delete(c *gin.Context, delete *dto.CouponDelete) {
	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		err := dao.NewCouponDao(tx).Delete(delete)
		if err != nil {
			return errors.Join(errors.New("删除优惠券失败："), err)
		}
		err = dao.NewCouponLadderRuleDao(tx).Delete(delete.Id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, []string{})
	return
}

// Update 修改Coupon的分类信息
func Update(c *gin.Context, update *dto.CouponUpdate) {
	// 1. 检查usage_scope是否存在
	if !dao.NewCouponUsageScopeDao(util.DBClient()).Exists(update.ScopeID) {
		response.Failure(c, "创建优惠券失败：请检查使用范围")
		return
	}

	// 还有当创建类型是折扣是，注意不能超过100或者小于1
	if update.DiscountType == dao.RateType && (update.DiscountValue > 100 || update.DiscountValue < 0) {
		response.Failure(c, "创建优惠券失败：请选择正确的折扣范围")
		return
	}

	if update.IsLadder == 1 && len(update.LadderRule) < 1 {
		response.Failure(c, "创建优惠券失败：多梯度优惠券创建失败")
		return
	}

	// todo 还需要检查，多梯度优惠券里面的，如果是打折内容，打折的比例也应该保持在 1-99范围内
	if update.IsLadder == 1 && update.DiscountType == dao.RateType {
		for _, ladder := range update.LadderRule {
			if ladder.DiscountValue > 100 || ladder.DiscountValue < 0 {
				response.Failure(c, "创建优惠券失败：请选择正确的多梯度优惠券折扣")
				return
			}
		}
	}

	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		err := dao.NewCouponDao(tx).Update(update)
		if err != nil {
			return errors.Join(errors.New("更新优惠券失败"), err)
		}
		err = dao.NewCouponLadderRuleDao(tx).Delete(update.ID)
		if err != nil {
			return err
		}
		err = dao.NewCouponLadderRuleDao(tx).Create(update.ID, &update.CouponCreate)
		return nil
	})

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, []string{})
	return
}

// Search 查询
func Search(c *gin.Context, search *dto.CouponSearch) {
	// id, name, page, limit
	res, err := dao.NewCouponDao(util.DBClient()).Search(search)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.PaginateSuccess(c, &res)
	return
}
