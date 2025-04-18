package active

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
)

// Create 根据type 来选择策略
func Create(c *gin.Context, create *dto.ActiveCreate) {
	// 检查优惠券
	if err := dao.NewCouponDao(util.DBClient()).Exists(create.Coupons...); err != nil {
		response.Error(c, errors.New("创建活动：请选择正确的优惠券"))
		return
	}

	// 创建前检查sku以及库存
	enough := make([]dao.ExistsAndStock, 0, len(create.SecKills))
	spuIds := make([]int, 0, len(create.SecKills))
	distinctM := make(map[int]struct{}, len(create.SecKills))
	for _, kill := range create.SecKills {
		enough = append(enough, dao.ExistsAndStock{
			ID:    kill.SkuID,
			Stock: kill.Stock,
		})
		if _, ok := distinctM[kill.SpuID]; !ok {
			distinctM[kill.SpuID] = struct{}{}
			spuIds = append(spuIds, kill.SpuID)
		}
	}

	// 创建前检查spu是否存在
	if !dao.NewSpuDao(util.DBClient()).Exists(spuIds...) {
		response.Error(c, errors.New("创建活动：请选择正确的商品"))
		return
	}

	if len(enough) == 0 {
		response.Failure(c, "创建活动失败：请选择正确的商品")
		return
	}

	if err := dao.NewSkuDao(util.DBClient()).Enough(enough); err != nil {
		response.Error(c, err)
		return
	}

	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		// 1. 创建活动
		activeId, err := dao.NewActiveDao(tx).Create(create)
		if err != nil {
			return err
		}
		// 2. 创建活动优惠券。扣减优惠券使用数量
		if create.Type != dao.ActiveSecKill {
			if err = dao.NewActiveCouponDao(tx).Create(activeId, create.Coupons...); err != nil {
				return err
			}
			if err = dao.NewCouponDao(tx).CouponUse(create.Coupons...); err != nil {
				return err
			}
		}

		// 3. 根据active类型创建活动
		ctx := &Context{
			ActiveID: activeId,
			Param:    create,
		}
		if err = GetStrategy(create.Type).Create(tx, ctx); err != nil {
			return errors.Join(errors.New("创建活动失败：创建秒杀以及秒杀商品失败"), err)
		}

		// 5. 扣减商品的库存
		if err = dao.NewSkuDao(tx).SecKillCreate(&create.SecKills); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		response.Error(c, errors.New("创建活动失败: "+err.Error()))
		return
	}

	response.Success(c, []string{})
	return
}

func Update(c *gin.Context, update *dto.ActiveUpdate) {
	if err := dao.NewActiveDao(util.DBClient()).Exists(update.ID); err != nil {
		response.Error(c, err)
		return
	}
	
	// 检查优惠券
	if err := dao.NewCouponDao(util.DBClient()).Exists(update.Coupons...); err != nil {
		response.Error(c, errors.New("更新活动失败：请选择正确的优惠券"))
		return
	}

	// 创建前检查sku以及库存
	enough := make([]dao.ExistsAndStock, 0, len(update.SecKills))
	spuIds := make([]int, 0, len(update.SecKills))
	distinctM := make(map[int]struct{}, len(update.SecKills))
	for _, kill := range update.SecKills {
		enough = append(enough, dao.ExistsAndStock{
			ID:    kill.SkuID,
			Stock: kill.Stock,
		})
		if _, ok := distinctM[kill.SpuID]; !ok {
			distinctM[kill.SpuID] = struct{}{}
			spuIds = append(spuIds, kill.SpuID)
		}
	}

	// 创建前检查spu是否存在
	if !dao.NewSpuDao(util.DBClient()).Exists(spuIds...) {
		response.Error(c, errors.New("更新活动失败：请选择正确的商品"))
		return
	}

	if len(enough) == 0 {
		response.Failure(c, "更新活动失败：请选择正确的商品")
		return
	}

	if err := dao.NewSkuDao(util.DBClient()).Enough(enough); err != nil {
		response.Error(c, err)
		return
	}

	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		// 1. 创建活动
		err := dao.NewActiveDao(tx).Update(update)
		if err != nil {
			return err
		}
		// 2. 创建活动优惠券。扣减优惠券使用数量
		if update.Type != dao.ActiveSecKill {
			if err = dao.NewActiveCouponDao(tx).Update(update); err != nil {
				return err
			}
			if err = dao.NewCouponDao(tx).Cancel(update.Coupons...); err != nil {
				return err
			}
		}

		// 3. 根据active类型创建活动
		ctx := &Context{
			ActiveID: update.ID,
			Param:    update,
		}
		if err = GetStrategy(update.Type).Update(tx, ctx); err != nil {
			return errors.Join(errors.New("更新活动失败：创建秒杀以及秒杀商品失败"), err)
		}

		var seckillID int
		if seckillID, err = dao.NewSecKillDao(tx).GetID(update.ID); err != nil {
			return errors.New("更新活动失败")
		}
		// 5. 扣减商品的库存
		if err = dao.NewSkuDao(tx).SecKillUpdate(seckillID, &update.SecKills); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		response.Error(c, errors.New("更新活动失败: "+err.Error()))
		return
	}

	response.Success(c, []string{})
	return
}
