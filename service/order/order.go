package order

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service/pay"
	"mall_backend/util"
	"net/http"
	"strconv"
)

// TODO 这里创建订单和下单需要两个步骤，先创建订单，然后才是订单支付。先不要考虑那么精细化

func Create(c *gin.Context, create *dto.OrderCreate) {
	//获取用户id
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Failure(c, "下单失败：请登录后操作")
		return
	}

	// 1. 事务启动
	// 2. 检查商品，检查规格
	// 3. 地址检查
	// 4. 验证优惠券并计算价格
	// 5. 构造数据
	spuIds := make([]int, 0, len(create.Product)) // 这里要去重
	skuIds := make([]int, 0, len(create.Product))
	mSku := make(map[int]int, len(create.Product))
	seed := map[int]struct{}{}
	for _, v := range create.Product {
		if _, ok := seed[v.SpuID]; !ok {
			seed[v.SpuID] = struct{}{}
			spuIds = append(spuIds, v.SpuID)
		}
		skuIds = append(skuIds, v.SkuID)
		mSku[v.SkuID] = v.Num
	}

	//使用事务来创建
	if create.OrderType == dao.OrderTypeSecKill && len(create.Coupons) > 0 {
		response.Failure(c, "秒杀商品不能合并下单")
		return
	}

	if !dao.NewSpuDao(util.DBClient()).Exists(spuIds...) {
		response.Error(c, errors.New("创建订单失败：请选择正确的商品"))
		return
	}
	if !dao.NewSkuDao(util.DBClient()).Exists(skuIds...) {
		response.Error(c, errors.New("创建订单失败：请选择正确的商品规格"))
		return
	}

	if !dao.NewAddrDao(util.DBClient()).ExistsWithUser(int(user.ID), create.AddressID) {
		response.Error(c, errors.New("创建订单失败：请选择正确的地址"))
		return
	}

	if err = dao.NewCouponDao(util.DBClient()).ExistsWithUser(int(user.ID), create.Coupons...); err != nil {
		response.Error(c, errors.Join(errors.New("创建订单失败："), err))
		return
	}

	if err = dao.NewUserCouponDao(util.DBClient()).Exists(int(user.ID), create.Coupons...); err != nil {
		response.Error(c, errors.Join(errors.New("创建订单失败：请选择检查个人优惠券："), err))
		return
	}

	ctx := &Context{
		Param:     create,
		UniqueSku: mSku,
		UserID:    int(user.ID),
	}
	res, err := GetStrategy(create.OrderType).CreateOrder(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, res)
	return
}

func OrderPay(c *gin.Context, order *dto.PayOrder) {
	//// TODO 这里是支付订单，看有没有第三方的对接
	//// TODO 这里是消息通知  可能需要拆分。先弄个邮件发送得了
	//// TODO
	//// TODO 检查并支付  直接用事务吧
	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	var aliPayBody []byte
	if err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		orderRes, err := dao.NewOrderDao(tx).CheckBeforePay(order, int(user.ID))
		if err != nil {
			return err
		}

		more, err := dao.NewOrderSpuDao(tx).More(order.OrderCode)
		if err != nil {
			return err
		}
		var stock []dao.StockUpdate
		for _, spu := range *more {
			stock = append(stock, dao.StockUpdate{
				ID:  int(spu.SkuID),
				Num: int(spu.Count),
			})
		}

		err = dao.NewSkuDao(tx).StockDecrease(&stock)
		if err != nil {
			return err
		}

		// 检查库存，支付的时候锁库存

		// TODO 支付对接
		// TODO 消息通知

		// 这里只是发起支付，还没支付成功，所以不需要更改订单状态。可以在redis里面锁定这个订单的支付状态
		aliPayStr := ""
		aliPayStr, aliPayBody, err = HandleAliPcPay(order, orderRes.PayAmount, "Iphone6 16G")
		if err != nil {
			return err
		}
		if err = dao.NewOrderPayLogDao(tx).Create(more, aliPayStr); err != nil {
			return err
		}
		// 更新订单的状态
		return nil
	}); err != nil {
		response.Error(c, err)
		return
	}

	response.HtmlSuccess(c, string(aliPayBody))
	//response.Success(c, []string{})
	return
}

func CancelOrder(c *gin.Context, order *dto.CancelOrder) {
	//user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	//if err != nil {
	//	response.Error(c, err)
	//	return
	//}
	//// 取消订单的步骤：
	//// 1. 查看订单是否支付。 IsPaid
	//// 2. 归还库存
	//// 3. 归还优惠券
	//// 4. 归还用户优惠券
	//// 5.  TODO 如果支付了的话，要走退款流程
	//// 6. 更改订单状态
	//res := &model.MmOrder{}
	//err = util.DBClient().Transaction(func(tx *gorm.DB) error {
	//	// 查询订单是否支付
	//	res, err = dao.NewOrderDao(tx).IsPaid(order, int(user.ID))
	//	if err != nil {
	//		return err
	//	}
	//	// TODO 退款流程
	//	time.Sleep(time.Second * 5)
	//
	//	// 通过主订单查询子订单，并修改子订单的状态以及归还商品库存
	//	subOrders, err := dao.NewOrderSpuDao(tx).More(res.OrderCode)
	//	if err != nil {
	//		return err
	//	}
	//	// 归还商品库存
	//	if err = dao.NewSkuDao(tx).IncreaseStock(subOrders); err != nil {
	//		return err
	//	}
	//	// 修改子订单状态
	//	if err = dao.NewOrderSpuDao(tx).Cancel(res.OrderCode); err != nil {
	//		return err
	//	}
	//	// 修改住主订单状态
	//	if err = dao.NewOrderDao(tx).Cancel(order, int(user.ID)); err != nil {
	//		return err
	//	}
	//	// 1. 归还优惠券
	//	// 2. 归还用户优惠券
	//	if len(res.Coupons) > 0 {
	//		tmp := strings.Split(res.Coupons, ",")
	//		ids := make([]int, 0, len(tmp))
	//		for _, couponId := range tmp {
	//			id, err := strconv.Atoi(couponId)
	//			if err != nil {
	//				return errors.New("归还优惠券失败，请联系客服")
	//			}
	//			ids = append(ids, id)
	//		}
	//		if err = dao.NewCouponDao(tx).Cancel(ids...); err != nil {
	//			return err
	//		}
	//		if err = dao.NewUserCouponDao(tx).Cancel(int(user.ID), ids...); err != nil {
	//			return err
	//		}
	//	}
	//	// 记录退款日志
	//	if err = dao.NewOrderCancelDao(tx).Create(res, int(user.ID)); err != nil {
	//		return err
	//	}
	//	return nil
	//})
	//if err != nil {
	//	response.Error(c, err)
	//	return
	//}
	//response.Success(c, []string{})
	//return
}

func OrderExpire(c *gin.Context, order *dto.OrderExpire) {
	//user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	//if err != nil {
	//	response.Error(c, err)
	//	return
	//}
	//
	//// 1.归还库存，
	//// 2.这玩意儿
	//// 取消订单的步骤：
	//// 1. 查看订单是否支付。 IsPaid
	//// 2. 归还库存
	//// 3. 归还优惠券
	//// 4. 归还用户优惠券
	//// 5.  TODO 如果支付了的话，要走退款流程
	//// 6. 更改订单状态
	//err = util.DBClient().Transaction(func(tx *gorm.DB) error {
	//	// 通过主订单查询子订单，并修改子订单的状态以及归还商品库存
	//	subOrders, err := dao.NewOrderSpuDao(tx).More(order.OrderCode)
	//	if err != nil {
	//		return err
	//	}
	//	// 归还商品库存
	//	if err = dao.NewSkuDao(tx).IncreaseStock(subOrders); err != nil {
	//		return err
	//	}
	//	// 修改子订单状态
	//	if err = dao.NewOrderSpuDao(tx).Cancel(order.OrderCode); err != nil {
	//		return err
	//	}
	//	// 修改主订单状态
	//	if err = dao.NewOrderDao(tx).Expire(order, int(user.ID)); err != nil {
	//		return err
	//	}
	//	// 1. 归还优惠券
	//	// 2. 归还用户优惠券
	//	res, err := dao.NewOrderDao(tx).One(order.OrderCode, order.ID)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if len(res.Coupons) > 0 {
	//		tmp := strings.Split(res.Coupons, ",")
	//		ids := make([]int, 0, len(tmp))
	//		for _, couponId := range tmp {
	//			id, err := strconv.Atoi(couponId)
	//			if err != nil {
	//				return errors.New("归还优惠券失败，请联系客服")
	//			}
	//			ids = append(ids, id)
	//		}
	//		if err = dao.NewCouponDao(tx).Cancel(ids...); err != nil {
	//			return err
	//		}
	//		if err = dao.NewUserCouponDao(tx).Cancel(int(user.ID), ids...); err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//})
	//if err != nil {
	//	response.Error(c, err)
	//	return
	//}
	//response.Success(c, []string{})
	//return
}

// HandleAliPcPay 创建支付宝网页端支付的页面
func HandleAliPcPay(order *dto.PayOrder, amount int64, subject string) (reqStr string, bodyStr []byte, err error) {
	ali := pay.NewAlipay(true)
	m := map[string]string{
		"out_trade_no": order.OrderCode,
		"total_amount": strconv.FormatInt(amount/100, 10),
		"subject":      subject, //"Iphone6 16G",
		"product_code": "FAST_INSTANT_TRADE_PAY",
	}
	aliPayStr := ali.Generate("alipay.trade.page.pay", m)
	reqStr = aliPayStr

	res, err := http.Get(aliPayStr)
	if err != nil {
		return
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		err = errors.New(fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body))
		return
	}
	if err != nil {
		return
	}
	bodyStr = body
	return
}
