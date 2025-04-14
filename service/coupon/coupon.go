package coupon

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
)

/*
优惠券的 crud
*/

// Create 创建Coupon分类
func Create(c *gin.Context, create *dto.CouponCreate) {
	// 1 全平台，2商家 3 类别 4 具体商品 5 某个规格
	// TODO 如果creat.UseRange 不是 1 检查具体的东西

	err := dao.NewCouponDao().Create(create)
	if err != nil {
		response.Error(c, errors.Join(errors.New("创建优惠券失败："), err))
		return
	}

	response.Success(c, []string{})
	return
}

// Delete 删除Coupon分类
func Delete(c *gin.Context, delete *dto.CouponDelete) {
	err := dao.NewCouponDao().Delete(delete)
	if err != nil {
		response.Error(c, errors.Join(errors.New("删除优惠券失败："), err))
		return
	}

	response.Success(c, []string{})
	return
}

// Update 修改Coupon的分类信息
func Update(c *gin.Context, update *dto.CouponUpdate) {
	err := dao.NewCouponDao().Update(update)
	if err != nil {
		response.Error(c, errors.Join(errors.New("更新优惠券失败"), err))
		return
	}

	response.Success(c, []string{})
	return
}

// Search 查询
func Search(c *gin.Context, search *dto.CouponSearch) {
	// id, name, page, limit
	res, err := dao.NewCouponDao().Search(search)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.PaginateSuccess(c, &res)
	return
}
