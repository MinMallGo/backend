package dao

import (
	"errors"
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

type AddrDao struct {
	db *gorm.DB
}

func NewAddrDao(db *gorm.DB) *AddrDao {
	return &AddrDao{
		db: db,
	}
}

func (a *AddrDao) ExistsWithUser(userId int, id ...int) bool {
	tmp := &[]model.MmAddress{}
	return a.db.Model(&model.MmAddress{}).
		Where("id in ?", id).
		Where("user_id = ?", userId).
		Where("status = ?", constants.NormalStatus).
		Find(tmp).RowsAffected == int64(len(id))
}

func (a *AddrDao) Create(create *dto.AddressCreate, userId int) error {
	param := &model.MmAddress{
		UserID:      int32(userId),
		Name:        create.Name,
		Phone:       create.Phone,
		Address:     create.Address,
		HouseNumber: create.HouseNumber,
		IsDefault:   create.IsDefault == 1,
		Tag:         create.Tag,
		Status:      true,
		CreateTime:  time.Now(),
		UpdateTime:  util.MinDateTime(),
		DeleteTime:  util.MinDateTime(),
	}
	tx := a.db.Model(&model.MmAddress{}).Create(param)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return errors.New("新建地址失败")
	}
	return nil
}
