package dao

import (
	"errors"
	"fmt"
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

// IsPaid 查询订单是否没有支付。
func (u *OrderDao) IsPaid(orderCode string) (bool, error) {
	tx := u.db.Model(&model.MmOrder{}).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("order_code = ? AND payment_status = ?",
			orderCode, OrderStatusNeedPay).
		First(&model.MmOrder{})
	if tx.Error != nil {
		return false, tx.Error
	}

	if tx.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (u *OrderDao) Expire(orderCode string) error {
	res := u.db.Model(&model.MmOrder{}).Where("order_code = ? AND payment_status = ?",
		orderCode, false, OrderStatusNeedPay).Updates(map[string]interface{}{
		"payment_status": false,
		"update_time":    time.Now(),
	})
	if res.Error != nil || res.RowsAffected == 0 {
		return errors.New("订单过期失败：更新主订单失败")
	}
	return nil
}

// IsExpire 订单少于5秒钟不让支付
func (u *OrderDao) IsExpire(orderCode string) error {
	res := &model.MmOrder{}
	tx := u.db.Model(&model.MmOrder{}).Where("order_code = ? AND payment_status = ? AND expire_time > ?",
		orderCode, OrderStatusNeedPay, time.Now()).
		Find(res)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("订单支付提交异常")
	}

	if res.PaymentStatus >= OrderStatusPaid {
		return errors.New("订单已中支付，请勿重复提交")
	}

	if time.Now().After(res.ExpireTime.Add(-3 * time.Second)) {
		return errors.New("订单已过期")
	}

	return nil
}

func (u *OrderDao) One(orderCode string) (*model.MmOrder, error) {
	res := &model.MmOrder{}
	tx := u.db.Model(&model.MmOrder{}).Where("order_code = ? AND status = ?", orderCode, true).Find(res)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return res, tx.Error
	}
	return res, nil
}

func (u *OrderDao) PaySuccess(orderCode string) error {
	tx := u.db.Model(&model.MmOrder{}).Where("order_code = ? AND payment_status = ?", orderCode, OrderStatusNeedPay).
		Updates(map[string]interface{}{
			"payment_status": OrderStatusPaid,
			"update_time":    time.Now(),
		})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("订单支付状态修改失败")
	}

	return nil
}

func (u *OrderDao) Refund(orderCode string, cancelAmount int64) error {
	tx := u.db.Model(&model.MmOrder{}).
		Where("pay_amount <= ?", cancelAmount).
		Where("order_code = ?", orderCode).
		Update("payment_status", gorm.Expr(`
        CASE 
            WHEN pay_amount = ? THEN `+fmt.Sprintf("%d", OrderStatusAllCancel)+` 
            WHEN pay_amount < ? THEN `+fmt.Sprintf("%d", OrderStatusPartCancel)+` 
            ELSE payment_status 
        END`, cancelAmount, cancelAmount))

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("订单退款失败：更新主订单失败")
	}

	return nil
}
