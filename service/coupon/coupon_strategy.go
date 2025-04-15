package coupon

import (
	"log"
	"mall_backend/dao"
	"mall_backend/structure"
)

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
	ApplyStrategy(ctx *structure.Context)
}

// 把总价传过来，然后把券的计算规则传过来。就可以进行计算了

type PlatformStrategy struct {
	//
}

func (t *PlatformStrategy) ApplyStrategy(ctx *structure.Context) {
	if ctx.Category == dao.DiscountType {
		FullDiscount(ctx)
	}

	if ctx.Category == dao.RateType {
		DiscountRate(ctx)
	}
	log.Println(ctx.Price)
}

type MerchantStrategy struct {
}

func (m *MerchantStrategy) ApplyStrategy(ctx *structure.Context) {
	log.Printf("%#v\n", ctx)
	if ctx.Category == dao.DiscountType {
		FullDiscount(ctx)
	}

	if ctx.Category == dao.RateType {
		DiscountRate(ctx)
	}
	log.Println(ctx.Price)
}

type CategoryStrategy struct{}

func (c *CategoryStrategy) ApplyStrategy(ctx *structure.Context) {
	if ctx.Category == dao.DiscountType {
		FullDiscount(ctx)
	}

	if ctx.Category == dao.RateType {
		DiscountRate(ctx)
	}
}

// FullDiscount 满减
func FullDiscount(ctx *structure.Context) {
	if ctx.Price > ctx.Rules.Threshold {
		ctx.Price -= ctx.Rules.DisCountPrice
	}
	return
}

func DiscountRate(ctx *structure.Context) {
	ctx.Price = ctx.Price * ctx.Rules.DiscountRate / 100
}
