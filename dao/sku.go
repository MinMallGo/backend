package dao

import (
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/util"
)

func SkuGetOneByIds(id int) (*[]model.MmSku, error) {
	res := &[]model.MmSku{}
	err := util.DBClient().Model(&model.MmSku{}).Where("status = ?", constants.NormalStatus).Where("id=?", id).First(res)
	if err != nil {
		return nil, err.Error
	}
	return res, nil
}
