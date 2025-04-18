package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
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
			SeckillID:  int64(seckillID),
			SpuID:      int64(item.SpuID),
			SkuID:      int64(item.SkuID),
			Price:      int64(item.Price),
			Stock:      int32(item.Stock),
			Sold:       0,
			Version:    0,
			Status:     true,
			StartTime:  item.StartTime,
			EndTime:    item.EndTime,
			CreateTime: util.MinDateTime(),
			UpdateTime: util.MinDateTime(),
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

func (s *SecKillProductDao) Exist(p *dto.ActivePurchase) error {
	// 需要注意的是，这里仅仅是需要检查 - 1 而不是减很多库存
	tx := s.db.Model(&model.MmSeckillProduct{}).Where("seckill_id = ? AND sku_id = ? AND status = ? AND stock - 1 >= 0",
		p.SeckillID, p.Product.SkuID, true).Find(&model.MmSeckillProduct{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected < 1 {
		return errors.New("秒杀库存不足")
	}
	return nil
}

func (s *SecKillProductDao) Price(seckillID int, skuID int) (int, error) {
	res := &model.MmSeckillProduct{}
	tx := s.db.Model(&model.MmSeckillProduct{}).Select("price").Where("seckill_ID = ? AND sku_id = ? AND status = ?",
		seckillID, skuID, true).Find(res)
	if tx.Error != nil || tx.RowsAffected < 1 {
		return 0, tx.Error
	}

	return int(res.Price), nil
}

func (s *SecKillProductDao) DecrStock(seckillID int, skuID int) error {
	log.Println(seckillID, skuID)
	tx := s.db.Model(&model.MmSeckillProduct{}).
		Where("seckill_id = ?", seckillID).
		Where("sku_id = ?", skuID).
		Where("status = ?", true).
		Where("stock >= ?", 1).
		Updates(map[string]interface{}{
			"stock": clause.Expr{SQL: "`stock` - 1"},
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected < 1 {
		return errors.New("秒杀失败：扣减库存失败")
	}
	return nil
}
