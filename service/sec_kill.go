package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	constants "mall_backend/constant"
	"mall_backend/dao"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
)

func SecKillCreate(c *gin.Context, create *dto.SecKillCreate) {
	// 判断活动，判断秒杀
	// 判断 sec_kill sku,spu,price,stock
	SecKillInfo := &dto.SecKillInfo{
		SpuId:            create.SpuId,
		SkuId:            create.SkuId,
		Stock:            create.Stock,
		SecKillStartTime: create.SecKillStartTime,
		SecKillEndTime:   create.SecKillEndTime,
		Price:            create.Price,
	}
	if err := secKillChecker(SecKillInfo); err != nil {
		response.Failure(c, err.Error())
		return
	}

	activeCreate := &dto.ActiveCreate{
		Name:            create.Name,
		Type:            create.Type,
		Desc:            "",
		ActiveStartTime: "",
		ActiveEndTime:   "",
		Coupons:         nil,
	}

	//// 这里要用事务包裹
	// 创建活动
	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		_, err := dao.ActiveCreate(activeCreate)
		if err != nil {
			return err
		}
		// 创建秒杀
		_, err = dao.SecKillCreate(create)
		if err != nil {
			return err
		}
		return nil
	})
	// TODO
	if err != nil {
		response.Failure(c, err.Error())
		return
	}

	//// 这里要用事务包裹
	response.Success(c, []string{})
	return
}
func SecKillUpdate(c *gin.Context, update *dto.SecKillUpdate) {}

// SecKillDelete 删除秒杀内容
func SecKillDelete(activeId int) error {
	err := util.DBClient().Model(&model.MmActiveSeckill{}).Select("status").Where("id = ?", activeId).
		Updates(model.MmActiveSeckill{Status: constants.BanStatus})
	if err != nil {
		return err.Error
	}
	return nil
}

func secKillChecker(param *dto.SecKillInfo) error {
	var err error
	if !SpuExists(param.SpuId) {
		err = errors.Join(err, errors.New("请选择正确的商品"))
	}

	if !SkuExists(param.SkuId) {
		err = errors.Join(err, errors.New("请选择正确的规格"))
	}

	spec, errx := dao.SpecGetOneByIds([]int{param.SkuId})
	if errx != nil || len(*spec) < 1 {
		err = errors.Join(err, errx)
	}
	*spec = (*spec)[:1]

	if int32((*spec)[0].Price) < int32(param.Price) {
		err = errors.Join(err, errors.New("秒杀价格比原价高"))
	}

	if int32((*spec)[0].Strock) < int32(param.Stock) {
		err = errors.Join(err, errors.New("库存不够"))
	}

	return err
}
