package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"time"
)

type OrderCancelDao struct {
	db *gorm.DB
}

func NewOrderCancelDao(db *gorm.DB) *OrderCancelDao {
	return &OrderCancelDao{
		db: db,
	}
}

func (d *OrderCancelDao) Create(order *model.MmOrder, userId int) error {
	param := &model.MmOrderCancelLog{
		OrderID:        order.ID,
		OrderCode:      order.OrderCode,
		BatchCode:      order.BatchCode,
		UserID:         int32(userId),
		CancelAmount:   order.PayAmount,
		CancelWay:      order.PaymentWay,
		CancelAt:       time.Now(),
		ThirdPartyCode: "",
		PayQueryData:   "",
	}
	tx := d.db.Model(&model.MmOrderCancelLog{}).Create(param)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return errors.New("退款失败：写入订单退款日志失败")
	}
	return nil
}
