package coupon

import (
	"log"
	"mall_backend/dao"
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

type Context struct {
	Price int               // 原价
	Num   int               // 如果是是件数，还需要传入件数，比如满几件打几折
	cl    dao.CouponLadders // 优惠券以及优惠券类型
}

type Strategies interface {
	ApplyStrategy(ctx *Context)
}

// 把总价传过来，然后把券的计算规则传过来。就可以进行计算了

type PlatformStrategy struct{}

func (t *PlatformStrategy) ApplyStrategy(ctx *Context) {
	if int(ctx.cl.DiscountType) == dao.DiscountType {
		FullDiscount(ctx)
	}

	if int(ctx.cl.DiscountType) == dao.RateType {
		DiscountRate(ctx)
	}
	log.Println(ctx.Price)
}

type MerchantStrategy struct{}

func (m *MerchantStrategy) ApplyStrategy(ctx *Context) {
	log.Printf("%#v\n", ctx)
	if int(ctx.cl.DiscountType) == dao.DiscountType {
		FullDiscount(ctx)
	}

	if int(ctx.cl.DiscountType) == dao.RateType {
		DiscountRate(ctx)
	}
	log.Println(ctx.Price)
}

type CategoryStrategy struct{}

func (c *CategoryStrategy) ApplyStrategy(ctx *Context) {
	if int(ctx.cl.DiscountType) == dao.DiscountType {
		FullDiscount(ctx)
	}

	if int(ctx.cl.DiscountType) == dao.RateType {
		DiscountRate(ctx)
	}
}

// FullDiscount 满减 // 满500 - 200 // 满3件-200
func FullDiscount(ctx *Context) {

	// 如果是金额满减则比较金额，否则比较数量
	if int(ctx.cl.ThresholdType) == dao.ThresholdAmount {
		//
		if ctx.Price >= int(ctx.cl.ThresholdValue) {
			ctx.Price -= int(ctx.cl.DiscountValue)
		}
	} else {
		if ctx.Num >= int(ctx.cl.ThresholdValue) {
			ctx.Price -= int(ctx.cl.DiscountValue)
		}
	}

	return
}

// DiscountRate 满折 满500打8折，满 几件打几折
func DiscountRate(ctx *Context) {
	// 判断是否满足条件
	if int(ctx.cl.ThresholdType) == dao.ThresholdAmount {
		if ctx.Price >= int(ctx.cl.ThresholdValue) {
			ctx.Price -= ctx.Price * int(ctx.cl.DiscountValue) / 100
		}
	} else {
		if ctx.Num >= int(ctx.cl.ThresholdValue) {
			ctx.Price -= ctx.Price * int(ctx.cl.DiscountValue) / 100
		}
	}

	return
}
