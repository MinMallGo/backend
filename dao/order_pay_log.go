package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/util"
	"time"
)

type OrderPayLogDao struct {
	db *gorm.DB
}

func NewOrderPayLogDao(db *gorm.DB) *OrderPayLogDao {
	return &OrderPayLogDao{
		db: db,
	}
}

func (d *OrderPayLogDao) Create(orders *[]model.MmOrderSpu, queryString string) error {
	var param []model.MmOrderPayLog
	for _, order := range *orders {
		param = append(param, model.MmOrderPayLog{
			OrderID:        order.ID,
			OrderCode:      order.OrderCode,
			UserID:         order.UserID,
			PayAmount:      order.PayAmount,
			PayWay:         order.PaymentWay,
			ThirdPartyCode: "",
			PayQueryData:   queryString,
			PayTime:        util.MinDateTime(),
		})
	}

	if tx := d.db.Model(&model.MmOrderPayLog{}).Create(param); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *OrderPayLogDao) PaySuccess(orderCode string) error {
	tx := d.db.Model(&model.MmOrderPayLog{}).Where("order_code = ?", orderCode).Updates(map[string]interface{}{
		"pay_time": time.Now(),
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
