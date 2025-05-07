package order

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"mall_backend/dao"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/service/coupon"
	"mall_backend/util"
	"strings"
	"time"
)

type NormalStrategy struct {
	finalPrice int
	totalPrice int
	orderCode  string
	userID     int
}

// 这里下单的时候，同一个商户下的多个商品属于同一个订单
// 这里直接当成同一个订单。操他妈的多商户。 以后再根据商户id来扩展

func (n *NormalStrategy) CreateOrder(ctx *Context) (res Result, err error) {
	// 断言获取create
	create, ok := ctx.Param.(*dto.OrderCreate)
	if !ok {
		err = errors.New("创建订单失败：assert *dto.Create fail")
		return
	}
	n.userID = ctx.UserID

	// 通过查询获取到商品的价格
	product, err := dao.NewSkuDao(util.DBClient()).MoreWithOrder(&create.Product)
	if err != nil {
		err = errors.New("创建订单失败：商品检查失败")
		return
	}

	idProduct := make(map[int]model.MmSku)
	for _, sku := range *product {
		idProduct[int(sku.ID)] = sku
	}

	// 计算订单总价
	totalPrice := 0
	for _, item := range *product {
		totalPrice += int(item.Price) * ctx.UniqueSku[int(item.ID)]
	}

	couponLadders, err := dao.NewCouponDao(util.DBClient()).CouponWithLadder(create.Coupons...)
	if err != nil {
		err = errors.New("创建订单失败：查询优惠券失败")
		return
	}

	// 计算数量
	totalUnit := 0
	for _, item := range create.Product {
		totalUnit += item.Num
	}

	// 写的不是很好，感觉这里还是不需要拆分也挺好？不过真需要改成配置么。算了，没必要；浪费时间
	// 看优惠券是不是互斥的。  只需要佛reach,如果 > 2 || == 2 但是不是 platform或者其他的
	if len(*couponLadders) > 0 {
		couponMap := make(map[string]string)
		// 看有没有重复不能叠加优惠券
		for _, item := range *couponLadders {
			if _, ok := couponMap[item.Scope.Code]; ok {
				err = errors.New("创建订单失败： 存在重复的优惠券")
				return
			}
			// 这里直接控制是否叠加使用
			couponMap[item.Scope.Code] = item.Scope.Code
		}
		log.Println(couponMap)
		if len(couponMap) > 2 {
			err = errors.New("创建订单失败： 多种优惠券不能同时叠加")
			return
		}

		if len(couponMap) == 2 {
			// 去平台和商家，取不到就报错
			if _, ok := couponMap[dao.CouponPlatform]; !ok {
				err = errors.New("创建订单失败：只有平台券和商家券可以同时叠加")
				return
			}

			if _, ok := couponMap[dao.CouponMerchant]; !ok {
				err = errors.New("创建订单失败：只有平台券和商家券可以同时叠加")
				return
			}
		}

	}

	// 使用策略模式里面的优惠券
	finalPrice := totalPrice
	log.Println("使用优惠券之前：", totalPrice)
	couponCtx := make(map[int]string) // 为了保存优惠券使用记录
	for _, item := range *couponLadders {
		tx := coupon.Context{
			Price: finalPrice,
			Num:   totalUnit,
			CL:    item,
			CM:    couponCtx,
		}
		item.Scope.Code = strings.ToUpper(item.Scope.Code)

		if item.Scope.Code == dao.CouponPlatform {
			strategy := &coupon.PlatformStrategy{}
			strategy.ApplyStrategy(&tx)
		}

		if item.Scope.Code == dao.CouponMerchant {
			strategy := &coupon.MerchantStrategy{}
			strategy.ApplyStrategy(&tx)
		}

		// TODO 分类券
		if item.Scope.Code == dao.CouponCategory {
		}
		// TODO 单个商品券
		if item.Scope.Code == dao.CouponProduct {
		}

		finalPrice = tx.Price
	}
	log.Println("使用优惠券之后：", finalPrice)
	// 1. 创建订单
	// 2. 扣减优惠券
	// 3. 扣减库存
	// 5. 扣减用户优惠券
	//

	orderCode := util.OrderCode()
	batchCode := util.BatchCode()
	n.orderCode = orderCode
	n.finalPrice = finalPrice
	n.totalPrice = totalPrice
	// 创建订单
	orderCreate := &model.MmOrder{
		BatchCode:     batchCode,
		OrderCode:     orderCode,
		OrderType:     dao.OrderTypeNormal,
		UserID:        int64(int32(ctx.UserID)),
		TotalAmount:   int64(totalPrice),
		PayAmount:     int64(finalPrice),
		AddressID:     int64(int32(create.AddressID)),
		PaymentStatus: dao.OrderStatusNeedPay,
		PaymentWay:    0,
		Source:        int32(create.Source),
		IsSign:        false,
		SignDate:      util.MinDateTime(),
		AfterSale:     0,
		Status:        true,
		CreateTime:    time.Now(),
		UpdateTime:    util.MinDateTime(),
		ExpireTime:    time.Now().Add(time.Minute * 30),
	}
	// 创建使用的订单
	orderCoupon := n.orderCoupon(couponCtx)
	// 拆分商品以及具体的信息到order_spu
	subOrder := n.orderSpu(create, product)

	orderId := 0
	err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		orderId, err = dao.NewOrderDao(tx).Create(orderCreate)
		if err != nil {
			return err
		}

		if err = dao.NewOrderSpuDao(tx).Create(subOrder); err != nil {
			return err
		}
		if err = dao.NewCouponDao(tx).CouponUse(create.Coupons...); err != nil {
			return err
		}

		if err = dao.NewOrderCouponDao(tx).Create(orderCoupon); err != nil {
			return err
		}
		// TODO 应该是是支付的时候才扣库存，只有秒杀是创建订单就扣减库存
		//if err = dao.NewSkuDao(tx).StockDecrease(&stock); err != nil {
		//	return err
		//}
		if err = dao.NewUserCouponDao(tx).Use(ctx.UserID, create.Coupons...); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	return Result{
		ID:        orderId,
		Code:      orderCode,
		BatchCode: batchCode,
	}, nil
}

// orderSpu 拆分订单的商品到子订单，为啥要这么拆分？因为order是每个商家存的同一个订单
// order_spu是为了保存明细，以及每个商品的价格之类的
func (n *NormalStrategy) orderSpu(create *dto.OrderCreate, product *[]model.MmSku) *[]model.MmOrderSpu {
	idProduct := make(map[int]model.MmSku)
	for _, sku := range *product {
		idProduct[int(sku.ID)] = sku
	}
	var subOrder []model.MmOrderSpu
	//stock := []dao.StockUpdate{}
	skuPayAmount := make(map[int]int64)
	var tmpSum int64
	for i, item := range create.Product {
		if i < len(create.Product)-1 {
			tmp := int64((item.Num) * int(idProduct[item.SkuID].Price) * n.finalPrice / n.totalPrice)
			skuPayAmount[item.SkuID] = tmp
			tmpSum += tmp
			continue
		}
		skuPayAmount[item.SkuID] = int64(n.finalPrice) - tmpSum
	}

	for _, item := range create.Product {
		subOrder = append(subOrder, model.MmOrderSpu{
			OrderCode:     n.orderCode,
			UserID:        int32(n.userID),
			SpuID:         int64(item.SpuID),
			SkuID:         int64(item.SkuID),
			UnitAmount:    int64(idProduct[item.SkuID].Price),
			Count:         int32(item.Num),
			PaymentStatus: dao.OrderStatusNeedPay,
			PaymentWay:    0,
			TotalAmount:   int64((item.Num) * int(idProduct[item.SkuID].Price)),
			PayAmount:     skuPayAmount[item.SkuID], // 当前商品 * 数量 * 优惠后的价格 / 总价
			RefundAmount:  0,
			Specs:         idProduct[item.SkuID].Spces,
			IsSign:        0,
			SignDate:      util.MinDateTime(),
			Status:        true,
			CreateTime:    time.Now(),
			UpdateTime:    util.MinDateTime(),
		})
	}
	return &subOrder
}

// orderCoupon 保存用户使用的优惠券  TODO 这里好像有点问题吗
func (n *NormalStrategy) orderCoupon(ladder map[int]string) *[]model.MmOrderCoupon {
	orderCoupon := make([]model.MmOrderCoupon, 0, len(ladder))
	for id, item := range ladder {
		orderCoupon = append(orderCoupon, model.MmOrderCoupon{
			OrderID:    0,
			OrderCode:  n.orderCode,
			CouponID:   int32(id),
			Context:    item,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
		})
	}
	return &orderCoupon
}

// 仅仅使用订单号进行过期处理
func (n *NormalStrategy) orderExpire(orderCode string) error {
	n.orderCode = orderCode
	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		order, err := dao.NewOrderDao(util.DBClient()).One(n.orderCode)
		if err != nil {
			return err
		}

		orderSku, err := dao.NewOrderSpuDao(util.DBClient()).More(n.orderCode)
		if err != nil {
			return err
		}

		// 1. 归还库存
		stocks := make([]dao.StockUpdate, 0, len(*orderSku))
		for _, item := range *orderSku {
			stocks = append(stocks, dao.StockUpdate{
				ID:  int(item.SkuID),
				Num: int(item.Count),
			})
		}
		err = dao.NewSkuDao(tx).StockIncrease(&stocks)
		if err != nil {
			return errors.Join(errors.New(fmt.Sprintf("订单号：%s归还库存失败：", n.orderCode)), err)
		}

		orderCoupon, err := dao.NewOrderCouponDao(util.DBClient()).More(n.orderCode)
		if err != nil {
			return err
		}
		if orderCoupon != nil && len(*orderCoupon) > 0 {
			log.Printf("%#v", orderCoupon)
			// 2. 归还用户优惠券
			coupons := make([]int, 0, len(*orderCoupon))
			for _, item := range *orderCoupon {
				coupons = append(coupons, int(item.CouponID))
			}
			err = dao.NewUserCouponDao(tx).Cancel(int(order.UserID), coupons...)
			if err != nil {
				return errors.Join(errors.New(fmt.Sprintf("订单号：%s归还优惠券失败：", n.orderCode)), err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
