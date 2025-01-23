// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmSku = "mm_sku"

// MmSku  规格详细表- 保存规格的详细信息
type MmSku struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title      string    `gorm:"column:title;not null;comment:规格的组合名称" json:"title"` // 规格的组合名称
	SpuID      int32     `gorm:"column:spu_id;not null;comment:商品ID" json:"spu_id"`  // 商品ID
	Price      int32     `gorm:"column:price;comment:price * 100" json:"price"`      // price * 100
	Strock     int32     `gorm:"column:strock;comment:商品库存" json:"strock"`           // 商品库存
	Spces      string    `gorm:"column:spces;comment:存规格的json格式数据" json:"spces"`     // 存规格的json格式数据
	Status     bool      `gorm:"column:status;not null;default:1" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
}

// TableName MmSku's table name
func (*MmSku) TableName() string {
	return TableNameMmSku
}
