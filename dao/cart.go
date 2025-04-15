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

type CartDao struct {
	db *gorm.DB
}

func NewCartDao(db *gorm.DB) *CartDao {
	return &CartDao{
		db: util.DBClient(),
	}
}

func (d *CartDao) Create(create *dto.CartCreate, userId int) error {
	// 判断用户是否存在
	// TODO 这里就先不写了吧
	// 判断spu，sku是否存在
	if !NewSpuDao(d.db).Exists(create.SpuId) {
		return errors.New("添加到购物车失败：商品不存在")
	}

	if !NewSkuDao(d.db).Exists(create.SkuId) {
		return errors.New("添加到购物车失败：规格不存在")
	}
	// 检查当前spu,sku 是否存在，存在则修改数量而不是新增

	tx := d.db.Begin()
	ex := &model.MmCart{}
	if res := d.db.Model(&model.MmCart{}).Where("user_id = ? and spu_id = ? and sku_id = ? and status = ?", userId, create.SpuId, create.SkuId, constants.NormalStatus).Find(ex); res.Error == nil && res.RowsAffected == 1 {
		// 添加数量
		res := d.db.Model(&model.MmCart{}).Where("user_id = ? and spu_id = ? and sku_id = ? and status = ?", userId, create.SpuId, create.SkuId, constants.NormalStatus).Update("nums", int(ex.Nums)+create.Num)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		tx.Commit()
		return nil
	}

	param := &model.MmCart{
		UserID:     int32(userId),
		SpuID:      int32(create.SpuId),
		SkuID:      int32(create.SkuId),
		Nums:       int32(create.Num),
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	if err := d.db.Create(param).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (d *CartDao) Update(update *dto.CartUpdate, userId int) error {
	if !NewSpuDao(d.db).Exists(update.SpuId) {
		return errors.New("添加到购物车失败：商品不存在")
	}

	if !NewSkuDao(d.db).Exists(update.SkuId) {
		return errors.New("更新购物车失败：规格不存在")
	}

	if err := d.db.
		Model(&model.MmCart{}).
		Where("id = ? and user_id  = ? and spu_id = ?", update.Id, userId, update.SpuId).
		Updates(model.MmCart{
			SkuID: int32(update.SkuId),
			Nums:  int32(update.Num),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (d *CartDao) Delete(delete *dto.CartDelete, userId int) error {
	if err := d.db.Model(&model.MmCart{}).
		Where("id IN ? and user_id = ?", delete.Id, userId).
		Updates(map[string]interface{}{
			"status":      constants.BanStatus,
			"delete_time": time.Now(),
		}).Error; err != nil {
		return err
	}
	return nil
}

type searchMore struct {
	model.MmCart
	Spu model.MmSpu `json:"spu" gorm:"foreignKey:SpuID;references:ID"`
	Sku model.MmSku `json:"sku" gorm:"foreignKey:SkuID;references:ID"`
}

func (d *CartDao) Search(search *dto.CartSearch, userId int) (*dto.PaginateCount, error) {
	// 关联查询spu，sku
	result := &dto.PaginateCount{}
	res := &[]searchMore{}
	fd := d.db.Model(&model.MmCart{}).
		Preload("Sku").
		Preload("Spu").
		Where("user_id = ?", userId).
		Offset((search.Page - 1) * search.Size).
		Limit(search.Size).
		Find(&res)
	if fd.Error != nil {
		return &dto.PaginateCount{}, nil
	}

	var count int64
	d.db.Model(&model.MmCart{}).Where("user_id = ?", userId).Count(&count)

	result.Data = res
	result.Page = search.Page
	result.Size = search.Size
	result.Count = int(count)
	return result, nil
}
