package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// 管理员注册
	AdminRegister(r)
	// 用户注册
	UserRegister(r)
	// spu 注册
	SpuRegister(r)
	// Sku注册
	SkuRegister(r)
	// 注册规格值路由
	SpecRegister(r)
	// 分类注册
	CategoryRegister(r)
	// 活动注册
	ActiveRegister(r)
	// 注册优惠券路由
	CouponRegister(r)
	// 注册购物车
	CartRegister(r)
}
