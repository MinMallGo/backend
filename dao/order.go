package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"time"
)

const (
	OrderStatusNeedPay       = iota + 1 // 1. 待支付
	OrderStatusPaid                     // 2. 已支付
	OrderStatusPartCancel               // 3. 部分退款成功
	OrderStatusAllCancel                //4. 全额退款成功
	OrderStatusPartCanceling            //5. 部分申请退款中
	OrderStatusAllCanceling             //6. 全额申请退款中
)

// 定义订单类型
const (
	OrderTypeNormal = 1 + iota
	OrderTypeSecKill
)

// 定义支付方式
const (
	PayByAlipay = 1 + iota
	PayByWechat
	PayByBankCARD
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

// CheckBeforePay 支付前检查
func (u *OrderDao) CheckBeforePay(order *dto.PayOrder, userId int) (*model.MmOrder, error) {
	res := &model.MmOrder{}
	tx := u.db.Model(&model.MmOrder{}).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("batch_code = ? AND order_code = ? AND user_id = ? AND payment_status = ?",
			order.BatchCode, order.OrderCode, userId, OrderStatusNeedPay).
		First(res)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return res, errors.New("支付失败：请检查订单是否已支付或已经取消")
	}
	return res, nil
}

func (u *OrderDao) OrderPaid(order *dto.PayOrder, userId int) error {
	// 更新订单状态
	if tx := u.db.Model(&model.MmOrder{}).
		Where("batch_code = ? AND order_code = ? AND user_id = ? AND payment_status = ?",
			order.BatchCode, order.OrderCode, userId, OrderStatusNeedPay).Updates(map[string]interface{}{
		"payment_status": OrderStatusPaid,
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
			order.ID, order.OrderCode, userId, true, OrderStatusPaid).
		First(res)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return res, errors.New("退款失败：请选择正确的订单")
	}
	return res, nil
}

func (u *OrderDao) Cancel(order *dto.CancelOrder, userId int) error {
	res := u.db.Model(&model.MmOrder{}).Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ? AND process_status = ?",
		order.ID, order.OrderCode, userId, true, OrderStatusPaid).Updates(map[string]interface{}{
		"payment_status": false,
		//"process_status": OrderCommitCancel,
		"update_time": time.Now(),
	})
	if res.Error != nil || res.RowsAffected == 0 {
		return errors.New("订单取消失败：更新主订单失败")
	}
	return nil
}

func (u *OrderDao) Expire(order *dto.OrderExpire, userId int) error {
	res := u.db.Model(&model.MmOrder{}).Where("id = ? AND order_code = ? AND user_id = ? AND payment_status = ? AND process_status = ?",
		order.ID, order.OrderCode, userId, false, OrderStatusNeedPay).Updates(map[string]interface{}{
		"payment_status": false,
		//"process_status": OrderExpired,
		"update_time": time.Now(),
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
