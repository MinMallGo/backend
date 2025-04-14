package dao

import (
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

// `discount_type`  '1 = 满减，2 = 折扣',
const (
	DiscountType int = iota
	RateType
)

// `use_range` '使用范围：1 全平台，2商家 3 类别 4 具体商品 5 某个规格',
const (
	RangeAll = iota
	RangeShop
	RangeCategory
	RangeSpu
	RangeSKU
)

type CouponDao struct {
	db *gorm.DB
	m  model.MmCoupon
}

func NewCouponDao() *CouponDao {
	return &CouponDao{
		db: util.DBClient(),
		m:  model.MmCoupon{},
	}
}

func (c CouponDao) Create(create *dto.CouponCreate) error {
	param := &model.MmCoupon{
		Name:             create.Name,
		DiscountType:     int32(create.DiscountType),
		UseRange:         int32(create.UseRange),
		Threshold:        int32(create.Threshold),
		DiscountPrice:    int32(create.DiscountPrice),
		DiscountRate:     int32(create.DiscountRate),
		Code:             util.CouponCode(),
		ReceiveCount:     int32(create.ReceiveCount),
		ReceiveStartTime: create.ReceiveStartTime,
		ReceiveEndTime:   create.ReceiveEndTime,
		UseStartTime:     create.UseStartTime,
		UseEndTime:       create.UseEndTime,
		TotalCount:       int32(create.TotalCount),
		AssginCount:      0,
		UseCount:         0,
		Status:           true,
		CreateTime:       time.Now().UTC(),
		UpdateTime:       util.MinDateTime(),
		DeleteTime:       util.MinDateTime(),
	}
	if tx := c.db.Model(c.m).Create(param); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (c CouponDao) Delete(delete *dto.CouponDelete) error {
	// TODO 删除的时候检查是否是自己的。不能乱删
	// TODO 这里需要更新用户 user_coupon
	return c.db.Model(c.m).Where("id = ?", delete.Id).Update("status", false).Error
}

func (c CouponDao) Update(update *dto.CouponUpdate) error {
	return c.db.Model(c.m).Where("id = ?", update.ID).Updates(map[string]interface{}{
		"name":               update.Name,
		"discount_type":      update.DiscountType,
		"use_range":          update.UseRange,
		"threshold":          update.Threshold,
		"discount_price":     update.DiscountPrice,
		"discount_rate":      update.DiscountRate,
		"receive_count":      update.ReceiveCount,
		"receive_start_time": update.ReceiveStartTime,
		"receive_end_time":   update.ReceiveEndTime,
		"use_start_time":     update.UseStartTime,
		"use_end_time":       update.UseEndTime,
		"total_count":        update.TotalCount,
	}).Error
}

func (c CouponDao) Search(search *dto.CouponSearch) (dto.PaginateCount, error) {
	data := &[]model.MmCoupon{}

	// 先写吧，优化等以后再来说
	whereStr := "status = ?"
	var param []interface{}
	param = append(param, constants.NormalStatus)

	if len(search.Name) > 0 {
		whereStr += " AND name LIKE ?"
		param = append(param, "%"+search.Name+"%")
	}

	if search.DiscountType > 0 {
		whereStr += " AND discount_type = ?"
		param = append(param, search.DiscountType)
	}

	if !search.ReceiveStartTime.IsZero() {
		whereStr += " AND receive_start_time > ?"
		param = append(param, search.ReceiveStartTime)
	}

	if !search.ReceiveEndTime.IsZero() {
		whereStr += " AND receive_end_time > ?"
		param = append(param, search.ReceiveStartTime)
	}

	if !search.UseStartTime.IsZero() {
		whereStr += " AND use_start_time > ?"
		param = append(param, search.ReceiveStartTime)
	}

	if !search.UseEndTime.IsZero() {
		whereStr += " AND use_end_time > ?"
		param = append(param, search.ReceiveStartTime)
	}

	tx := c.db.Model(c.m).Debug().
		Offset((search.Page-1)*search.Size).
		Limit(search.Size).
		Where(whereStr, param...).
		Find(&data)
	if tx.Error != nil {
		return dto.PaginateCount{}, tx.Error
	}

	var count int64
	tx = c.db.Model(c.m).
		Debug().
		Where(whereStr, param...).
		Count(&count)
	if tx.Error != nil {
		return dto.PaginateCount{}, tx.Error
	}

	res := dto.PaginateCount{
		Data:  data,
		Page:  search.Page,
		Size:  search.Size,
		Count: int(count),
	}

	return res, nil
}

func CouponGetByIds(ids []int) (*[]model.MmCoupon, error) {
	res := &[]model.MmCoupon{}
	err := util.DBClient().Model(&model.MmCoupon{}).Debug().Where("status = ?", constants.NormalStatus).Find(res, ids)
	if err.Error != nil {
		return nil, err.Error
	}

	return res, nil
}

// CouponCreate 创建
func CouponCreate(create *dto.CouponCreate) error {

	return nil
}
