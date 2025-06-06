// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmAddress = "mm_address"

// MmAddress 用户地址表
type MmAddress struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID      int32     `gorm:"column:user_id;not null;comment:用户id" json:"user_id"`                              // 用户id
	Name        string    `gorm:"column:name;not null;comment:收货人姓名" json:"name"`                                   // 收货人姓名
	Phone       string    `gorm:"column:phone;not null;comment:收货人手机号" json:"phone"`                                // 收货人手机号
	Address     string    `gorm:"column:address;not null;comment:收货人地址 === 额这里以后以后再说怎么拆分地址" json:"address"`         // 收货人地址 === 额这里以后以后再说怎么拆分地址
	HouseNumber string    `gorm:"column:house_number;not null;comment:门牌号" json:"house_number"`                     // 门牌号
	IsDefault   bool      `gorm:"column:is_default;not null;default:1;comment:是否默认地址" json:"is_default"`            // 是否默认地址
	Tag         string    `gorm:"column:tag;comment:地址标签。以后再扩展" json:"tag"`                                         // 地址标签。以后再扩展
	Status      bool      `gorm:"column:status;not null;default:1;comment:状态： 0 删除 2 禁用 3 账号异常 1 正常" json:"status"` // 状态： 0 删除 2 禁用 3 账号异常 1 正常
	CreateTime  time.Time `gorm:"column:create_time;comment:创建时间" json:"create_time"`                               // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;comment:更新时间" json:"update_time"`                               // 更新时间
	DeleteTime  time.Time `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                               // 删除时间
}

// TableName MmAddress's table name
func (*MmAddress) TableName() string {
	return TableNameMmAddress
}
