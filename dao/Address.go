package dao

import (
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/util"
)

type AddrDao struct {
	db *gorm.DB
}

func NewAddrDao() *AddrDao {
	return &AddrDao{
		db: util.DBClient(),
	}
}

func (a *AddrDao) ExistsWithUser(userId int, id ...int) bool {
	tmp := &[]model.MmAddress{}
	return a.db.Model(&model.MmAddress{}).
		Where("id in ?", id).
		Where("user_id = ?", userId).
		Where("status = ?", constants.NormalStatus).
		Find(tmp).RowsAffected == int64(len(id))
}
