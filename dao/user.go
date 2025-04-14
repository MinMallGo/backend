package dao

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao() *UserDao {
	return &UserDao{
		db: util.DBClient(),
	}
}

func (u *UserDao) CurrentUser(token string) (*model.MmUser, error) {
	jsonStr, err := util.CacheClient().Get(context.TODO(), token).Result()
	if err != nil {
		return nil, err
	}
	res := &model.MmUser{}
	err = json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *UserDao) Exists(id ...int) bool {
	res := &[]model.MmUser{}
	return u.db.Model(&model.MmUser{}).Where("id in ?", id).Find(res).RowsAffected == int64(len(id))
}

func (u *UserDao) Create(param *dto.CartCreate) {}

func (u *UserDao) Update(param *dto.CartCreate) {}

func (u *UserDao) Delete(param *dto.CartCreate) {}

func (u *UserDao) Search(param *dto.CartCreate) {}
