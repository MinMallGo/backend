package dao

import (
	"gorm.io/gorm"
	"mall_backend/dao/model"
)

type MenuDao struct {
	db *gorm.DB
}
type RoleMenuDAO struct {
	db *gorm.DB
}

func NewMenuDao(db *gorm.DB) *MenuDao {
	return &MenuDao{
		db: db,
	}
}
func NewRoleDao(db *gorm.DB) *RoleMenuDAO {
	return &RoleMenuDAO{
		db: db,
	}
}

func (d *MenuDao) Search() ([]model.MmMenu, error) {
	// 关联查询spu，sku
	var users []model.MmMenu
	fd := d.db.Model(&model.MmMenu{}).
		Find(&users)
	if fd.Error != nil {
		return nil, fd.Error
	}
	return users, nil
}

func (r *RoleMenuDAO) SearchRole() ([]model.MmRoleMenu, error) {
	var roleMenus []model.MmRoleMenu
	if err := r.db.Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}
