package queue

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"mall_backend/service/order"
	"mall_backend/service/pay"
	"mall_backend/util"
	"strconv"
	"time"
)

func OrderExpireQueue() {
	fmt.Println("订单过期结果查询队列启动成功")
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			// 扫描zSet里面过期的键值，并进行过期处理
			// 1. 因为zSet里面我是通过权重来设置成员，然后我现在这里可以通过当前时间来查询订单号。
			// 2. 然后我这里进行一个数据库的查询，筛选出支付完成的订单
			// 3. 这里只负责扫描过期的订单然后给他归还库存

			cmd := util.CacheClient().ZRangeByScore(context.TODO(), util.OrderPayWaitKey(), &redis.ZRangeBy{
				Min: strconv.Itoa(int(util.MinDateTime().Unix())),
				Max: strconv.Itoa(int(util.OrderExpireTime())),
			})

			waitOrders, err := cmd.Result()

			if err != nil {
				continue
			}
			for _, orderCode := range waitOrders {
				fmt.Printf("检查订单是否过期：订单号：%s\n", orderCode)
				go func(code string) {
					err := order.Expire(code)
					log.Printf("传入的订单号%s\n", code)
					if err != nil {
						return
					}
					// 这里就需要移除支付过期的内容
					_, err = util.CacheClient().ZRem(context.TODO(), util.OrderPayWaitKey(), code).Result()
					if err != nil {
						fmt.Printf("移除过期订单号失败：%s\n", code)
						return
					}
					fmt.Printf("移除过期订单号成功：%s\n", code)
				}(orderCode)

				// 传入到ali支付检查这里去看，支付有没有成功
				// 如果成功了，移除，然后修改订单状态。不管有没有成功，这些都该移除了
			}
		}
	}
}

// TODO 有个bug。就是已经支付的订单如果没有被及时处理，这里会被移除，所以你需要检查订单是否支付，如果没有支付才能移除
// TODO bug2 这里归还用户优惠券似乎出现了问题，需要你去看一下，你能做到吗
// OrderPayQueue 以同步的方式检查订单是否支付
func OrderPayQueue() {
	fmt.Println("同步支付结果查询队列启动成功")
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			cmd := util.CacheClient().ZRangeByScore(context.TODO(), util.OrderPayWaitKey(), &redis.ZRangeBy{
				Min: strconv.Itoa(int(time.Now().Add(time.Second).Unix())),
				Max: strconv.Itoa(int(time.Now().Add(util.OrderTTL()).Unix())),
			})
			waitOrders, err := cmd.Result()

			if err != nil {
				continue
			}
			for _, orderCode := range waitOrders {
				go func(orderCode string) {
					log.Printf("<UNK>%s\n", err)
					// 查询是否支付
					body, err := pay.NewAlipay(true).PayStatus("alipay.trade.query", map[string]string{
						"out_trade_no": orderCode,
					})

					if err != nil {
						log.Printf("pay status check is failed:%s\n", err)
						return
					}

					fmt.Printf("ali_pay response : %s\n", body)

					// 这里需要加上订单完成之后的逻辑操作：
					// 1.
					err = order.PaySuccess(orderCode)
					if err != nil {
						return
					}
				}(orderCode)
			}

		}
	}
}
