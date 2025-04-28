package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mall_backend/dao/model"
	"time"
)

type OrderSpuDao struct {
	db *gorm.DB
}

func NewOrderSpuDao(db *gorm.DB) *OrderSpuDao {
	return &OrderSpuDao{db: db}
}

func (d *OrderSpuDao) Create(create *[]model.MmOrderSpu) error {
	err := d.db.Model(&model.MmOrderSpu{}).Create(create).Error
	if err != nil {
		return errors.New("创建订单失败：拆分子订单失败")
	}
	return nil
}

func (d *OrderSpuDao) More(orderCode string) (res *[]model.MmOrderSpu, err error) {
	err = d.db.Model(&model.MmOrderSpu{}).Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("order_code = ?", orderCode).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *OrderSpuDao) Cancel(orderCode string) (err error) {
	err = d.db.Model(&model.MmOrderSpu{}).Where("order_code = ?", orderCode).Updates(map[string]interface{}{
		"payment_status": OrderStatusNeedPay,
		"update_time":    time.Now(),
	}).Error
	return
}
