package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dto"
	"mall_backend/util"
)

type OrderDao struct {
	db *gorm.DB
}

func NewOrderDao() *OrderDao {
	return &OrderDao{
		db: util.DBClient(),
	}
}

func (u *OrderDao) Create(create *dto.OrderCreate, userId int) error {
	// 1. 事务启动
	// 2. 检查商品，检查规格
	// 3. 地址检查
	// 4. 验证优惠券并计算价格
	// 5. 构造数据
	spuIds := make([]int, 0, len(create.Product))
	skuIds := make([]int, 0, len(create.Product))
	mSku := make(map[int]int, len(create.Product))
	for _, v := range create.Product {
		spuIds = append(spuIds, v.SpuID)
		skuIds = append(skuIds, v.SkuID)
		mSku[v.SkuID] = v.Num
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
	if !NewCouponDao().ExistsWithUser(userId, create.Coupons...) {
		return errors.New("创建订单失败：请选择正确的优惠券")
	}
	// 通过查询获取到商品的价格
	product, err := NewSkuDao().MoreWithOrder(&create.Product)
	if err != nil {
		return errors.New("创建订单失败：优惠券检查失败")
	}

	totalPrice := 0
	for _, item := range *product {
		totalPrice += int(item.Price) * mSku[int(item.ID)]
	}

	coupons, err := NewCouponDao().More(create.Coupons...)
	if err != nil {
		return errors.New("创建订单失败：查询优惠券失败")
	}

	// TODO 检查商户优惠券是否能用。 这里商户还没写欸，先不加

	for _, item := range *coupons {

	}

	// calculate final price

	return nil
}
