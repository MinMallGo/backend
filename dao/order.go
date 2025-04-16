package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mall_backend/dao/model"
	"mall_backend/dto"
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

func (u *OrderDao) Create(create *model.MmOrder) (int, error) {
	err := u.db.Model(&model.MmOrder{}).Create(create).Error
	if err != nil {
		return 0, err
	}
	return int(create.ID), nil
}

func (u *OrderDao) CheckBeforePay(order *dto.PayOrder, userId int) (*model.MmOrder, error) {
	res := &model.MmOrder{}
	tx := u.db.Model(&model.MmOrder{}).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ?", order.ID, order.OrderCode, userId, false).
		First(res)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return res, errors.New("请选择正确的订单")
	}
	return res, nil
}

func (u *OrderDao) OrderPaid(order *dto.PayOrder, userId int) error {
	// 更新订单状态
	if tx := u.db.Model(&model.MmOrder{}).
		Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ?",
			order.ID, order.OrderCode, userId, false).Updates(map[string]interface{}{
		"payment_status": true,
		"update_time":    time.Now(),
	}); tx.Error != nil {
		return tx.Error
	}
	return nil
}
