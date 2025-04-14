package dao

import (
	"fmt"
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
}

func NewSkuDao() *SkuDao {
	return &SkuDao{db: util.DBClient()}
}

func (d *SkuDao) Create(create *dto.SkuCreate) error {
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

func (d *SkuDao) Update(update *dto.SkuUpdate) error {
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

func (d *SkuDao) Exists(id ...int) bool {
	return d.db.
		Model(&model.MmSku{}).
		Select("id").
		Where("status = ?", constants.NormalStatus).
		Where("id in ?", id).
		First(&model.MmCategory{}).
		RowsAffected == int64(len(id))
}

func (d *SkuDao) Delete(id ...int) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		err := d.db.Model(&model.MmSku{}).Where("id in ?", id).Updates(map[string]interface{}{
			"status":      constants.BanStatus,
			"delete_time": time.Now().Format("2006-01-02 15:04:05"),
		}).Error
		if err != nil {
			return err
		}

		err = d.db.Model(&model.MmSpec{}).Debug().Where("sku_id in ?", id).Updates(map[string]interface{}{
			"status":      constants.BanStatus,
			"delete_time": time.Now().Format("2006-01-02 15:04:05"),
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *SkuDao) OneById(id int) (*model.MmSku, error) {
	res := &model.MmSku{}
	if err := d.db.Model(&model.MmSku{}).Where("status = ?", constants.NormalStatus).Where("spu_id = ?").First(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (d *SkuDao) More(search *dto.SkuSearch) (*dto.PaginateCount, error) {
	// TODO
	return nil, nil
}

func (d *SkuDao) MoreWithOrder(items *[]dto.ShoppingItem) (*[]model.MmSku, error) {
	res := &[]model.MmSku{}
	skuIds := make([]int, 0, len(*items))
	for _, item := range *items {
		skuIds = append(skuIds, item.SkuID)
	}

	tx := d.db.Model(&model.MmSku{}).Where("id in ?", skuIds).Where("status = ?", true).Find(&res)
	if tx.Error != nil {
		return res, tx.Error
	}

	return res, nil
}

// TODO 加库存减库存这里需要优化，不过以后再说
// TODO 改造成一个sql。用case when then

func (d *SkuDao) DecreaseStock(items *[]dto.ShoppingItem) bool {
	for _, item := range *items {
		tx := d.db.Model(&model.MmSku{}).
			Where("id = ?", item.SkuID).
			Where("status = ?", true).
			Update("stock", fmt.Sprintf("stock - %d", item.Num))
		if tx.Error != nil {
			return false
		}
	}
	return true
}

func (d *SkuDao) IncreaseStock(items *[]dto.ShoppingItem) bool {
	for _, item := range *items {
		tx := d.db.Model(&model.MmSku{}).
			Where("id = ?", item.SkuID).
			Where("status = ?", true).
			Update("stock", fmt.Sprintf("stock + %d", item.Num))
		if tx.Error != nil {
			return false
		}
	}
	return true
}
