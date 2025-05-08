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

func NewOrderRefundDao(db *gorm.DB) *OrderCancelDao {
	return &OrderCancelDao{
		db: db,
	}
}

func (d *OrderCancelDao) Create(order *model.MmOrder) error {
	param := &model.MmOrderRefundLog{
		OrderID:        order.ID,
		OrderCode:      order.OrderCode,
		BatchCode:      order.BatchCode,
		UserID:         int32(order.UserID),
		CancelAmount:   order.PayAmount,
		CancelWay:      order.PaymentWay,
		CancelAt:       time.Now(),
		ThirdPartyCode: "",
	}
	tx := d.db.Model(&model.MmOrderRefundLog{}).Create(param)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return errors.New("退款失败：写入订单退款日志失败")
	}
	return nil
}
