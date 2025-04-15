package dao

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

type OrderDao struct {
	db *gorm.DB
}

func NewOrderDao(db *gorm.DB) *OrderDao {
	return &OrderDao{
		db: db,
	}
}

func (u *OrderDao) Create(create *dto.OrderCreate, userId int, originalPrice, finalPrice int) error {
	// 1. 创建订单
	// 2. 扣减优惠券
	// 3. 扣减库存
	// 5. 扣减用户优惠券 // TODO 别想这么多啊，一个优惠券就只能有一次
	//

	coupons := ""
	for i, v := range create.Coupons {
		coupons += fmt.Sprintf("%d", v)
		if i < len(create.Coupons)-1 {
			coupons += ","
		}
	}

	orderCode := util.OrderCode()

	order := &model.MmOrder{
		OrderCode:     orderCode,
		UserID:        int32(userId),
		OriginalPrice: int32(originalPrice),
		AddressID:     int32(create.AddressID),
		FinalPrice:    int32(finalPrice),
		Coupons:       coupons,
		PaymentStatus: 0,
		PaymentWay:    0,
		Source:        int32(create.Source),
		IsSign:        0,
		SignDate:      util.MinDateTime(),
		Status:        true,
		CreateTime:    time.Now(),
		UpdateTime:    util.MinDateTime(),
		DeleteTime:    util.MinDateTime(),
	}

	subOrder := []model.MmSubOrder{}
	stock := []StockUpdate{}
	for _, item := range create.Product {
		subOrder = append(subOrder, model.MmSubOrder{
			OrderCode:       orderCode,
			OrderUniqueCode: util.SubOrderCode(),
			SpuID:           int32(item.SpuID),
			SkuID:           int32(item.SkuID),
			Nums:            int32(item.Num),
			OriginalPrice:   int32(originalPrice),
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
		stock = append(stock, StockUpdate{
			ID:  item.SkuID,
			Num: item.Num,
		})
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.MmOrder{}).Create(&order).Error
		if err != nil {
			return err
		}
		if err = tx.Model(&model.MmSubOrder{}).Create(&subOrder).Error; err != nil {
			return err
		}
		if err = NewCouponDao(u.db).CouponUse(create.Coupons...); err != nil {
			return err
		}
		if err = NewSkuDao(u.db).StockDecrease(&stock); err != nil {
			return err
		}
		if err = NewUserCouponDao(u.db).Use(userId, create.Coupons...); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *OrderDao) OrderPay(c *gin.Context) {}
