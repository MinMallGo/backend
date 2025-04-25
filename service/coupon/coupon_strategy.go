package coupon

import (
	"encoding/json"
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
	CL    dao.CouponLadders // 优惠券以及优惠券类型,优惠券梯度
	CM    map[int]string    // 保存哪个优惠券使用的哪个优惠内容
}

type Strategies interface {
	ApplyStrategy(ctx *Context)
}

// 把总价传过来，然后把券的计算规则传过来。就可以进行计算了

type PlatformStrategy struct{}

func (t *PlatformStrategy) ApplyStrategy(ctx *Context) {
	log.Println("类型：", int(ctx.CL.DiscountType))
	if int(ctx.CL.DiscountType) == dao.DiscountType {
		FullDiscount(ctx)
	}
	log.Println("满减后：", int(ctx.CL.DiscountType))

	if int(ctx.CL.DiscountType) == dao.RateType {
		DiscountRate(ctx)
	}
	log.Println("满折扣后：", int(ctx.CL.DiscountType))
	log.Println(ctx.Price)
}

type MerchantStrategy struct{}

func (m *MerchantStrategy) ApplyStrategy(ctx *Context) {
	log.Println("类型：", int(ctx.CL.DiscountType))
	if int(ctx.CL.DiscountType) == dao.DiscountType {
		FullDiscount(ctx)
	}
	log.Println("满减后：", int(ctx.CL.DiscountType))
	if int(ctx.CL.DiscountType) == dao.RateType {
		DiscountRate(ctx)
	}
	log.Println("满折扣后：", int(ctx.CL.DiscountType))
	log.Println(ctx.Price)
}

type CategoryStrategy struct{}

func (c *CategoryStrategy) ApplyStrategy(ctx *Context) {
	// TODO 2025年4月25日 10点36分 如果是分类则需要单独计算商品下的分类
	// TODO 然后从总价里面减去，计算的时候也是通过优惠券算
	if int(ctx.CL.DiscountType) == dao.DiscountType {
		FullDiscount(ctx)
	}

	if int(ctx.CL.DiscountType) == dao.RateType {
		DiscountRate(ctx)
	}
}

// FullDiscount 满减 // 满500 - 200 // 满3件-200
func FullDiscount(ctx *Context) {
	// 如果有梯度，还要根据梯度来选择最优解
	thresholdTValue, disCountValue := getDisCount(ctx)

	//log.Println("<<<<<<<<<<<<", thresholdTValue, disCountValue)
	// 如果是金额满减则比较金额，否则比较数量
	if int(ctx.CL.ThresholdType) == dao.ThresholdAmount {
		//
		if ctx.Price >= thresholdTValue {
			ctx.Price -= disCountValue
		}
	}
	//log.Println("<<<<<<<<<<<<", ctx.Price)
	if int(ctx.CL.ThresholdType) == dao.ThresholdUnit {
		if ctx.Num >= thresholdTValue {
			ctx.Price -= disCountValue
		}
	}
	//log.Println("<<<<<<<<<<<<", ctx.Price)
	return
}

// DiscountRate 满折 满500打8折，满 几件打几折
// 优惠券还可能有梯度，这里需要算用哪个梯度，否则用原来的
func DiscountRate(ctx *Context) {
	// 如果有梯度，还要根据梯度来选择最优解
	thresholdTValue, disCountValue := getDisCount(ctx)
	//log.Println(">>>>>>", thresholdTValue, disCountValue)
	// 判断是否满足条件 这里不能合并到一起，因为一个是看数量，一个是看价格。别瞎吉儿改了
	if int(ctx.CL.ThresholdType) == dao.ThresholdAmount {
		if ctx.Price >= thresholdTValue {
			ctx.Price = ctx.Price * disCountValue / 100
		}
	}

	// 判断是不是满几件商品打折 TODO 这里以后如果是多商户的话，还需要进一步细分
	if int(ctx.CL.ThresholdType) == dao.ThresholdUnit {
		if ctx.Num >= thresholdTValue {
			ctx.Price = ctx.Price * disCountValue / 100
		}
	}

	//log.Println(">>>>>>", ctx.Price)

	return
}

// 返传入多梯度的优惠券的规则（ctx.CL.Ladder），返回门槛和折扣。
func getDisCount(ctx *Context) (int, int) {
	thresholdTValue := int(ctx.CL.ThresholdValue)
	disCountValue := int(ctx.CL.DiscountValue)
	if len(ctx.CL.Ladders) > 0 {
		// 这里就循环来获取值咯
		for _, l := range ctx.CL.Ladders {
			if ctx.Price > int(l.ThresholdValue) {
				thresholdTValue = int(l.ThresholdType)
				disCountValue = int(l.DiscountValue)
			}
		}
	}
	// 记录一下用的那张优惠券，放在context里面
	recordCoupon(ctx, thresholdTValue, disCountValue)
	return thresholdTValue, disCountValue
}

// 记录用的优惠券信息，以及如果是多梯度的话，记录用的哪个梯度的优惠券
func recordCoupon(ctx *Context, thresholdValue, disCountValue int) {
	x := map[int]map[string]int{
		int(ctx.CL.ID): {
			"thresholdType":   int(ctx.CL.ThresholdType),
			"thresholdTValue": thresholdValue,
			"discountType":    int(ctx.CL.DiscountType),
			"disCountValue":   disCountValue,
		},
	}
	res, _ := json.Marshal(x)
	ctx.CM[int(ctx.CL.ID)] = string(res)
}
