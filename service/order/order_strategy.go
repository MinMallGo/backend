package order

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mall_backend/dao"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/service/coupon"
	"mall_backend/structure"
	"mall_backend/util"
	"time"
)

const (
	NormalOrder = iota
	SeckillOrder
)

type Context struct {
	Param     interface{}
	UniqueSku map[int]int
	UserID    int
}

type Result struct {
	ID   int
	Code string
}

type Strategy interface {
	Purchase(ctx *Context) (Result, error)
}

func GetStrategy(orderType int) Strategy {
	switch orderType {
	case NormalOrder:
		return &NormalStrategy{}
	case SeckillOrder:
		return &SeckillStrategy{}
	default:
		panic("not support order")
	}
}

type NormalStrategy struct{}

func (n *NormalStrategy) Purchase(ctx *Context) (res Result, err error) {
	// 断言获取create
	create, ok := ctx.Param.(*dto.OrderCreate)
	if !ok {
		err = errors.New("创建订单失败：assert *dto.OrderCreate fail")
		return
	}

	// 通过查询获取到商品的价格
	product, err := dao.NewSkuDao(util.DBClient()).MoreWithOrder(&create.Product)
	if err != nil {
		err = errors.New("创建订单失败：优惠券检查失败")
		return
	}

	// 计算订单总价
	totalPrice := 0
	for _, item := range *product {
		totalPrice += int(item.Price) * ctx.UniqueSku[int(item.ID)]
	}

	coupons, err := dao.NewCouponDao(util.DBClient()).More(create.Coupons...)
	if err != nil {
		err = errors.New("创建订单失败：查询优惠券失败")
		return
	}

	// 使用策略模式里面的优惠券
	finalPrice := totalPrice
	for _, item := range *coupons {
		tx := structure.Context{
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
			strategy.ApplyStrategy(&tx)
		}
		if item.UseRange == dao.RangeMerchant {
			strategy := &coupon.MerchantStrategy{}
			strategy.ApplyStrategy(&tx)
		}

		finalPrice = tx.Price
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
	orderCreate := &model.MmOrder{
		OrderCode:     orderCode,
		UserID:        int32(ctx.UserID),
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
		orderId, err = dao.NewOrderDao(tx).Create(orderCreate)
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
		if err = dao.NewUserCouponDao(tx).Use(int(ctx.UserID), create.Coupons...); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	return Result{
		ID:   orderId,
		Code: orderCode,
	}, nil
}

type SeckillStrategy struct{}

func (s *SeckillStrategy) Purchase(ctx *Context) (res Result, err error) {

	return Result{}, nil
}
