package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
)

type SecKillProductDao struct {
	db *gorm.DB
}

func NewSecKillProductDao(db *gorm.DB) *SecKillProductDao {
	return &SecKillProductDao{db: db}
}

func (s *SecKillProductDao) Create(seckillID int, create *[]dto.SecKillCreate) error {
	param := make([]model.MmSeckillProduct, 0, len(*create))
	for _, item := range *create {
		param = append(param, model.MmSeckillProduct{
			SeckillID: int64(seckillID),
			SpuID:     int64(item.SpuID),
			SkuID:     int64(item.SkuID),
			Price:     int64(item.Price),
			Stock:     int32(item.Stock),
			Sold:      0,
			Version:   0,
			Status:    true,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			CreatedAt: util.MinDateTime(),
			UpdatedAt: util.MinDateTime(),
		})
	}

	createRes := s.db.Model(model.MmSeckillProduct{}).Create(param)
	if createRes.Error != nil || createRes.RowsAffected < int64(len(*create)) {
		return createRes.Error
	}

	return nil
}

func (s *SecKillProductDao) Update(secKillID int, update *[]dto.SecKillCreate) error {
	tx := s.db.Model(&model.MmSeckillProduct{}).Where("seckill_id = ?", secKillID).Update("status", false)
	if tx.Error != nil {
		return tx.Error
	}
	return s.Create(secKillID, update)
}
