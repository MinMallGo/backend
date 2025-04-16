package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mall_backend/dao/model"
	"time"
)

type SubOrderDao struct {
	db *gorm.DB
}

func NewSubOrderDao(db *gorm.DB) *SubOrderDao {
	return &SubOrderDao{db: db}
}

func (d *SubOrderDao) Create(create *[]model.MmSubOrder) error {
	err := d.db.Model(&model.MmSubOrder{}).Create(create).Error
	if err != nil {
		return errors.New("创建订单失败：拆分子订单失败")
	}
	return nil
}

func (d *SubOrderDao) More(orderCode string) (res *[]model.MmSubOrder, err error) {
	err = d.db.Model(&model.MmSubOrder{}).Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("order_code = ?", orderCode).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *SubOrderDao) Cancel(orderCode string) (err error) {
	err = d.db.Model(&model.MmSubOrder{}).Where("order_code = ?", orderCode).Updates(map[string]interface{}{
		"payment_status": OrderNeedPay,
		"update_time":    time.Now(),
	}).Error
	return
}
