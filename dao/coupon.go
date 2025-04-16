package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

// `discount_type`  '1 = 满减，2 = 折扣',
const (
	DiscountType int = 1 + iota
	RateType
)

// `use_range` '使用范围：1 全平台，2商家 3 类别 4 具体商品 5 某个规格',
const (
	RangeAll = 1 + iota
	RangeMerchant
	RangeCategory
	RangeSpu
	RangeSKU
)

type CouponDao struct {
	db *gorm.DB
}

func NewCouponDao(db *gorm.DB) *CouponDao {
	return &CouponDao{
		db: db,
	}
}

func (c *CouponDao) Create(create *dto.CouponCreate) error {
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
	if tx := c.db.Model(&model.MmCoupon{}).Create(param); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (c *CouponDao) Delete(delete *dto.CouponDelete) error {
	// TODO 删除的时候检查是否是自己的。不能乱删
	// TODO 这里需要更新用户 user_coupon
	return c.db.Model(&model.MmCoupon{}).Where("id = ?", delete.Id).Update("status", false).Error
}

func (c *CouponDao) Update(update *dto.CouponUpdate) error {
	return c.db.Model(&model.MmCoupon{}).Where("id = ?", update.ID).Updates(map[string]interface{}{
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

func (c *CouponDao) Search(search *dto.CouponSearch) (dto.PaginateCount, error) {
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
		whereStr += " AND receive_start_time >= ?"
		param = append(param, search.ReceiveStartTime)
	}

	if !search.ReceiveEndTime.IsZero() {
		whereStr += " AND receive_end_time <= ?"
		param = append(param, search.ReceiveEndTime)
	}

	if !search.UseStartTime.IsZero() {
		whereStr += " AND use_start_time >= ?"
		param = append(param, search.UseStartTime)
	}

	if !search.UseEndTime.IsZero() {
		whereStr += " AND use_end_time <= ?"
		param = append(param, search.UseEndTime)
	}

	tx := c.db.Model(&model.MmCoupon{}).
		Offset((search.Page-1)*search.Size).
		Limit(search.Size).
		Where(whereStr, param...).
		Find(&data)
	if tx.Error != nil {
		return dto.PaginateCount{}, tx.Error
	}

	var count int64
	tx = c.db.Model(&model.MmCoupon{}).
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

func (c *CouponDao) ExistsWithUser(userId int, ids ...int) error {
	// 检查优惠券是否可用
	coupons := &[]model.MmCoupon{}
	res := c.db.Model(&model.MmCoupon{}).
		Where("status", true).
		Where("id in ?", ids).
		Where("total_count > use_count"). // 没有使用完才能用
		Distinct("discount_type").
		Find(coupons)
	if res.Error != nil || res.RowsAffected < int64(len(ids)) {
		return res.Error
	}

	distinctCoupon := make(map[int]struct{})
	for _, coupon := range *coupons {
		distinctCoupon[int(coupon.DiscountType)] = struct{}{}

	}
	if len(distinctCoupon) < len(ids) {
		return errors.New("同类券不能多次叠加使用")
	}

	// 检查用户优惠券
	var count int64
	res = c.db.Model(&model.MmUserCoupon{}).
		Where("id in ?", ids).
		Where("user_id = ?", userId).
		Where("is_used = ?", false).
		Where("status = ?", 1).
		Count(&count)
	if res.Error != nil || (count) < int64(len(ids)) {
		return res.Error
	}
	return nil
}

func (c *CouponDao) More(ids ...int) (res *[]model.MmCoupon, err error) {
	err = c.db.Model(&model.MmCoupon{}).Where("id in ?", ids).Where("status = ?", true).Find(&res).Error
	if err != nil {
		return
	}
	return
}

type CouponUpdate struct {
	ID  int
	Num int
}

func (c *CouponDao) CouponUse(u ...int) error {
	// update xx set stock = case when id = 1 then stock - 1 when id = 2 then stock - 2 end where id in (1,2)
	if len(u) == 0 {
		return nil
	}
	str := "WHEN id = %d THEN use_count + %d "
	when := ""
	orStr := "(id = %d AND total_count >= use_count + %d )"
	where := ""
	for index, item := range u {
		when += fmt.Sprintf(str, item, 1)
		where += fmt.Sprintf(orStr, item, 1)
		if index < len(u)-1 {
			where += " OR "
		}
	}

	sql := `UPDATE mm_coupon SET use_count = CASE %s END WHERE %s`
	sql = fmt.Sprintf(sql, when, where)
	//res := &[]model.MmSku{}

	tx := c.db.Exec(sql)
	if tx.Error != nil || tx.RowsAffected < int64(len(u)) {
		return errors.New("使用优惠券失败")
	}
	return nil
}

func CouponGetByIds(ids []int) (*[]model.MmCoupon, error) {
	res := &[]model.MmCoupon{}
	err := util.DBClient().Model(&model.MmCoupon{}).Where("status = ?", constants.NormalStatus).Find(res, ids)
	if err.Error != nil {
		return nil, err.Error
	}

	return res, nil
}
