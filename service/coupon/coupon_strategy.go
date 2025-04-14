package coupon

import "mall_backend/dao"

/**
策略模式来计算使用优惠券之后的价格
应该是从券的类型去抽象而不是满减或者折扣
原价
优惠券类型： 平台商家可以叠加
商品参数： 包含商品id，商品数量，商品分类，所属商户，

多次调用以计算价格
先满减，再打折。如果可以的话，先平台，后商户
*/

type Strategies interface {
	ApplyStrategy(ctx *Context)
}

// 把总价传过来，然后把券的计算规则传过来。就可以进行计算了

type Context struct {
	Price    int  // 原价
	Category int  // 优惠券类型
	Rules    Rule //  规则类型
}

type Rule struct {
	Threshold     int
	DisCountPrice int
	DiscountRate  int
}

type PlatformStrategy struct {
	//
}

func (t *PlatformStrategy) ApplyStrategy(ctx *Context) {
	if ctx.Category == dao.DiscountType {
		FullDiscount(ctx)
	}

	if ctx.Category == dao.RateType {
		DiscountRate(ctx)
	}
}

type MerchantStrategy struct {
}

func (m *MerchantStrategy) ApplyStrategy(ctx *Context) {
	if ctx.Category == dao.DiscountType {
		FullDiscount(ctx)
	}

	if ctx.Category == dao.RateType {
		DiscountRate(ctx)
	}
}

type CategoryStrategy struct{}

func (c *CategoryStrategy) ApplyStrategy(ctx *Context) {
	if ctx.Category == dao.DiscountType {
		FullDiscount(ctx)
	}

	if ctx.Category == dao.RateType {
		DiscountRate(ctx)
	}
}

// FullDiscount 满减
func FullDiscount(ctx *Context) {
	if ctx.Price > ctx.Rules.Threshold {
		ctx.Price -= ctx.Rules.DisCountPrice
	}
	return
}

func DiscountRate(ctx *Context) {
	ctx.Price = ctx.Price * ctx.Rules.DiscountRate / 100
}
