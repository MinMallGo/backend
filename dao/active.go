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

func (d *ActiveDao) Update(update *dto.ActiveUpdate) error {
	tx := d.db.Model(&model.MmActive{}).Where("id = ?", update.ID).Updates(map[string]interface{}{
		"name":       update.Name,
		"type":       update.Type,
		"desc":       update.Desc,
		"start_time": update.StartTime,
		"end_time":   update.EndTime,
	})
	if tx.Error != nil || tx.RowsAffected == 0 {
		return tx.Error
	}

	return nil
}

func (d *ActiveDao) Exists(id int) error {
	var count int64
	tx := d.db.Model(&model.MmActive{}).Where("id = ?", id).Count(&count)
	if tx.Error != nil || count == 0 {
		return errors.New("更新活动失败：请选择正确的活动")
	}
	return nil
}

func (d *ActiveDao) BeforePurchase(purchase *dto.ActivePurchase) error {
	tx := d.db.Model(&model.MmActive{}).Where("id = ? AND start_time <= ? AND end_time >= ? AND status = ?",
		purchase.ActiveID, time.Now(), time.Now(), true,
	).First(&model.MmActive{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("购买失败：活动不存在或者不在活动时间内")
	}
	return nil
}
