package Coupon

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service/coupon"
)

func CreateCoupon(c *gin.Context) {
	create := &dto.CouponCreate{}
	if err := c.ShouldBindJSON(create); err != nil {
		response.Error(c, err)
		return
	}
	coupon.Create(c, create)
}

func UpdateCoupon(c *gin.Context) {
	update := &dto.CouponUpdate{}
	if err := c.ShouldBindJSON(update); err != nil {
		response.Error(c, err)
		return
	}
	coupon.Update(c, update)
}

func SearchCoupon(c *gin.Context) {
	search := &dto.CouponSearch{}
	if err := c.ShouldBindJSON(search); err != nil {
		response.Error(c, err)
		return
	}
	coupon.Search(c, search)
}

func DeleteCoupon(c *gin.Context) {
	del := &dto.CouponDelete{}
	if err := c.ShouldBindJSON(del); err != nil {
		response.Error(c, err)
		return
	}
	coupon.Delete(c, del)
}
