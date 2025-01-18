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

// SkuCreate 创建Sku分类
func SkuCreate(c *gin.Context, create *dto.SkuCreate) {
	if !SpuExists(create.SpuID) {
		response.Failure(c, "请选择正确的商品")
		return
	}

	param := &model.MmSku{
		SpuID:      int32(create.SpuID),
		Name:       create.Name,
		Desc:       create.Desc,
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	tx := util.DBClient().Model(&model.MmSku{}).Create(&param)
	if tx.Error != nil {
		response.Failure(c, "创建失败")
		return
	}

	response.Success(c, []string{})
	return
}

// SkuDelete 删除Sku分类
func SkuDelete(c *gin.Context, delete *dto.SkuDelete) {
	update := &model.MmSku{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmSku{}).Select("status", "delete_time").Where("id = ?", delete.Id).Updates(update)
	if tx.Error != nil {
		response.Failure(c, "删除失败")
		return
	}

	response.Success(c, []string{})
	return
}

// SkuUpdate 修改Sku的分类信息
func SkuUpdate(c *gin.Context, update *dto.SkuUpdate) {
	if !SpuExists(update.SpuID) {
		response.Failure(c, "请选择正确的商品")
		return
	}

	param := &model.MmSku{
		Name:       update.Name,
		Desc:       update.Desc,
		UpdateTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmSku{}).Where("status = ?", constants.NormalStatus).Where("id = ?", update.Id).Updates(param)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, []string{})
	return
}

// SkuSearch 查询
func SkuSearch(c *gin.Context, search *dto.SkuSearch) {
	// id, name, page, limit
	category := &[]dto.SkuSearchResponse{}

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

	if search.SpuID > 0 {
		whereStr += " AND spu_id = ?"
		param = append(param, search.SpuID)
	}

	tx := util.DBClient().Model(&model.MmSku{}).Debug().
		Offset((search.Page-1)*search.Limit).
		Limit(search.Limit).
		Where(whereStr, param...).
		Preload("Specs").
		Find(&category)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, category)
	return
}

// SkuExists 通过ID查询该商品是否可用
func SkuExists(spuID int) bool {
	tx := util.DBClient().Model(&model.MmSku{}).
		Select("id").
		Where("status = ?", constants.NormalStatus).
		Where("id = ?", spuID).
		First(&model.MmCategory{})
	return tx.RowsAffected != 0
}
