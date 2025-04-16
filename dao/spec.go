package dao

import (
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

// SpecCreate 创建spec，保存关联关系
func SpecCreate(spuId, skuId int, create *[]dto.Specs) error {
	param := &[]model.MmSpec{}
	for _, specs := range *create {
		*param = append(*param, model.MmSpec{
			SpuID:      int32(spuId),
			SkuID:      int32(skuId),
			KeyID:      int32(specs.KeyID),
			ValueID:    int32(specs.ValID),
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		})
	}

	return util.DBClient().Create(param).Error
}

// SpecUpdate 通过商品ID以及规格ID进行编辑
func SpecUpdate(spuId, skuId int, create *[]dto.Specs) (err error) {
	param := &[]model.MmSpec{}
	for _, specs := range *create {
		*param = append(*param, model.MmSpec{
			SpuID:      int32(spuId),
			SkuID:      int32(skuId),
			KeyID:      int32(specs.KeyID),
			ValueID:    int32(specs.ValID),
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		})
	}

	err = util.DBClient().Model(&model.MmSpec{}).Where("spu_id = ?", spuId).Where("sku_id = ?", skuId).Delete(&model.MmSpec{}).Error
	if err != nil {
		return
	}
	return util.DBClient().Model(&model.MmSpec{}).Create(param).Error
}
