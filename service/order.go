package service

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
)

// TODO 这里创建订单和下单需要两个步骤，先创建订单，然后才是订单支付。先不要考虑那么精细化

func OrderCreate(c *gin.Context, create *dto.OrderCreate) {
	//AddressID string         `json:"address_id" binding:"required"`
	//Coupons   []string       `json:"coupons" binding:"required"`
	//Product   []ShoppingItem `json:"product" binding:"required,dive"`
	//PayWay    string         `json:"pay_way" binding:"required,oneof=alipay wechatpay"`
	//Source    string         `json:"source" binding:"required,oneof=H5 min-program"`
	// 进行数据的校验
	// TODO 这里需要进行测试就是我在下单的时候，同时修改可用性，看能否使用成功
	// TODO 商品验证，商品规格验证
	// TODO 优惠券
	// TODO 满减
	// 这里只是创建订单
	// 获取用户id
	user, err := dao.NewUserDao().CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Failure(c, "下单失败：请登录后操作")
		return
	}

	err = dao.NewOrderDao().Create(create, int(user.ID))
}

func OrderPay(c *gin.Context) {
	// TODO 这里是支付订单，看有没有第三方的对接
	// TODO 这里是消息通知  可能需要拆分。先弄个邮件发送得了
}
