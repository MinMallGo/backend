package dao

import (
	"encoding/json"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
)

type SpecKey struct {
	model.MmSpecKey
	Value []model.MmSpecValue ` json:"value" gorm:"foreignKey:KeyID;references:ID"`
}

// GenSkuData 通过spec_key_id,spec_id 来查询和构造 sku.Title 以及sku.Specs
func GenSkuData(spec *[]dto.Specs) (title string, specs []byte, err error) {
	// 通过一条sql关联查询出来 spec_key => []spec_value然后进行筛选
	// 先for range 获取出来所有的spec_key_id,然后获取所有的 for range res进行筛选
	var keyIds, valueIds []int

	for _, spec := range *spec {
		keyIds = append(keyIds, spec.KeyID)
		valueIds = append(valueIds, spec.ValID)
	}
	res := &[]SpecKey{}
	err = util.DBClient().Model(&model.MmSpecKey{}).Debug().Preload("Value", func(db *gorm.DB) *gorm.DB {
		return db.Find(&model.MmSpecValue{}, valueIds)
	}).Find(res, keyIds).Error
	if err != nil {
		return
	}

	// for range 构造title
	for _, v := range *res {
		for _, v2 := range v.Value {
			title += v2.Name
		}
	}

	specs, err = json.Marshal(res)
	return
}
