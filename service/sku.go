package service

/**
TODO Admin或者商户才有查询和编辑的功能吧。额这些好像都是admin或者商户才能进行的操作
*/

import (
	"github.com/gin-gonic/gin"
	constants "mall_backend/constant"
	"mall_backend/dao"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
	"time"
)

// SkuCreate 创建Sku分类
// sku验证，spec_key验证，value_id验证
// 感觉有点繁琐啊，创建规格名字可以理解，但是为什么要创建了规格值才能传？
func SkuCreate(c *gin.Context, create *dto.SkuCreate) {
	if !SpuExists(create.SpuID) {
		response.Failure(c, "创建失败：请选择正确的商品")
		return
	}

	idx := make([]int, 0, len(create.Spec))
	for _, specs := range create.Spec {
		idx = append(idx, specs.KeyID)
	}

	if len(idx) == 0 || !dao.NewSpecKeyDao().Exists(idx) {
		response.Failure(c, "请选择正确的商品规格名")
		return
	}

	ids := make([]int, 0, len(create.Spec))
	for _, specs := range create.Spec {
		ids = append(ids, specs.ValID)
	}
	if len(ids) == 0 || !dao.NewSpecValueDao().Exists(ids...) {
		response.Failure(c, "创建失败：请选择正确的商品值")
		return
	}

	// 使用事务包裹
	if err := dao.NewSkuDao().Create(create); err != nil {
		response.Failure(c, err.Error())
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
		response.Failure(c, "修改失败：请选择正确的商品")
		return
	}

	if !dao.NewSkuDao().Exists(update.SpuID) {
		response.Failure(c, "修改失败：请选择正确的商品规格")
		return
	}

	// 使用事务包裹
	if err := dao.NewSkuDao().Update(update); err != nil {
		response.Failure(c, err.Error())
		return
	}

	response.Success(c, []string{})
	return
}

// SkuSearch TODO 这里是要设计查询
//
//		  商品
//		 /   \
//		sku1  sku2
//
//	 通过ID 来查询有哪些规格；编辑就需要点到具体的规格里面
func SkuSearch(c *gin.Context, search *dto.SkuSearch) {
	res := &model.MmSku{}
	if err := util.DBClient().Model(&model.MmSku{}).Where("status = ?", constants.NormalStatus).Where("spu_id = ?", search.SpuID).First(res).Error; err != nil {
		response.Failure(c, err.Error())
		return
	}
	response.Success(c, res)
	return
}
