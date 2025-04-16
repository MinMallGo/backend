package order

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

func Create(c *gin.Context) {
	param := &dto.OrderCreate{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.OrderCreate(c, param)
}

func Pay(c *gin.Context) {
	param := &dto.PayOrder{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.OrderPay(c, param)
}

// TODO 这里要优化成消息队列，因为如果通过时间段去查询，会造成表锁或者范围锁，因为没有对时间加索引，

func Expire(c *gin.Context) {
	// 通过传入的订单id以及订单唯一编号来确定订单是否过期。额这玩意儿要么让队列来整吧，定时就行了
	param := &dto.OrderExpire{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.OrderExpire(c, param)
}

// CancelOrder 取消订单
func CancelOrder(c *gin.Context) {
	param := &dto.CancelOrder{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	service.CancelOrder(c, param)
}

// 退款

func Update(c *gin.Context) {}

func Delete(c *gin.Context) {}

func Search(c *gin.Context) {}
