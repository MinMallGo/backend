package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
	"math"
	"time"
)

// SpecCreate 创建Spec分类
func SpecCreate(c *gin.Context, create *dto.SpecCreate) {
	if !SkuExists(create.SkuID) {
		response.Failure(c, "请选择正确的规格")
		return
	}

	reasonable := &dto.ReasonableValue{
		Price: create.Price,
		Stock: create.Stock,
	}
	if err := reasonableChecker(reasonable); err != nil {
		response.Failure(c, err.Error())
		return
	}

	param := &model.MmSpec{
		SkuID:      int32(create.SkuID),
		Name:       create.Name,
		Price:      int32(create.Price), // TODO 到时候写前端的时候，注意下这个玩意儿要扩大100倍，然后传整数过来
		Strock:     int32(create.Stock),
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	tx := util.DBClient().Model(&model.MmSpec{}).Create(&param)
	if tx.Error != nil {
		response.Failure(c, "创建失败")
		return
	}

	response.Success(c, []string{})
	return
}

// SpecDelete 删除Spec分类
func SpecDelete(c *gin.Context, delete *dto.SpecDelete) {
	update := &model.MmSpec{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmSpec{}).Select("status", "delete_time").Where("id = ?", delete.Id).Updates(update)
	if tx.Error != nil {
		response.Failure(c, "删除失败")
		return
	}

	response.Success(c, []string{})
	return
}

// SpecUpdate 修改Spec的分类信息
func SpecUpdate(c *gin.Context, update *dto.SpecUpdate) {
	if !SkuExists(update.SkuID) {
		response.Failure(c, "请选择正确的规格")
		return
	}

	reasonable := &dto.ReasonableValue{
		Price: update.Price,
		Stock: update.Stock,
	}
	if err := reasonableChecker(reasonable); err != nil {
		response.Failure(c, err.Error())
		return
	}

	param := &model.MmSpec{
		SkuID:      int32(update.SkuID),
		Name:       update.Name,
		Price:      int32(update.Price),
		Strock:     int32(update.Stock),
		UpdateTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmSpec{}).Where("status = ?", constants.NormalStatus).Where("id = ?", update.Id).Updates(param)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, []string{})
	return
}

// SpecSearch 查询
func SpecSearch(c *gin.Context, search *dto.SpecSearch) {
	// id, name, page, limit
	category := &[]model.MmSpec{}

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

	if search.SkuID > 0 {
		whereStr += " AND sku_id = ?"
		param = append(param, search.SkuID)
	}

	tx := util.DBClient().Model(&model.MmSpec{}).Debug().
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

// 合理值验证
func reasonableChecker(create *dto.ReasonableValue) error {
	if create.Price > math.MaxInt32 {
		return errors.New("请输入合理的价格")
	}

	if create.Stock > math.MaxInt32 {
		return errors.New("请输入合理的库存")
	}

	return nil
}
