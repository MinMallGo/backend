package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
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
