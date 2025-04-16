package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"time"
)

// 订单进程状态：0 待支付，1已支付 2已过期 3申请退款 4已退款 5订单完成

const (
	OrderNeedPay = iota
	OrderPaid
	OrderExpired
	OrderCommitCancel
	OrderCanceled
	OrderCompleted
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
		Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ? AND process_status = ?",
			order.ID, order.OrderCode, userId, false, OrderNeedPay).
		First(res)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return res, errors.New("支付失败：请检查订单是否已支付或已经取消")
	}
	return res, nil
}

func (u *OrderDao) OrderPaid(order *dto.PayOrder, userId int) error {
	// 更新订单状态
	if tx := u.db.Model(&model.MmOrder{}).
		Where("id = ? AND order_code = ? AND user_id = ? AND process_status = ?",
			order.ID, order.OrderCode, userId, false).Updates(map[string]interface{}{
		"payment_status": true,
		"process_status": OrderPaid,
		"update_time":    time.Now(),
	}); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (u *OrderDao) IsPaid(order *dto.CancelOrder, userId int) (*model.MmOrder, error) {
	res := &model.MmOrder{}
	tx := u.db.Model(&model.MmOrder{}).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ? AND process_status = ?",
			order.ID, order.OrderCode, userId, true, OrderPaid).
		First(res)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return res, errors.New("退款失败：请选择正确的订单")
	}
	return res, nil
}

func (u *OrderDao) Cancel(order *dto.CancelOrder, userId int) error {
	res := u.db.Model(&model.MmOrder{}).Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ? AND process_status = ?",
		order.ID, order.OrderCode, userId, true, OrderPaid).Updates(map[string]interface{}{
		"payment_status": false,
		"process_status": OrderCommitCancel,
		"update_time":    time.Now(),
	})
	if res.Error != nil || res.RowsAffected == 0 {
		return errors.New("订单取消失败：更新主订单失败")
	}
	return nil
}

func (u *OrderDao) Expire(order *dto.OrderExpire, userId int) error {
	res := u.db.Model(&model.MmOrder{}).Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ? AND process_status = ?",
		order.ID, order.OrderCode, userId, false, OrderNeedPay).Updates(map[string]interface{}{
		"payment_status": false,
		"process_status": OrderExpired,
		"update_time":    time.Now(),
	})
	if res.Error != nil || res.RowsAffected == 0 {
		return errors.New("订单过期失败：更新主订单失败")
	}
	return nil
}

func (u *OrderDao) One(orderCode string, id int) (*model.MmOrder, error) {
	res := &model.MmOrder{}
	tx := u.db.Model(&model.MmOrder{}).Where("id = ? AND order_code = ?", id, orderCode).Find(res)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return res, tx.Error
	}
	return res, nil
}
