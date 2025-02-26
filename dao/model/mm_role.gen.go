// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmRole = "mm_role"

// MmRole 角色表，powers用 , 分割表示多个权限
type MmRole struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       string    `gorm:"column:name;comment:角色名" json:"name"`                     // 角色名
	Powers     string    `gorm:"column:powers;comment:包含的权限， 用 逗号 分割" json:"powers"`      // 包含的权限， 用 逗号 分割
	IsMenu     int32     `gorm:"column:is_menu;comment:是否菜单角色" json:"is_menu"`            // 是否菜单角色
	Status     bool      `gorm:"column:status;default:1;comment:1 正常 0 删除" json:"status"` // 1 正常 0 删除
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
}

// TableName MmRole's table name
func (*MmRole) TableName() string {
	return TableNameMmRole
}
