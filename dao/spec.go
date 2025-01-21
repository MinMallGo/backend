package dao

import (
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/util"
)

func SpecGetOneByIds(ids []int) (*[]model.MmSpec, error) {
	res := &[]model.MmSpec{}
	err := util.DBClient().Model(&model.MmSpec{}).Where("status = ?", constants.NormalStatus).First(res, ids)
	if err != nil {
		return nil, err.Error
	}
	return res, nil
}
