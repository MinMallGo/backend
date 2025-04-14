package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
)

type OrderDao struct {
	db *gorm.DB
	m  model.MmOrder
}

func NewOrderDao() *OrderDao {
	return &OrderDao{
		db: util.DBClient(),
		m:  model.MmOrder{},
	}
}

func (u OrderDao) Create(create dto.OrderCreate, userId int) error {
	// 1. 事务启动
	// 2. 检查商品，检查规格
	// 3. 地址检查
	// 4. 验证优惠券并计算价格
	// 5. 构造数据
	spuIds := make([]int, 0, len(create.Product))
	skuIds := make([]int, 0, len(create.Product))
	for _, v := range create.Product {
		spuIds = append(spuIds, v.SpuID)
		skuIds = append(skuIds, v.SkuID)
	}

	if !NewSpuDao().Exists(spuIds...) {
		return errors.New("创建订单失败：请选择正确的商品")
	}
	if !NewSkuDao().Exists(skuIds...) {
		return errors.New("创建订单失败：请选择正确是商品规格")
	}
	if !NewAddrDao().ExistsWithUser(userId, create.AddressID) {
		return errors.New("创建订单失败：请选择正确的地址")
	}

	// calculate final price

	return nil
}
