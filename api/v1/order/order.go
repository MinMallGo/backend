package order

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service/order"
)

func Create(c *gin.Context) {
	param := &dto.OrderCreate{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	order.Create(c, param)
}

func Pay(c *gin.Context) {
	param := &dto.PayOrder{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	order.Pay(c, param)
}

func Pay2(c *gin.Context) {
	param := &dto.PayOrder{}
	if err := c.ShouldBindQuery(&param); err != nil {
		response.Error(c, err)
		return
	}
	order.Pay(c, param)
}

// TODO expire 这里要优化成消息队列，因为如果通过时间段去查询，会造成表锁或者范围锁，因为没有对时间加索引，

// CancelOrder 取消订单
func CancelOrder(c *gin.Context) {
	param := &dto.CancelOrder{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	order.CancelOrder(c, param)
}

// 退款

func Update(c *gin.Context) {}

func Delete(c *gin.Context) {}

func Search(c *gin.Context) {}
