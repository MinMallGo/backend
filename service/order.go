package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service/coupon"
	"mall_backend/structure"
	"mall_backend/util"
)

// TODO 这里创建订单和下单需要两个步骤，先创建订单，然后才是订单支付。先不要考虑那么精细化

func OrderCreate(c *gin.Context, create *dto.OrderCreate) {
	//AddressID string         `json:"address_id" binding:"required"`
	//Coupons   []string       `json:"coupons" binding:"required"`
	//Product   []ShoppingItem `json:"product" binding:"required,dive"`
	//PayWay    string         `json:"pay_way" binding:"required,oneof=alipay wechatpay"`
	//Source    string         `json:"source" binding:"required,oneof=H5 min-program"`
	// 进行数据的校验
	// TODO 这里需要进行测试就是我在下单的时候，同时修改可用性，看能否使用成功
	// TODO 商品验证，商品规格验证
	// TODO 优惠券
	// TODO 满减
	// 这里只是创建订单
	// 获取用户id
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Failure(c, "下单失败：请登录后操作")
		return
	}

	// 1. 事务启动
	// 2. 检查商品，检查规格
	// 3. 地址检查
	// 4. 验证优惠券并计算价格
	// 5. 构造数据
	spuIds := make([]int, 0, len(create.Product))
	skuIds := make([]int, 0, len(create.Product))
	mSku := make(map[int]int, len(create.Product))
	for _, v := range create.Product {
		spuIds = append(spuIds, v.SpuID)
		skuIds = append(skuIds, v.SkuID)
		mSku[v.SkuID] = v.Num
	}

	if !dao.NewSpuDao(util.DBClient()).Exists(spuIds...) {
		response.Error(c, errors.New("创建订单失败：请选择正确的商品"))
		return
	}
	if !dao.NewSkuDao(util.DBClient()).Exists(skuIds...) {
		response.Error(c, errors.New("创建订单失败：请选择正确是商品规格"))
		return
	}

	if !dao.NewAddrDao(util.DBClient()).ExistsWithUser(int(user.ID), create.AddressID) {
		response.Error(c, errors.New("创建订单失败：请选择正确的地址"))
		return
	}
	if err = dao.NewCouponDao(util.DBClient()).ExistsWithUser(int(user.ID), create.Coupons...); err != nil {
		response.Error(c, errors.Join(errors.New("创建订单失败：请选择正确的优惠券："), err))
		return
	}

	// 通过查询获取到商品的价格
	product, err := dao.NewSkuDao(util.DBClient()).MoreWithOrder(&create.Product)
	if err != nil {
		response.Error(c, errors.New("创建订单失败：优惠券检查失败"))
		return
	}

	// 计算订单总价
	totalPrice := 0
	for _, item := range *product {
		totalPrice += int(item.Price) * mSku[int(item.ID)]
	}

	coupons, err := dao.NewCouponDao(util.DBClient()).More(create.Coupons...)
	if err != nil {
		response.Error(c, errors.New("创建订单失败：查询优惠券失败"))
		return
	}

	finalPrice := totalPrice
	for _, item := range *coupons {
		ctx := structure.Context{
			Price:    finalPrice,
			Category: int(item.DiscountType),
			Rules: structure.Rule{
				Threshold:     int(item.Threshold),
				DisCountPrice: int(item.DiscountPrice),
				DiscountRate:  int(item.DiscountRate),
			},
		}

		if item.UseRange == dao.RangeAll {
			strategy := &coupon.PlatformStrategy{}
			strategy.ApplyStrategy(&ctx)
		}
		if item.UseRange == dao.RangeMerchant {
			strategy := &coupon.MerchantStrategy{}
			strategy.ApplyStrategy(&ctx)
			log.Println("after use coupon:", ctx.Price)
		}

		finalPrice = ctx.Price
	}

	err = dao.NewOrderDao(util.DBClient()).Create(create, int(user.ID), totalPrice, finalPrice)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, []string{})
	return
}

func OrderPay(c *gin.Context) {
	// TODO 这里是支付订单，看有没有第三方的对接
	// TODO 这里是消息通知  可能需要拆分。先弄个邮件发送得了
}
