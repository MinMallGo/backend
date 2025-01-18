package service

import (
	"github.com/gin-gonic/gin"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
	"strings"
	"time"
)

// SpuCreate 创建spu分类
func SpuCreate(c *gin.Context, create *dto.SpuCreate) {
	if !CategoryExists(create.CategoryId) {
		response.Failure(c, "请选择正确的分类")
		return
	}

	if create.MerchantId != 0 && !MerchantExists(create.MerchantId) {
		response.Failure(c, "请选择正确的商户")
		return
	}

	images := ""
	if len(create.Img) > 0 {
		images += strings.Join(create.Img, ",")
	}

	param := &model.MmSpu{
		Name:       create.Name,
		Desc:       create.Desc,
		Imgs:       images,
		Code:       util.SpuCode(),
		BrandID:    int32(create.BrandId), // TODO 需要验证和绑定品牌
		CategoryID: int32(create.CategoryId),
		MerchantID: int32(create.MerchantId),
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	tx := util.DBClient().Model(&model.MmSpu{}).Create(&param)
	if tx.Error != nil {
		response.Failure(c, "创建失败")
		return
	}

	response.Success(c, []string{})
	return
}

// SpuDelete 删除spu分类
func SpuDelete(c *gin.Context, delete *dto.SpuDelete) {
	update := &model.MmSpu{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmSpu{}).Select("status", "delete_time").Where("id = ?", delete.Id).Updates(update)
	if tx.Error != nil {
		response.Failure(c, "删除失败")
		return
	}

	response.Success(c, []string{})
	return
}

// SpuUpdate 修改spu的分类信息
func SpuUpdate(c *gin.Context, update *dto.SpuUpdate) {
	if !CategoryExists(update.CategoryId) {
		response.Failure(c, "请选择正确的分类")
		return
	}

	if update.CategoryId != 0 && !MerchantExists(update.MerchantId) {
		response.Failure(c, "请选择正确的商户")
		return
	}

	images := ""
	if len(update.Img) > 0 {
		images += strings.Join(update.Img, ",")
	}

	param := &model.MmSpu{
		Name:       update.Name,
		Desc:       update.Desc,
		Imgs:       images,
		BrandID:    int32(update.BrandId),
		CategoryID: int32(update.CategoryId),
		MerchantID: int32(update.MerchantId),
		UpdateTime: time.Now(),
	}

	tx := util.DBClient().Model(&model.MmSpu{}).Where("status = ?", constants.NormalStatus).Where("id = ?", update.Id).Updates(param)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, []string{})
	return
}

// SpuSearch 查询
func SpuSearch(c *gin.Context, search *dto.SpuSearch) {
	// id, name, page, limit
	category := &[]dto.SpuSearchResponse{}

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

	if search.MerchantId != 0 {
		whereStr += " AND merchant_id = ?"
		param = append(param, search.MerchantId)
	}

	if search.CategoryId != 0 {
		whereStr += " AND category_id = ?"
		param = append(param, search.CategoryId)
	}

	tx := util.DBClient().Model(&model.MmSpu{}).Debug().
		Offset((search.Page-1)*search.Limit).
		Limit(search.Limit).
		Where(whereStr, param...).
		Preload("Category").
		Find(&category)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, category)
	return
}

// SpuExists 通过ID查询该商品是否可用
func SpuExists(spuID int) bool {
	tx := util.DBClient().Model(&model.MmSpu{}).
		Select("id").
		Where("status = ?", constants.NormalStatus).
		Where("id = ?", spuID).
		First(&model.MmCategory{})
	return tx.RowsAffected != 0
}
