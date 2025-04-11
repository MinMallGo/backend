package dao

import (
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

type SpecValueDao struct {
	db *gorm.DB
	m  model.MmSpecValue
}

func NewSpecValueDao() *SpecValueDao {
	return &SpecValueDao{
		db: util.DBClient(),
		m:  model.MmSpecValue{},
	}
}

// Create 这里创建需要等待到获取spu的id然后循环创建
func (d SpecValueDao) Create(specId int, create ...string) ([]dto.SpecKeyCreateRespValue, error) {
	// 先查询有无，有就直接拿到id，无则创建拿到id
	result := make([]dto.SpecKeyCreateRespValue, 0, len(create))
	for _, name := range create {
		tmp := model.MmSpecValue{}
		if d.db.Model(d.m).Where("name = ?", name).Select("id").Find(&tmp).RowsAffected != 0 {
			// 获取到id
			result = append(result, dto.SpecKeyCreateRespValue{
				Name: name,
				Id:   int(tmp.ID),
			})
		} else {
			param := &model.MmSpecValue{
				KeyID:      int32(specId),
				Name:       name,
				CreateTime: time.Now(),
				UpdateTime: util.MinDateTime(),
				DeleteTime: util.MinDateTime(),
			}
			res := d.db.Create(param)
			if res.Error != nil {
				return result, res.Error
			}
			result = append(result, dto.SpecKeyCreateRespValue{
				Name: name,
				Id:   int(param.ID),
			})
		}
	}
	return result, nil
}

//	func (d SpecValueDao) Update(update *dto.SpecValueUpdate) error {
//		param := &model.MmSpecValue{
//			Name:       update.Name,
//			KeyID:      int32(update.SpecKeyId),
//			UpdateTime: time.Now(),
//		}
//		err := util.DBClient().Model(&model.MmSpecValue{}).Where("id = ?", update.Id).Updates(param).Error
//		return err
//	}
//
//	func (d SpecValueDao) Delete(del *dto.SpecValueDelete) error {
//		param := &model.MmSpecValue{
//			Status:     constants.BanStatus,
//			DeleteTime: time.Now(),
//		}
//
//		return util.DBClient().Debug().Select("status", "delete_time").Model(&model.MmSpecValue{}).Where("id = ?", del.Id).Updates(param).Error
//	}
//
// func (d SpecValueDao) Get() {
//
// }
//
// func (d SpecValueDao) Gets(id ...int) {
//
// }
func (d SpecValueDao) Exists(id ...int) bool {
	return util.DBClient().Model(&model.MmSpecValue{}).
		Where("status = ?", constants.NormalStatus).
		Where("id in ?", id).
		Find(&model.MmSpecValue{}).RowsAffected == int64(len(id))
}

//
//func (d SpecValueDao) Paginate(search *dto.SpecValueSearch) (*dto.PaginateCount, error) {
//	whereStr := "status = ?"
//	var params []interface{}
//	params = append(params, constants.NormalStatus)
//
//	if len(search.Name) > 0 {
//		whereStr += "AND name LIKE ?"
//		params = append(params, "%"+search.Name+"%")
//	}
//
//	if search.SpecKeyId > 0 {
//		whereStr += "AND key_id = ?"
//		params = append(params, search.SpecKeyId)
//	}
//
//	param := &[]model.MmSpecValue{}
//	err := util.DBClient().Model(&model.MmSpecValue{}).Debug().Offset((search.Page-1)*search.Size).Limit(search.Size).Where(whereStr, params).Find(param).Error
//	if err != nil {
//		return nil, err
//	}
//
//	var count int64 = 0
//	err = util.DBClient().Model(&model.MmSpecValue{}).Debug().Where(whereStr, params).Count(&count).Error
//	if err != nil {
//		return nil, err
//	}
//	res := &dto.PaginateCount{
//		Data:  param,
//		Page:  search.Page,
//		Size:  search.Size,
//		Count: int(count),
//	}
//	return res, err
//}
