package dao

import (
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

// SecKillCreate 创建秒杀数据
func SecKillCreate(create *dto.SecKillCreate) (int, error) {
	// 数据验证，数据清洗，数据写入
	param := &model.MmActiveSeckill{
		SpuID:      int32(create.SpuId),
		SkuID:      int32(create.SkuId),
		Stock:      int32(create.Stock),
		StartTime:  time.Time{},
		EndTime:    time.Time{},
		Name:       create.Name,
		Desc:       create.Desc,
		Price:      0,
		Status:     false,
		CreateTime: time.Now(),
		UpdateTime: time.Time{},
		DeleteTime: time.Time{},
	}
	res := util.DBClient().Model(model.MmActiveSeckill{}).Create(param)
	if res.Error != nil {
		return 0, res.Error
	}

	return int(param.ID), nil
}
