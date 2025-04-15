package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"time"
)

type OrderPaidLogDao struct {
	db *gorm.DB
}

func NewOrderPaidLogDao(db *gorm.DB) *OrderPaidLogDao {
	return &OrderPaidLogDao{
		db: db,
	}
}

func (d *OrderPaidLogDao) Create(order *model.MmOrder) error {
	param := &model.MmOrderPaidLog{
		OrderID:   order.ID,
		OrderCode: order.OrderCode,
		UserID:    order.UserID,
		PaidPrice: order.FinalPrice,
		PaidAt:    time.Now(),
	}
	if tx := d.db.Model(&model.MmOrderPaidLog{}).Create(param); tx.Error != nil {
		return tx.Error
	}
	return nil
}
