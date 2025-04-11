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
	m  model.MmSpecKey
}

func NewSpecKeyDao() *SpecKeyDao {
	return &SpecKeyDao{
		db: util.DBClient(),
		m:  model.MmSpecKey{},
	}
}

func (d SpecKeyDao) Create(create *dto.SpecKeyCreate) (res dto.SpecKeyCreateResp, err error) {
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
		if err = d.db.Model(&d.m).Create(param).Error; err != nil {
			tx.Rollback()
			return
		}
		spuId = int(param.ID)
	}

	values, err := NewSpecValueDao().Create(int(spuId), create.Value...)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	res.Id = int(spuId)
	res.Values = values

	return
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

func SpecKeyCreate(create *dto.SpecKeyCreate) error {
	// 先查询和判断这个名字是不是存在，存在就提示存在，不存在呢，就新增
	exists := &model.MmSpecKey{}
	err := util.DBClient().Model(&model.MmSpecKey{}).Where("name = ?", create.Name).Find(exists).Error
	if err != nil {
		return err
	}

	// 存在跳过新增

	param := &model.MmSpecKey{
		Name:       create.Name,
		Unit:       create.Uint,
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	err = util.DBClient().Model(&model.MmSpecKey{}).Create(param).Error
	return err
}

func SpecKeyDelete(del *dto.SpecKeyDelete) error {
	param := &model.MmSpecKey{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}

	err := util.DBClient().Debug().Select("status", "delete_time").Model(&model.MmSpecKey{}).Where("id = ?", del.Id).Updates(param).Error
	return err
}

func SpecKeyUpdate(update *dto.SpecKeyUpdate) error {
	param := &model.MmSpecKey{
		Name:       update.Name,
		Unit:       update.Uint,
		UpdateTime: time.Now(),
	}
	err := util.DBClient().Model(&model.MmSpecKey{}).Where("id = ?", update.Id).Updates(param).Error
	return err
}

func SpecKeySearch(search *dto.SpecKeySearch) (*dto.PaginateCount, error) {
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

	param := &[]model.MmSpecKey{}
	err := util.DBClient().Model(&model.MmSpecKey{}).Debug().Offset((search.Page-1)*search.Size).Limit(search.Size).Where(whereStr, params).Find(param).Error
	if err != nil {
		return nil, err
	}
	var count int64 = 0
	err = util.DBClient().Model(&model.MmSpecKey{}).Debug().Where(whereStr, params).Count(&count).Error
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

// SpecKeyExists 判断是否存在
func SpecKeyExists(specId int) bool {
	e := util.DBClient().Model(&model.MmSpecKey{}).Debug().Where("status = ?", constants.NormalStatus).Where("id = ?", specId).Find(&model.MmSpecKey{}).RowsAffected
	return e == 1
}

// SpecKeyExists2 SpecKeyExists 判断全部是否存在
func SpecKeyExists2(specId []int) bool {
	return util.DBClient().Model(&model.MmSpecKey{}).
		Debug().
		Where("status = ?", constants.NormalStatus).
		Where("id = ?", specId).Find(&model.MmSpecKey{}).RowsAffected == int64(len(specId))
}
