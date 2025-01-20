package service

import (
	"github.com/gin-gonic/gin"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
	"time"
)

// CouponCreate 创建Coupon分类
func CouponCreate(c *gin.Context, create *dto.CouponCreate) {
	startTime, err := time.Parse("2006-01-02 15:04:05", create.StartTime)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", create.StartTime)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}

	// 如果是折扣则需要判断折扣的范围，必须是1-100之间
	if create.Type == constants.CouponDiscount && (create.Price <= 0 || create.Price > 100) {
		response.Failure(c, "请输入正确的折扣")
		return
	}

	param := &model.MmCoupon{
		Name:       create.Name,
		Type:       int32(create.Type),
		MinSpend:   int32(create.MinSpend),
		Price:      int32(create.Price),
		Code:       util.CouponCode(),
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	tx := util.DBClient().Model(&model.MmCoupon{}).Create(&param)
	if tx.Error != nil {
		response.Failure(c, "创建失败")
		return
	}

	response.Success(c, []string{})
	return
}

// CouponDelete 删除Coupon分类
func CouponDelete(c *gin.Context, delete *dto.CouponDelete) {
	update := &model.MmCoupon{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmCoupon{}).Select("status", "delete_time").Where("id = ?", delete.Id).Updates(update)
	if tx.Error != nil {
		response.Failure(c, "删除失败")
		return
	}

	response.Success(c, []string{})
	return
}

// CouponUpdate 修改Coupon的分类信息
func CouponUpdate(c *gin.Context, update *dto.CouponUpdate) {
	startTime, err := time.Parse("2006-01-02 15:04:05", update.StartTime)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", update.EndTime)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}

	param := &model.MmCoupon{
		Name:       update.Name,
		Type:       int32(update.Type),
		MinSpend:   int32(update.MinSpend),
		Price:      int32(update.Price),
		Code:       util.CouponCode(),
		StartTime:  startTime,
		EndTime:    endTime,
		UpdateTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmCoupon{}).Where("status = ?", constants.NormalStatus).Where("id = ?", update.Id).Updates(param)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, []string{})
	return
}

// CouponSearch 查询
func CouponSearch(c *gin.Context, search *dto.CouponSearch) {
	// id, name, page, limit
	category := &[]model.MmCoupon{}

	// 先写吧，优化等以后再来说
	whereStr := "status = ?"
	var param []interface{}
	param = append(param, constants.NormalStatus)

	if search.Id > 0 {
		whereStr += " AND id = ?"
		param = append(param, search.Id)
	}

	if len(search.Name) > 0 {
		whereStr += " AND name LIKE ?"
		param = append(param, "%"+search.Name+"%")
	}

	if search.Type > 0 {
		whereStr += " AND type = ?"
		param = append(param, search.Type)
	}

	if len(search.StartTime) > 0 {
		startTime, err := time.Parse("2006-01-02 15:04:05", search.StartTime)
		if err != nil {
			response.Failure(c, err.Error())
			return
		}
		whereStr += " AND create_time > ?"
		param = append(param, startTime)
	}

	if len(search.EndTime) > 0 {
		EndTime, err := time.Parse("2006-01-02 15:04:05", search.EndTime)
		if err != nil {
			response.Failure(c, err.Error())
			return
		}
		whereStr += " AND end_time < ?"
		param = append(param, EndTime)
	}

	tx := util.DBClient().Model(&model.MmCoupon{}).Debug().
		Offset((search.Page-1)*search.Limit).
		Limit(search.Limit).
		Where(whereStr, param...).
		Find(&category)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, category)
	return
}
