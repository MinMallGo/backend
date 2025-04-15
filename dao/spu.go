package dao

import (
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"strings"
	"time"
)

type SpuDao struct {
	db *gorm.DB
}

func NewSpuDao(db *gorm.DB) *SpuDao {
	return &SpuDao{
		db: db,
	}
}

func (u *SpuDao) Exists(id ...int) bool {
	if len(id) == 0 {
		return false
	}

	res := &[]model.MmSpu{}
	return u.db.Model(&model.MmSpu{}).Where("id in ?", id).Find(res).RowsAffected == int64(len(id))
}

func (u *SpuDao) Create(create *dto.SpuCreate) error {
	images := ""
	if len(create.Img) > 0 {
		images += strings.Join(create.Img, ",")
	}

	param := &model.MmSpu{
		Name:       create.Name,
		Desc:       create.Desc,
		Imgs:       images,
		Code:       util.SpuCode(),
		BrandID:    int32(create.BrandId), // TODO 需要验证和绑定品牌
		CategoryID: int32(create.CategoryId),
		MerchantID: int32(create.MerchantId),
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}
	return u.db.Model(&model.MmSpu{}).Create(&param).Error
}

func (u *SpuDao) Update(update *dto.SpuUpdate) error {
	images := ""
	if len(update.Img) > 0 {
		images += strings.Join(update.Img, ",")
	}

	param := &model.MmSpu{
		Name:       update.Name,
		Desc:       update.Desc,
		Imgs:       images,
		BrandID:    int32(update.BrandId),
		CategoryID: int32(update.CategoryId),
		MerchantID: int32(update.MerchantId),
		UpdateTime: time.Now(),
	}

	return u.db.Model(&model.MmSpu{}).Where("status = ?", constants.NormalStatus).Where("id = ?", update.Id).Updates(param).Error

}

func (u *SpuDao) Delete(delete *dto.SpuDelete) error {
	update := &model.MmSpu{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}

	return u.db.Model(&model.MmSpu{}).Select("status", "delete_time").Where("id = ?", delete.Id).Updates(update).Error
}

func (u *SpuDao) Search(search *dto.SpuSearch) (*dto.PaginateCount, error) {
	// id, name, page, limit
	result := &dto.PaginateCount{}
	category := &[]dto.SpuSearchResponse{}

	// 先写吧，优化等以后再来说

	whereStr := "status = ?"
	var param []interface{}
	param = append(param, constants.NormalStatus)

	if search.Id > 0 {
		whereStr += " AND id = ?"
		param = append(param, search.Id)
	}

	if len(search.Name) > 0 {
		whereStr += " AND name LIKE ?"
		param = append(param, "%"+search.Name+"%")
	}

	if search.MerchantId != 0 {
		whereStr += " AND merchant_id = ?"
		param = append(param, search.MerchantId)
	}

	if search.CategoryId != 0 {
		whereStr += " AND category_id = ?"
		param = append(param, search.CategoryId)
	}

	tx := u.db.Model(&model.MmSpu{}).Debug().
		Offset((search.Page-1)*search.Size).
		Limit(search.Size).
		Where(whereStr, param...).
		Preload("Category").
		Preload("Sku").
		Find(&category)
	if tx.Error != nil {
		return result, tx.Error
	}

	var count int64
	err := u.db.Model(&model.MmSpu{}).Debug().Where(whereStr, param...).Count(&count).Error
	if err != nil {
		return result, err
	}

	result = &dto.PaginateCount{
		Data:  category,
		Page:  search.Page,
		Size:  search.Size,
		Count: int(count),
	}

	return result, nil
}
