package dao

import (
	"encoding/json"
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

type SpecKey struct {
	model.MmSpecKey
	Value []model.MmSpecValue ` json:"value" gorm:"foreignKey:KeyID;references:ID"`
}

type SpecKeyDao struct {
	db *gorm.DB
}

func NewSpecKeyDao(db *gorm.DB) *SpecKeyDao {
	return &SpecKeyDao{
		db: db,
	}
}

func (d *SpecKeyDao) Create(create *dto.SpecKeyCreate) (res dto.SpecKeyCreateResp, err error) {
	// 先查询和判断这个名字是不是存在，存在就提示存在，不存在呢，就新增
	exists := &model.MmSpecKey{}
	result := d.db.Model(&model.MmSpecKey{}).Where("name = ?", create.Name).Find(exists)
	if result.Error != nil {
		return
	}

	tx := d.db.Begin()
	var spuId int
	if result.RowsAffected > 0 {
		spuId = int(exists.ID)
	} else {
		param := &model.MmSpecKey{
			Name:       create.Name,
			Unit:       create.Uint,
			Status:     constants.NormalStatus,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		}
		if err = d.db.Model(&model.MmSpecKey{}).Create(param).Error; err != nil {
			tx.Rollback()
			return
		}
		spuId = int(param.ID)
	}

	values, err := NewSpecValueDao(d.db).Create(int(spuId), create.Value...)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	res.Id = int(spuId)
	res.Values = values

	return
}

func (d *SpecKeyDao) Update(update *dto.SpecKeyUpdate) (err error) {
	return d.db.Model(&model.MmSpecKey{}).
		Where("id = ?", update.Id).Updates(map[string]interface{}{
		"name":        update.Name,
		"unit":        update.Uint,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	}).Error
}

func (d *SpecKeyDao) Delete(delete *dto.SpecKeyDelete) (err error) {
	return d.db.Model(&model.MmSpecKey{}).Where("id = ?", delete.Id).Updates(map[string]interface{}{
		"status":      constants.BanStatus,
		"delete_time": time.Now().Format("2006-01-02 15:04:05"),
	}).Error
}

func (d *SpecKeyDao) OneById(keyId int) (res *SpecKey, err error) {
	err = d.db.Model(model.MmSpecKey{}).Where("id = ?", keyId).Preload("Value").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *SpecKeyDao) More(search *dto.SpecKeySearch) (*dto.PaginateCount, error) {
	whereStr := "status = ?"
	var params []interface{}
	params = append(params, constants.NormalStatus)

	if len(search.Name) > 0 {
		whereStr += "AND name LIKE ?"
		params = append(params, "%"+search.Name+"%")
	}

	if len(search.Uint) > 0 {
		whereStr += "AND unit = ?"
		params = append(params, search.Uint)
	}

	param := &[]SpecKey{}
	err := d.db.Model(&model.MmSpecKey{}).Offset((search.Page-1)*search.Size).Preload("Value").Limit(search.Size).Where(whereStr, params).Find(param).Error
	if err != nil {
		return nil, err
	}
	var count int64 = 0
	err = d.db.Model(&model.MmSpecKey{}).Where(whereStr, params).Count(&count).Error
	if err != nil {
		return nil, err
	}
	res := &dto.PaginateCount{
		Data:  param,
		Page:  search.Page,
		Size:  search.Size,
		Count: int(count),
	}
	return res, err
}

func (d *SpecKeyDao) Exists(specId ...int) bool {
	if len(specId) == 0 {
		return false
	}

	return d.db.
		Model(&model.MmSpecKey{}).
		Where("status = ?", constants.NormalStatus).
		Where("id in ?", specId).
		Find(&[]model.MmSpecKey{}).
		RowsAffected == int64(len(specId))

}

func (d *SpecKeyDao) GenSkuData(spec *[]dto.Specs) (title string, specs []byte, err error) {
	// 通过一条sql关联查询出来 spec_key => []spec_value然后进行筛选
	// 先for range 获取出来所有的spec_key_id,然后获取所有的 for range res进行筛选
	var keyIds, valueIds []int

	for _, spec := range *spec {
		keyIds = append(keyIds, spec.KeyID)
		valueIds = append(valueIds, spec.ValID)
	}
	res := &[]SpecKey{}
	err = d.db.Model(&model.MmSpecKey{}).Preload("Value", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,key_id,name").Find(&model.MmSpecValue{}, valueIds)
	}).Select("id,name,unit").Find(res, keyIds).Error
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
