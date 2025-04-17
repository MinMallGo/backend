package active

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"mall_backend/dao"
	"mall_backend/dto"
)

type Context struct {
	ActiveID int
	Param    interface{}
}

type Strategy interface {
	Create(*gorm.DB, *Context) error
}

func GetStrategy(activeType int) Strategy {
	switch activeType {
	case dao.ActiveSecKill:
		return &SecKill{}
	case dao.ActiveGroupBuying:
		panic("need complete")
	default:
		panic("unknown activeType need to complete")
	}
}

type SecKill struct{}

func (s *SecKill) Create(tx *gorm.DB, ctx *Context) error {
	secKill, ok := ctx.Param.(*dto.ActiveCreate)
	if !ok {
		return errors.New("assert ctx with error: param not a type of dto.ActiveCreate")
	}
	log.Printf("secKill param:%#v,%#v", secKill, ctx.ActiveID)
	// 调用secKill的创建。其他活动自己实现
	seckillID, err := dao.NewSecKillDao(tx).Create(ctx.ActiveID, secKill)
	if err != nil {
		return err
	}

	err = dao.NewSecKillProductDao(tx).Create(seckillID, &secKill.SecKills)
	if err != nil {
		return err
	}
	return err
}
