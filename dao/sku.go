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

// SkuCreateTransaction 创建商品规格。使用事务保证一致性
func SkuCreateTransaction(create *dto.SkuCreate) error {
	return util.DBClient().Transaction(func(tx *gorm.DB) error {
		title, specs, err := GenSkuData(&create.Spec)
		if err != nil {
			return err
		}

		// title = 规格值的名字
		// specs = 规格上下级的关系
		param := &model.MmSku{
			Title:      title,
			SpuID:      int32(create.SpuID),
			Price:      int32(create.Price),
			Strock:     int32(create.Stock),
			Spces:      string(specs),
			Status:     constants.NormalStatus,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		}
		err = util.DBClient().Model(&model.MmSku{}).Create(&param).Error
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

func SkuUpdateTransaction(update *dto.SkuUpdate) error {
	return util.DBClient().Transaction(func(tx *gorm.DB) error {
		title, specs, err := GenSkuData(&update.Spec)
		if err != nil {
			return err
		}

		// title = 规格值的名字
		// specs = 规格上下级的关系
		param := &model.MmSku{
			Title:      title,
			SpuID:      int32(update.SpuID),
			Price:      int32(update.Price),
			Strock:     int32(update.Stock),
			Spces:      string(specs),
			UpdateTime: time.Now(),
		}
		// 更新sku
		err = util.DBClient().Model(&model.MmSku{}).Where("id = ?", update.Id).Updates(param).Error
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

// SkuExists 通过ID查询该商品是否可用
func SkuExists(spuID int) bool {
	return util.DBClient().Model(&model.MmSku{}).Select("id").Where("status = ?", constants.NormalStatus).
		Where("id = ?", spuID).First(&model.MmCategory{}).RowsAffected != 0
}
