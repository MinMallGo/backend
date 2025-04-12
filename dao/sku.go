package dao

import (
	"gorm.io/gorm"
	"log"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

type SkuDao struct {
	db *gorm.DB
	m  model.MmSku
}

func NewSkuDao() *SkuDao {
	return &SkuDao{db: util.DBClient(), m: model.MmSku{}}
}

func (d SkuDao) Create(create *dto.SkuCreate) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		title, specs, err := NewSpecKeyDao().GenSkuData(&create.Spec)
		if err != nil {
			return err
		}

		// title = 规格值的名字
		// specs = 规格上下级的关系
		param := &model.MmSku{
			Title:      title,
			SpuID:      int32(create.SpuID),
			Price:      int32(create.Price),
			Stock:      int32(create.Stock),
			Spces:      string(specs),
			Status:     constants.NormalStatus,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		}
		err = d.db.Model(&model.MmSku{}).Create(&param).Error
		if err != nil {
			return err
		}

		err = SpecCreate(create.SpuID, int(param.ID), &create.Spec)
		if err != nil {
			return err
		}

		return nil
	})
}

func (d SkuDao) Update(update *dto.SkuUpdate) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		title, specs, err := NewSpecKeyDao().GenSkuData(&update.Spec)
		if err != nil {
			return err
		}

		// title = 规格值的名字
		// specs = 规格上下级的关系
		param := &model.MmSku{
			Title:      title,
			SpuID:      int32(update.SpuID),
			Price:      int32(update.Price),
			Stock:      int32(update.Stock),
			Spces:      string(specs),
			UpdateTime: time.Now(),
		}
		// 更新sku
		err = d.db.Model(&model.MmSku{}).Where("id = ?", update.Id).Updates(param).Error
		if err != nil {
			log.Println("update sku with error:", err)
			return err
		}
		// 更新spec
		err = SpecUpdate(update.SpuID, update.SkuID, &update.Spec)
		if err != nil {
			log.Println("updateSpec  with error:", err)
			return err
		}

		return nil

	})
}

func (d SkuDao) Exists(id ...int) bool {
	return d.db.
		Model(&model.MmSku{}).
		Select("id").
		Where("status = ?", constants.NormalStatus).
		Where("id in ?", id).
		First(&model.MmCategory{}).
		RowsAffected == int64(len(id))
}

func (d SkuDao) Delete(id ...int) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return d.db.Model(&model.MmSku{}).Where("id in ?", id).Update("status", false).Error
	})
}
