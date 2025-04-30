package order

import (
	"mall_backend/dao"
)

type Context struct {
	Param     interface{}
	UniqueSku map[int]int
	UserID    int
}

type Result struct {
	ID   int
	Code string
}

type Strategy interface {
	CreateOrder(ctx *Context) (Result, error)
	orderExpire() error
}

func GetStrategy(orderType int) Strategy {
	switch orderType {
	case dao.OrderTypeNormal:
		return &NormalStrategy{}
	case dao.OrderTypeSecKill:
		return &SeckillStrategy{}
	default:
		panic("not support order")
	}
}
