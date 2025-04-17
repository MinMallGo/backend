package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

type SecKillDao struct {
	db *gorm.DB
}

func NewSecKillDao(db *gorm.DB) *SecKillDao {
	return &SecKillDao{db: db}
}

func (s *SecKillDao) Create(activeID int, create *dto.ActiveCreate) (int, error) {
	param := &model.MmActiveSeckill{
		ActiveID:   int32(activeID),
		Name:       create.Name,
		StartTime:  create.StartTime,
		EndTime:    create.EndTime,
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}

	createRes := s.db.Model(&model.MmActiveSeckill{}).Create(param)
	if createRes.Error != nil {
		return 0, createRes.Error
	}

	return int(param.ID), nil
}
