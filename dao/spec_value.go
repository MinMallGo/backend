package dao

import (
	"errors"
	"fmt"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

func SpecValueCreate(create *dto.SpecValueCreate) error {
	// 先查询和判断这个名字是不是存在，存在就提示存在，不存在呢，就新增
	if !SpecKeyExists(create.SpecKeyId) {
		return errors.New("请选择正确的规格名称编号")
	}

	exists := &model.MmSpecValue{}
	err := util.DBClient().Model(&model.MmSpecValue{}).Where("id = ?", create.SpecKeyId).Where("name = ?", create.Name).Find(exists).Error
	if err != nil {
		return err
	}

	if exists.ID != 0 {
		return errors.New(fmt.Sprintf("要新增的名称已经存在：%s", create.Name))
	}

	param := &model.MmSpecValue{
		Name:       create.Name,
		KeyID:      int32(create.SpecKeyId),
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	err = util.DBClient().Model(&model.MmSpecValue{}).Create(param).Error
	return err
}

func SpecValueDelete(del *dto.SpecValueDelete) error {
	param := &model.MmSpecValue{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}

	err := util.DBClient().Debug().Select("status", "delete_time").Model(&model.MmSpecValue{}).Where("id = ?", del.Id).Updates(param).Error
	return err
}

func SpecValueUpdate(update *dto.SpecValueUpdate) error {
	param := &model.MmSpecValue{
		Name:       update.Name,
		KeyID:      int32(update.SpecKeyId),
		UpdateTime: time.Now(),
	}
	err := util.DBClient().Model(&model.MmSpecValue{}).Where("id = ?", update.Id).Updates(param).Error
	return err
}

func SpecValueSearch(search *dto.SpecValueSearch) (*dto.PaginateCount, error) {
	whereStr := "status = ?"
	var params []interface{}
	params = append(params, constants.NormalStatus)

	if len(search.Name) > 0 {
		whereStr += "AND name LIKE ?"
		params = append(params, "%"+search.Name+"%")
	}

	if search.SpecKeyId > 0 {
		whereStr += "AND key_id = ?"
		params = append(params, search.SpecKeyId)
	}

	param := &[]model.MmSpecValue{}
	err := util.DBClient().Model(&model.MmSpecValue{}).Debug().Offset((search.Page-1)*search.Size).Limit(search.Size).Where(whereStr, params).Find(param).Error
	if err != nil {
		return nil, err
	}

	var count int64 = 0
	err = util.DBClient().Model(&model.MmSpecValue{}).Debug().Where(whereStr, params).Count(&count).Error
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
