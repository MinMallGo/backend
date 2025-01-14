// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMmUser = "mm_user"

// MmUser 用户表
type MmUser struct {
	ID            int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Account       string    `gorm:"column:account;not null;comment:账号" json:"account"`                                // 账号
	Name          string    `gorm:"column:name;not null;comment:名字" json:"name"`                                      // 名字
	Password      string    `gorm:"column:password;not null;comment:密码" json:"password"`                              // 密码
	Salt          string    `gorm:"column:salt;not null;comment:密码盐" json:"salt"`                                     // 密码盐
	Type          int32     `gorm:"column:type;not null;comment:1 普通用户 2 商户 999 表示管理员" json:"type"`                   // 1 普通用户 2 商户 999 表示管理员
	Role          string    `gorm:"column:role;not null;comment:用户角色，一个用户可以对应多个角色" json:"role"`                       // 用户角色，一个用户可以对应多个角色
	Phone         string    `gorm:"column:phone;comment:手机号" json:"phone"`                                            // 手机号
	ThirdParty    string    `gorm:"column:third_party;comment:保存第三方登录方式的信息，比如什么微信登录，邮箱登录之类的" json:"third_party"`      // 保存第三方登录方式的信息，比如什么微信登录，邮箱登录之类的
	LastLoginIP   string    `gorm:"column:last_login_ip;comment:最新登录IP" json:"last_login_ip"`                         // 最新登录IP
	LastLoginTime time.Time `gorm:"column:last_login_time;comment:最新登录时间" json:"last_login_time"`                     // 最新登录时间
	Status        bool      `gorm:"column:status;not null;default:1;comment:状态： 0 删除 2 禁用 3 账号异常 1 正常" json:"status"` // 状态： 0 删除 2 禁用 3 账号异常 1 正常
	CreateTime    time.Time `gorm:"column:create_time;comment:创建时间" json:"create_time"`                               // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time;comment:更新时间" json:"update_time"`                               // 更新时间
	DeleteTime    time.Time `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                               // 删除时间
}

// TableName MmUser's table name
func (*MmUser) TableName() string {
	return TableNameMmUser
}
