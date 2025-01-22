package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	constants "mall_backend/constant"
	"mall_backend/dao"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
)

func SecKillCreate(c *gin.Context, create *dto.SecKillCreate) {
	// 构造秒杀插入的数据，然后进行检查
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

	// 构造活动创建的数据
	activeCreate := &dto.ActiveCreate{
		Name:            create.Name,
		Type:            create.Type,
		Desc:            create.Desc,
		ActiveStartTime: create.ActiveStartTime,
		ActiveEndTime:   create.ActiveEndTime,
		Coupons:         create.Coupons,
	}
	//// 这里要用事务包裹
	// 创建活动
	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		activeId, err := dao.ActiveCreate(activeCreate)
		if err != nil {
			return err
		}

		// 创建活动优惠券
		if len(create.Coupons) > 0 {
			err = dao.ActiveCouponCreate(activeId, &create.Coupons)
			log.Println(activeId, err)
			if err != nil {
				return err
			}
		}

		// 创建秒杀
		err = dao.SecKillCreate(create)
		if err != nil {
			return err
		}

		// 扣减库存
		return nil
	})

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

// 检查秒杀活动是否选择了正确是商品，规格，以及规格，规格价格，以及库存
func secKillChecker(param *dto.SecKillInfo) error {
	var err error
	if !SpuExists(param.SpuId) {
		err = errors.Join(err, errors.New("请选择正确的商品"))
	}

	if !SkuExists(param.SkuId) {
		err = errors.Join(err, errors.New("请选择正确的规格"))
	}

	spec, errx := dao.SpecGetOneByIds([]int{param.SkuId})
	if errx != nil || spec == nil || len(*spec) < 1 {
		err = errors.Join(err, errx)
	} else {
		*spec = (*spec)[:1]

		if int32((*spec)[0].Price) < int32(param.Price) {
			err = errors.Join(err, errors.New("秒杀价格比原价高"))
		}

		if int32((*spec)[0].Strock) < int32(param.Stock) {
			err = errors.Join(err, errors.New("库存不够"))
		}
	}

	return err
}
