package dao

import (
	"errors"
	"fmt"
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

func (d *OrderSpuDao) PaySuccess(orderCode string) error {
	tx := d.db.Model(&model.MmOrderSpu{}).Where("order_code = ? AND payment_status = ?", orderCode, OrderStatusNeedPay).
		Updates(map[string]interface{}{
			"payment_status": OrderStatusPaid,
			"update_time":    time.Now(),
		})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("子订单支付状态修改失败")
	}

	return nil
}

func (d *OrderSpuDao) MoreRefundOrder(orderCode string, skus ...int) (res *[]model.MmOrderSpu, err error) {
	// 需要全部传或者部分传
	if len(skus) == 0 {
		return nil, errors.New("请传入正确的商品规格")
	}
	tx := d.db.Model(&model.MmOrderSpu{}).Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("order_code = ? AND sku_id in ?", orderCode, skus).Find(&res)
	if tx.Error != nil {
		return nil, err
	}

	if tx.RowsAffected != int64(len(skus)) {
		return nil, errors.New("退款查询失败：请选择正确的商品")
	}

	return res, nil
}

func (d *OrderSpuDao) Refund(refundOrder *[]model.MmOrderSpu) error {
	// when then
	wt := " WHEN sku_id = ? THEN refund_amount + ? "
	wtStr := ""
	var args []interface{}
	var skuIDs []int
	for _, item := range *refundOrder {
		wtStr += wt
		args = append(args, item.SkuID, item.PayAmount)
		skuIDs = append(skuIDs, int(item.SkuID))
	}

	tx := d.db.Model(&model.MmOrderSpu{}).
		Where("order_code = ?", (*refundOrder)[0].OrderCode).
		Where("sku_id in ?", skuIDs).
		Update("refund_amount", gorm.Expr(fmt.Sprintf(`
        CASE 
            %s
        END`, wtStr), args...))

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("订单退款失败：更新子订单失败")
	}

	return nil
}
