package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
)

// 定义优惠券的类别。需要注意的是这里以后需要手动维护，如果数据库新增了
const (
	PlatformScope string = `platform`
	StoreScope    string = `store`
	CategoryScope string = `category`
	ProductScope  string = `product`
)

type CouponUsageScopeDao struct {
	db *gorm.DB
}

func NewCouponUsageScopeDao(db *gorm.DB) *CouponUsageScopeDao {
	return &CouponUsageScopeDao{db: db}
}

// Exists 检查是否存在
func (s *CouponUsageScopeDao) Exists(id ...int) bool {
	if len(id) == 0 {
		return false
	}
	var count int64
	tx := s.db.Model(&model.MmCouponUsageScope{}).Where("id in ?", id).Count(&count)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return false
	}

	return count == int64(len(id))
}
