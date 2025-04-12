package service

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
)

// SpuCreate 创建spu分类
func SpuCreate(c *gin.Context, create *dto.SpuCreate) {
	if !CategoryExists(create.CategoryId) {
		response.Failure(c, "请选择正确的分类")
		return
	}

	// TODO 商户
	//if create.MerchantId != 0 && !MerchantExists(create.MerchantId) {
	//	response.Failure(c, "请选择正确的商户")
	//	return
	//}

	err := dao.NewSpuDao().Create(create)
	if err != nil {
		response.Failure(c, "创建失败")
		return
	}

	response.Success(c, []string{})
	return
}

// SpuDelete 删除spu分类
func SpuDelete(c *gin.Context, delete *dto.SpuDelete) {
	err := dao.NewSpuDao().Delete(delete)
	if err != nil {
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

	//if update.CategoryId != 0 && !MerchantExists(update.MerchantId) {
	//	response.Failure(c, "请选择正确的商户")
	//	return
	//}

	err := dao.NewSpuDao().Update(update)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, []string{})
	return
}

// SpuSearch 查询
func SpuSearch(c *gin.Context, search *dto.SpuSearch) {
	res, err := dao.NewSpuDao().Search(search)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}

	response.Success(c, res)
	return
}
