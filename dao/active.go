package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

const (
	ActiveSecKill     int = 1 + iota // 秒杀
	ActiveGroupBuying                // 拼团
)

type ActiveDao struct {
	db *gorm.DB
}

func NewActiveDao(db *gorm.DB) *ActiveDao {
	return &ActiveDao{db: db}
}

func (d *ActiveDao) Create(create *dto.ActiveCreate) (int, error) {
	param := &model.MmActive{
		Name:       create.Name,
		Type:       int32(create.Type),
		Desc:       create.Desc,
		StartTime:  create.StartTime,
		EndTime:    create.EndTime,
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	err := d.db.Model(&model.MmActive{}).Create(param)
	if err.Error != nil || err.RowsAffected == 0 {
		return 0, errors.New("创建获取失败：创建主活动失败")
	}
	return int(param.ID), nil
}
