package queue

import (
	"context"
	"fmt"
	"mall_backend/service/order"
	"mall_backend/util"
	"time"
)

func OrderExpireQueue() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			// 扫描zSet里面过期的键值，并进行过期处理
			// 1. 因为zSet里面我是通过权重来设置成员，然后我现在这里可以通过当前时间来查询订单号。
			// 2. 然后我这里进行一个数据库的查询，筛选出支付完成的订单
			// 3. 这里只负责扫描过期的订单然后给他归还库存
			waitOrders, err := util.CacheClient().ZRange(context.TODO(), util.OrderPayWaitKey(), 0, time.Now().Add(-time.Minute).Unix()).Result()
			if err != nil {
				continue
			}
			for _, orderCode := range waitOrders {
				fmt.Printf("订单号：%s\n", orderCode)
				go func(orderCode string) {
					err := order.Expire(orderCode)
					if err != nil {
						return
					}
					// 这里就需要移除支付过期的内容
					_, err = util.CacheClient().ZRem(context.TODO(), util.OrderPayWaitKey(), orderCode).Result()
					if err != nil {
						fmt.Printf("移除支付过期key失败：%s\n", orderCode)
						return
					}
					fmt.Printf("移除支付过期key成功：%s\n", orderCode)
				}(orderCode)
				fmt.Println(time.Now().Add(time.Minute).Unix())
				// 传入到ali支付检查这里去看，支付有没有成功
				// 如果成功了，移除，然后修改订单状态。不管有没有成功，这些都该移除了
			}
		}
	}
}
