package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mall_backend/dao"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service/coupon"
	"mall_backend/structure"
	"mall_backend/util"
	"strconv"
	"strings"
	"time"
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
	spuIds := make([]int, 0, len(create.Product)) // 这里要去重
	skuIds := make([]int, 0, len(create.Product))
	mSku := make(map[int]int, len(create.Product))
	seed := map[int]struct{}{}
	hasSecKill := false
	for _, v := range create.Product {
		if _, ok := seed[v.SpuID]; !ok {
			seed[v.SpuID] = struct{}{}
			spuIds = append(spuIds, v.SpuID)
		}
		skuIds = append(skuIds, v.SkuID)
		mSku[v.SkuID] = v.Num
		if v.IsSecKill == 1 {
			hasSecKill = true
		}
	}

	if hasSecKill && len(create.Coupons) > 0 {
		response.Failure(c, "秒杀商品不能合并下单")
		return
	}

	if !dao.NewSpuDao(util.DBClient()).Exists(spuIds...) {
		response.Error(c, errors.New("创建订单失败：请选择正确的商品"))
		return
	}
	if !dao.NewSkuDao(util.DBClient()).Exists(skuIds...) {
		response.Error(c, errors.New("创建订单失败：请选择正确的商品规格"))
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

	if err = dao.NewUserCouponDao(util.DBClient()).Exists(int(user.ID), create.Coupons...); err != nil {
		response.Error(c, errors.Join(errors.New("创建订单失败：请选择检查个人优惠券："), err))
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

	// 使用策略模式里面的优惠券
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
		}

		finalPrice = ctx.Price
	}

	// 1. 创建订单
	// 2. 扣减优惠券
	// 3. 扣减库存
	// 5. 扣减用户优惠券 // TODO 别想这么多啊，一个优惠券就只能有一次
	//

	couponIds := ""
	for i, v := range create.Coupons {
		couponIds += fmt.Sprintf("%d", v)
		if i < len(create.Coupons)-1 {
			couponIds += ","
		}
	}

	orderCode := util.OrderCode()

	order := &model.MmOrder{
		OrderCode:     orderCode,
		UserID:        int32(user.ID),
		OriginalPrice: int32(totalPrice),
		AddressID:     int32(create.AddressID),
		FinalPrice:    int32(finalPrice),
		Coupons:       couponIds,
		PaymentStatus: 0,
		PaymentWay:    0,
		ProcessStatus: dao.OrderNeedPay,
		Source:        int32(create.Source),
		IsSign:        0,
		SignDate:      util.MinDateTime(),
		Status:        true,
		CreateTime:    time.Now(),
		UpdateTime:    util.MinDateTime(),
		DeleteTime:    util.MinDateTime(),
	}

	subOrder := []model.MmSubOrder{}
	stock := []dao.StockUpdate{}
	for _, item := range create.Product {
		subOrder = append(subOrder, model.MmSubOrder{
			OrderCode:       orderCode,
			OrderUniqueCode: util.SubOrderCode(),
			SpuID:           int32(item.SpuID),
			SkuID:           int32(item.SkuID),
			Nums:            int32(item.Num),
			OriginalPrice:   int32(totalPrice),
			FinalPrice:      int32(finalPrice),
			AddressID:       int32((*create).AddressID),
			CouponID:        0,
			PaymentStatus:   0,
			PaymentWay:      0,
			IsSign:          0,
			SignDate:        util.MinDateTime(),
			Status:          false,
			CreateTime:      time.Now(),
			UpdateTime:      util.MinDateTime(),
			DeleteTime:      util.MinDateTime(),
		})
		stock = append(stock, dao.StockUpdate{
			ID:  item.SkuID,
			Num: item.Num,
		})
	}

	orderId := 0
	err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		orderId, err = dao.NewOrderDao(tx).Create(order)
		if err != nil {
			return err
		}
		if err = dao.NewSubOrderDao(tx).Create(&subOrder); err != nil {
			return err
		}
		if err = dao.NewCouponDao(tx).CouponUse(create.Coupons...); err != nil {
			return err
		}
		if err = dao.NewSkuDao(tx).StockDecrease(&stock); err != nil {
			return err
		}
		if err = dao.NewUserCouponDao(tx).Use(int(user.ID), create.Coupons...); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, map[string]interface{}{
		"orderId":   orderId,
		"orderCode": orderCode,
	})
	return
}

func OrderPay(c *gin.Context, order *dto.PayOrder) {
	// TODO 这里是支付订单，看有没有第三方的对接
	// TODO 这里是消息通知  可能需要拆分。先弄个邮件发送得了
	// TODO
	// TODO 检查并支付  直接用事务吧
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	if err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		res, err := dao.NewOrderDao(tx).CheckBeforePay(order, int(user.ID))
		if err != nil {
			return err
		}

		// TODO 支付对接
		time.Sleep(time.Second * 5)
		// TODO 消息通知

		if err = dao.NewOrderDao(tx).OrderPaid(order, int(user.ID)); err != nil {
			return err
		}

		if err = dao.NewOrderPaidLogDao(tx).Create(res); err != nil {
			return err
		}

		// 更新订单的状态
		return nil
	}); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, []string{})
	return
}

func CancelOrder(c *gin.Context, order *dto.CancelOrder) {
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}
	// 取消订单的步骤：
	// 1. 查看订单是否支付。 IsPaid
	// 2. 归还库存
	// 3. 归还优惠券
	// 4. 归还用户优惠券
	// 5.  TODO 如果支付了的话，要走退款流程
	// 6. 更改订单状态
	res := &model.MmOrder{}
	err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		// 查询订单是否支付
		res, err = dao.NewOrderDao(tx).IsPaid(order, int(user.ID))
		if err != nil {
			return err
		}
		// TODO 退款流程
		time.Sleep(time.Second * 5)

		// 通过主订单查询子订单，并修改子订单的状态以及归还商品库存
		subOrders, err := dao.NewSubOrderDao(tx).More(res.OrderCode)
		if err != nil {
			return err
		}
		// 归还商品库存
		if err = dao.NewSkuDao(tx).IncreaseStock(subOrders); err != nil {
			return err
		}
		// 修改子订单状态
		if err = dao.NewSubOrderDao(tx).Cancel(res.OrderCode); err != nil {
			return err
		}
		// 修改住主订单状态
		if err = dao.NewOrderDao(tx).Cancel(order, int(user.ID)); err != nil {
			return err
		}
		// 1. 归还优惠券
		// 2. 归还用户优惠券
		if len(res.Coupons) > 0 {
			tmp := strings.Split(res.Coupons, ",")
			ids := make([]int, 0, len(tmp))
			for _, couponId := range tmp {
				id, err := strconv.Atoi(couponId)
				if err != nil {
					return errors.New("归还优惠券失败，请联系客服")
				}
				ids = append(ids, id)
			}
			if err = dao.NewCouponDao(tx).Cancel(ids...); err != nil {
				return err
			}
			if err = dao.NewUserCouponDao(tx).Cancel(int(user.ID), ids...); err != nil {
				return err
			}
		}
		// 记录退款日志
		if err = dao.NewOrderCancelDao(tx).Create(res, int(user.ID)); err != nil {
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

func OrderExpire(c *gin.Context, order *dto.OrderExpire) {
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	// 1.归还库存，
	// 2.这玩意儿
	// 取消订单的步骤：
	// 1. 查看订单是否支付。 IsPaid
	// 2. 归还库存
	// 3. 归还优惠券
	// 4. 归还用户优惠券
	// 5.  TODO 如果支付了的话，要走退款流程
	// 6. 更改订单状态
	err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		// 通过主订单查询子订单，并修改子订单的状态以及归还商品库存
		subOrders, err := dao.NewSubOrderDao(tx).More(order.OrderCode)
		if err != nil {
			return err
		}
		// 归还商品库存
		if err = dao.NewSkuDao(tx).IncreaseStock(subOrders); err != nil {
			return err
		}
		// 修改子订单状态
		if err = dao.NewSubOrderDao(tx).Cancel(order.OrderCode); err != nil {
			return err
		}
		// 修改主订单状态
		if err = dao.NewOrderDao(tx).Expire(order, int(user.ID)); err != nil {
			return err
		}
		// 1. 归还优惠券
		// 2. 归还用户优惠券
		res, err := dao.NewOrderDao(tx).One(order.OrderCode, order.ID)
		if err != nil {
			return err
		}

		if len(res.Coupons) > 0 {
			tmp := strings.Split(res.Coupons, ",")
			ids := make([]int, 0, len(tmp))
			for _, couponId := range tmp {
				id, err := strconv.Atoi(couponId)
				if err != nil {
					return errors.New("归还优惠券失败，请联系客服")
				}
				ids = append(ids, id)
			}
			if err = dao.NewCouponDao(tx).Cancel(ids...); err != nil {
				return err
			}
			if err = dao.NewUserCouponDao(tx).Cancel(int(user.ID), ids...); err != nil {
				return err
			}
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
