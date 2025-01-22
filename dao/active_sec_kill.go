package dao

import (
	"errors"
	"log"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"time"
)

// SecKillCreate 创建秒杀数据
func SecKillCreate(create *dto.SecKillCreate) error {
	// 数据验证，数据清洗，数据写入
	startTime, err := time.Parse("2006-01-02 15:04:05", create.ActiveStartTime)
	if err != nil {
		return err
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", create.ActiveEndTime)
	if err != nil {
		return err
	}

	// 批量写入了spec
	// 如果指定了spec_id 的话查询的时候就要只差这个
	res := &[]model.MmSpec{}

	whereStr := "status = ?"
	param := []interface{}{}
	param = append(param, constants.NormalStatus)

	if create.SpecId > 0 {
		whereStr += " AND id = ?"
		param = append(param, create.SpecId)
	}

	errx := util.DBClient().Model(&model.MmSpec{}).Where(whereStr, param...).
		Where("id = ?", create.SkuId).Find(res)
	if errx.Error != nil {
		return errx.Error
	}

	if len(*res) < 1 {
		return errors.New("请选择正确的商品规格")
	}
	insert := &[]model.MmActiveSeckill{}
	for _, spec := range *res {
		tmp := model.MmActiveSeckill{
			SpuID:      int32(create.SpuId),
			SkuID:      int32(create.SkuId),
			Stock:      int32(create.Stock),
			SpecID:     spec.ID,
			StartTime:  startTime,
			EndTime:    endTime,
			Name:       create.Name,
			Desc:       create.Desc,
			Price:      int32(create.Price),
			Status:     false,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		}
		*insert = append(*insert, tmp)
	}
	log.Println(insert)
	createRes := util.DBClient().Model(model.MmActiveSeckill{}).Create(insert)
	if createRes.Error != nil {
		return createRes.Error
	}

	return nil
}
