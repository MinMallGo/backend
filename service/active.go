package service

import (
	"mall_backend/dao"
	"mall_backend/dto"
)

func ActiveCreate(create *dto.ActiveCreate) (id int, err error) {
	id, err = dao.ActiveCreate(create)
	return
}
