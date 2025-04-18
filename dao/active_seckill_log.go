package dao

import (
	"errors"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"time"
)

// 为了验证用一个用户同一场秒杀，同一个商品只能下单一次

type SeckillLog struct {
	db *gorm.DB
}

func NewSeckillLog(db *gorm.DB) *SeckillLog {
	return &SeckillLog{db: db}
}

func (s *SeckillLog) Create(userId int, p *dto.ActivePurchase) error {
	param := &model.MmSeckillLog{
		ActiveID:   int32(p.ActiveID),
		ActiveType: int32(ActiveSecKill),
		SkuID:      int32(p.Product.SkuID),
		UserID:     int32(userId),
		CreateAt:   time.Now(),
	}
	tx := s.db.Model(&model.MmSeckillLog{}).Create(param)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *SeckillLog) Exists(userId int, purchase *dto.ActivePurchase) error {
	var count int64
	tx := s.db.Model(&model.MmSeckillLog{}).Where("user_id = ? AND active_id = ? AND sku_id = ? AND active_type = ?",
		userId, purchase.ActiveID, purchase.Product.SkuID, ActiveSecKill).Count(&count)
	if tx.Error != nil {
		return tx.Error
	}
	if count >= 1 {
		return errors.New("已经参与过秒杀活动")
	}
	return nil
}
