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

// Refund 退款
func Refund(c *gin.Context) {
	param := &dto.RefundOrder{}
	if err := c.ShouldBind(param); err != nil {
		response.Error(c, err)
		return
	}
	order.Refund(c, param)
}
