package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"io"
	"log"
	"mall_backend/dao"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service/pay"
	"mall_backend/util"
	"net/http"
	"strconv"
	"time"
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
		response.Error(c, errors.New(fmt.Sprintf("创建订单失败：请选择正确的商品:%v", spuIds)))
		return
	}
	if !dao.NewSkuDao(util.DBClient()).Exists(skuIds...) {
		response.Error(c, errors.New(fmt.Sprintf("创建订单失败：请选择正确的商品规格:%v", skuIds)))
		return
	}

	if !dao.NewAddrDao(util.DBClient()).ExistsWithUser(int(user.ID), create.AddressID) {
		response.Error(c, errors.New("创建订单失败：请选择正确的地址"))
		return
	}

	// 使用了优惠券再查询
	if len(create.Coupons) > 0 {
		if err = dao.NewCouponDao(util.DBClient()).ExistsWithUser(int(user.ID), create.Coupons...); err != nil {
			response.Error(c, errors.Join(errors.New("创建订单失败："), err))
			return
		}

		if err = dao.NewUserCouponDao(util.DBClient()).Exists(int(user.ID), create.Coupons...); err != nil {
			response.Error(c, errors.Join(errors.New("创建订单失败：请选择检查个人优惠券："), err))
			return
		}
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

// OrderSetNx 对某个订单进行加锁，以免订单被重复/多次支付
func OrderSetNx(lockName string, acquireTimeout, expireTime time.Duration) error {
	endTime := time.Now().Add(acquireTimeout)
	for time.Now().Before(endTime) {
		result, err := util.CacheClient().SetNX(context.Background(), lockName, "lock", expireTime).Result()
		if err != nil {
			return errors.New("订单支付：无法获取订单的锁，稍后重试。")
		}
		if result {
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}
	return errors.New("订单支付：无法获取订单的锁，稍后重试。")
}

// OrderReleaseLock 释放订单锁
func OrderReleaseLock(lockName string) error {
	return util.CacheClient().Del(context.Background(), lockName).Err()
}

// Pay 需要考虑支付成功以及支付失败，支付超时这三种可能行
//
//	使用zSet保存支付的order_code，避免多次重复发起支付请求。
//
// 1. 为什么要用zSet来保存？1. 他是个他是个有序的集合，然后里面的成员不能重复。然后我我这里根据过期时间设置他的权重。然后通过定时任务去让他过期
// 2. 需要用定时任务是扫，然后去掉过期的。避免阻塞正常下单业务流程
// TODO 这里好像会有问题，因为这里是返回一个html页面，如果这个页面被意外关闭，这里应该怎么来返回一个支付的页面呢
// TODO 真需要就加一个justPay的页面，先不管
func Pay(c *gin.Context, order *dto.PayOrder) {
	var err error

	// 1. 查询订单是否过期
	err = dao.NewOrderDao(util.DBClient()).IsExpire(order.OrderCode)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 使用分布式锁进行锁订单，获取到锁才行执行下面的步骤
	err = OrderSetNx(util.OrderSetNXLockName(order.OrderCode), util.OrderTTL(), util.OrderTTL())
	if err != nil {
		response.Error(c, err)
		return
	}

	defer func() {
		err := OrderReleaseLock(util.OrderSetNXLockName(order.OrderCode))
		if err != nil {
			// 记录日志，处理释放锁失败的情况
			log.Println("订单支付释放锁失败")
		}
	}()

	// 通过查询zSet中的订单号是否存在来判断是否能否继续走面的支付流程
	// 1. 错误则返回错误，error.Is(err,redis.Nil) 这里需要忽略，因为是没找到
	// 2. 如果权重（score）大于0说明正在支付中
	score, err := util.CacheClient().ZScore(context.TODO(), util.OrderPayWaitKey(), order.OrderCode).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		response.Error(c, err)
		return
	}

	if int64(score) > time.Now().Add(-util.OrderTTL()).Unix() {
		// 考虑一下，如果这里返回的地址用户点击后没有支付然乎关闭了。再点进来。订单还没有支付，应该继续返回支付页面才对撒
		// 这时候肯定再order_pay_log里面存在
		one, err := dao.NewOrderPayLogDao(util.DBClient()).One(order.OrderCode)
		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, one.PayQueryData)
		return
	}

	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	aliPayStr := ""
	if err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		orderRes, err := dao.NewOrderDao(tx).CheckBeforePay(order, int(user.ID))
		if err != nil {
			return err
		}

		more, err := dao.NewOrderSpuDao(tx).More(order.OrderCode)
		if err != nil {
			return err
		}
		// 请求支付宝支付地址
		aliPayStr, _, err = HandleAliPcPay(order, orderRes.PayAmount, fmt.Sprintf("min_mall:%s", orderRes.SubjectTitle))
		if err != nil && !errors.Is(err, redis.Nil) {
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
	response.Success(c, map[string]string{"url": aliPayStr})
	return
}

// Expire 这里来的都是过期的订单并且创建了支付请求的订单，如果过期了说明支付也过期了。需要归还库存
func Expire(orderCode string) error {
	//  订单过期，主要是看有没有支付，如果支付了就直接返回
	//  如果发起了支付但是没有支付就要归还库存
	//  怎么查看有没有发起支付？在zSet里面进行查询，如果支付了就存在，如果不存在说明已经支付成功，要么就是没支付但是需要过期了。因为当你下单的时候会同时调用pay接口，所以就要么还没支付
	order, err := dao.NewOrderDao(util.DBClient()).One(orderCode)
	if err != nil || order.ID == 0 {
		return err
	}

	// 如果支付状态大于当前状态为已支付就不管他
	if order.PaymentStatus >= dao.OrderStatusPaid {
		return nil
	}
	err = GetStrategy(int(order.OrderType)).orderExpire(orderCode)
	if err != nil {
		return err
	}
	return nil
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

// PaySuccess 订单支付完成之后，对主订单，子订单，支付记录进行成功操作
func PaySuccess(orderCode string) error {
	// 订单支付成功之后的逻辑处理
	// 1. 先检查订单是否存在
	// 2. 查看订单是否已经完成了修改，这里看需不需要给支付宝给个回调表示已经收到了
	// 3. 写入order_pay_log里面
	// 4. 修改order的支付状态
	// 5. 修改order_spu的状态

	err := util.DBClient().Transaction(func(tx *gorm.DB) error {
		isPaid, err := dao.NewOrderDao(tx).IsPaid(orderCode)
		if err != nil {
			return err
		}

		if !isPaid {
			return errors.New("order not found")
		}
		// 修改主订单状态
		err = dao.NewOrderDao(tx).PaySuccess(orderCode)
		if err != nil {
			return err
		}
		// 修改子订单状态
		err = dao.NewOrderSpuDao(tx).PaySuccess(orderCode)
		if err != nil {
			return err
		}
		// 修改订单支付记录的状态
		err = dao.NewOrderPayLogDao(tx).PaySuccess(orderCode)
		if err != nil {
			return err
		}

		_, err = util.CacheClient().ZRem(context.TODO(), util.OrderPayWaitKey(), orderCode).Result()
		if err != nil {
			return errors.New(fmt.Sprintf("移除支付成功订单号失败：%s\n", orderCode))
		}

		// TODO 通知

		// TODO 发货

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// 1. 通过订单号以及sku来进行验证和确认退款的金额
// 2. 进来操作的时候要加锁。避免
// 3. 在事务里面进行
// 2. 先检查订单，订单的sku是否都存在  返回 *[]model.MmOrderSpu{},error
// 3. 然后检查金额
// 4. 调用ali的退款接口
// 5. todo。这里还需要一个退款查询的队列

func Refund(c *gin.Context, refund *dto.RefundOrder) {
	// 1. 分布式锁整一个
	// 使用分布式锁进行锁订单，获取到锁才行执行下面的步骤

	user, err := dao.NewUserDao(util.DBClient()).CurrentUser(c.GetHeader("token"))
	if err != nil {
		response.Error(c, err)
		return
	}

	err = OrderSetNx(util.OrderSetNXLockName(refund.OrderCode), util.OrderTTL(), util.OrderTTL())
	if err != nil {
		response.Error(c, err)
		return
	}

	defer func() {
		err := OrderReleaseLock(util.OrderSetNXLockName(refund.OrderCode))
		if err != nil {
			// 记录日志，处理释放锁失败的情况
			log.Println("订单退款释放锁失败")
		}
	}()

	// 等哈儿退款成功了要进行，order，orderSpu支付状态的修改。额其实是sku
	// 计算当前这些订单的价格，并获取
	err = util.DBClient().Transaction(func(tx *gorm.DB) error {
		one, err := dao.NewOrderDao(tx).One(refund.OrderCode)
		if err != nil {
			return err
		}

		if one.UserID != int64(user.ID) {
			return errors.New("退款失败：请选择正确的订单")
		}

		if (*one).PaymentStatus == dao.OrderStatusNeedPay {
			return errors.New("退款失败：请选择正确的订单")
		}

		if (*one).PaymentStatus == dao.OrderStatusAllCancel {
			return errors.New("退款失败：当前订单已经全部退款")
		}

		refundOrder, err := dao.NewOrderSpuDao(tx).MoreRefundOrder(refund.OrderCode, refund.Skus...)
		if err != nil {
			return err
		}
		var amount int64
		stocks := make([]dao.StockUpdate, 0, len(*refundOrder))
		for _, v := range *refundOrder {
			amount += v.PayAmount
			stocks = append(stocks, dao.StockUpdate{
				ID:  int(v.SkuID),
				Num: int(v.Count),
			})
		}

		err = pay.NewAlipay(true).Cancel(map[string]string{
			"refund_amount": fmt.Sprintf("%.2f", float64(amount)/100),
			"out_trade_no":  refund.OrderCode,
		})
		if err != nil {
			return err
		}

		err = dao.NewOrderDao(tx).Refund(refund.OrderCode, amount)
		if err != nil {
			return err
		}

		err = dao.NewOrderSpuDao(tx).Refund(refundOrder)
		if err != nil {
			return err
		}

		// 创建退款记录
		// 这里应该修改退款的金额为实际退款金额，不一定就是订单的总金额
		one.PayAmount = amount
		err = dao.NewOrderRefundDao(tx).Create(one)
		if err != nil {
			return err
		}

		// 归还子订单的库存
		err = dao.NewSkuDao(tx).StockIncrease(&stocks)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "退单成功")
	return
}
