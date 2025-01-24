package routes

import (
	"github.com/gin-gonic/gin"
	coupon "mall_backend/api/v1/coupon"
	"mall_backend/middleware"
)

func CouponRegister(r *gin.Engine) {
	v1 := r.Group("v1").Use(middleware.AuthVerify())
	{
		// coupon的操作 // TODO 这个玩意儿应该是放在Admin下面的。或者说是商户下面
		v1.POST("/coupon/create", coupon.CreateCoupon)
		v1.POST("/coupon/delete", coupon.DeleteCoupon)
		v1.POST("/coupon/update", coupon.UpdateCoupon)
		v1.POST("/coupon/search", coupon.SearchCoupon)
		// spec的操作
	}
}
