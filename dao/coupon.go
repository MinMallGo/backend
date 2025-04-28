package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"strings"
	"time"
)

// `discount_type`  '1 = 满减，2 = 折扣',
const (
	DiscountType int = 1 + iota
	RateType
)

const (
	ThresholdAmount = 1
	ThresholdUnit   = 2
)

const (
	CouponPlatform = "PLATFORM"
	CouponMerchant = "MERCHANT"
	CouponCategory = "CATEGORY"
	CouponProduct  = "PRODUCT"
)

type CouponDao struct {
	db *gorm.DB
}

func NewCouponDao(db *gorm.DB) *CouponDao {
	return &CouponDao{
		db: db,
	}
}

func (c *CouponDao) Create(create *dto.CouponCreate) (int, error) {
	param := &model.MmCoupon{
		Name:             create.Name,
		Code:             util.CouponCode(),
		UsageScopeID:     int32(create.ScopeID),
		ThresholdType:    int32(create.ThresholdType),
		ThresholdValue:   int32(create.ThresholdValue),
		DiscountType:     int32(create.DiscountType),
		DiscountValue:    int32(create.DiscountValue),
		UserReceiveLimit: create.ReceiveLimit == 0,
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
	}
	if tx := c.db.Model(&model.MmCoupon{}).Create(param); tx.Error != nil {
		return 0, tx.Error
	}
	return int(param.ID), nil
}

func (c *CouponDao) Delete(delete *dto.CouponDelete) error {
	// TODO 删除的时候检查是否是自己的。不能乱删
	// TODO 这里需要更新用户 user_coupon
	return c.db.Model(&model.MmCoupon{}).Where("id = ?", delete.Id).Update("status", false).Error
}

func (c *CouponDao) Update(update *dto.CouponUpdate) error {
	return c.db.Model(&model.MmCoupon{}).Where("id = ?", update.ID).Updates(model.MmCoupon{
		Name:             update.Name,
		Code:             util.CouponCode(),
		UsageScopeID:     int32(update.ScopeID),
		ThresholdType:    int32(update.ThresholdType),
		ThresholdValue:   int32(update.ThresholdValue),
		DiscountType:     int32(update.DiscountType),
		DiscountValue:    int32(update.DiscountValue),
		UserReceiveLimit: update.ReceiveLimit == 0,
		ReceiveStartTime: update.ReceiveStartTime,
		ReceiveEndTime:   update.ReceiveEndTime,
		UseStartTime:     update.UseStartTime,
		UseEndTime:       update.UseEndTime,
		TotalCount:       int32(update.TotalCount),
		UpdateTime:       time.Now().UTC(),
	}).Error
}

func (c *CouponDao) Search(search *dto.CouponSearch) (dto.PaginateCount, error) {
	data := &[]CouponLadders{}

	// 先写吧，优化等以后再来说
	whereStr := "status = ?"
	var param []interface{}
	param = append(param, constants.NormalStatus)

	if len(search.Name) > 0 {
		whereStr += " AND name LIKE ?"
		param = append(param, "%"+search.Name+"%")
	}

	if search.ThresholdType > 0 {
		whereStr += " AND threshold_type = ?"
		param = append(param, search.ThresholdType)
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
		Preload("Ladders", "status = ?", constants.NormalStatus).
		Preload("Scope").
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

type CouponWithScope struct {
	model.MmCoupon
	Scope model.MmCouponUsageScope `gorm:"foreignKey:UsageScopeID;references:ID" json:"scope"`
}

func (c *CouponDao) ExistsWithUser(userId int, ids ...int) error {
	// 检查优惠券是否可用
	coupons := &[]CouponWithScope{}
	res := c.db.Model(&model.MmCoupon{}).
		Preload("Scope").
		Where("status", true).
		Where("id in ?", ids).
		Where("total_count > use_count").         // 没有使用完才能用
		Where("use_start_time <= ?", time.Now()). // 检查优惠券是否过期
		Where("use_end_time >= ?", time.Now()).   // 检查优惠券是否过期
		Find(coupons)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected < int64(len(ids)) {
		return errors.New("请选择正确的优惠券")
	}

	distinctCoupon := make(map[int]int)
	distinctScope := make(map[string]int)
	for _, coupon := range *coupons {
		distinctCoupon[int(coupon.DiscountType)]++
		distinctScope[strings.ToUpper(coupon.Scope.Code)]++
	}
	// 只有在针对打折券，才是多张不能一起用，
	if distinctCoupon[RateType] > 1 {
		return errors.New("打折券不能叠加使用")
	}

	// 这里限制一下使用规则吧，最多只能使用两张券// 平台可以与任何券叠加。但是其他的不能两两叠加。
	dsLen := len(distinctScope)
	// 最多只能使用两张优惠券，而且必须有一张是平台券才行，否则不能一起用

	if _, ok := distinctScope[CouponPlatform]; !ok && dsLen >= 2 {
		return errors.New("只有平台券才能和其他券叠加使用")
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

func (c *CouponDao) Cancel(ids ...int) error {
	if len(ids) == 0 {
		return nil
	}
	str := "WHEN id = %d THEN use_count - %d "
	when := ""
	whereStr := " id in ( %s ) "
	where := ""
	for index, item := range ids {
		when += fmt.Sprintf(str, item, 1)
		where += fmt.Sprintf(" '%d' ", item)
		if index < len(ids)-1 {
			where += ","
		}
	}

	sql := `UPDATE mm_coupon SET use_count = CASE %s END WHERE %s`
	sql = fmt.Sprintf(sql, when, fmt.Sprintf(whereStr, where))
	//res := &[]model.MmSku{}

	tx := c.db.Exec(sql)
	if tx.Error != nil || tx.RowsAffected < int64(len(ids)) {
		return errors.New("归还优惠券失败")
	}
	return nil
}

func (c *CouponDao) Exists(ids ...int) error {
	if len(ids) == 0 {
		return nil
	}
	res := c.db.Model(&model.MmCoupon{}).
		Where("status", true).
		Where("id in ?", ids).
		Distinct("discount_type").
		Find(&[]model.MmCoupon{})
	if res.Error != nil || res.RowsAffected < int64(len(ids)) {
		return res.Error
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

type CouponLadders struct {
	model.MmCoupon
	Ladders []model.MmCouponLadderRule `gorm:"foreignKey:CouponID;references:ID" json:"ladders"`
	Scope   model.MmCouponUsageScope   `gorm:"foreignKey:UsageScopeID;references:ID" json:"scope"`
}

func (c *CouponDao) CouponWithLadder(ids ...int) (*[]CouponLadders, error) {
	res := &[]CouponLadders{}
	tx := c.db.Model(&model.MmCoupon{}).
		Where("id in ? AND status = ?", ids, constants.NormalStatus).
		Preload("Ladders", "status = ?", constants.NormalStatus).
		Preload("Scope").
		Order("usage_scope_id desc").
		Find(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return res, nil
}
