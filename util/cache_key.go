package util

import (
	"fmt"
	"time"
)

// OrderPayWaitKey 返回正在支付订单的key
func OrderPayWaitKey() string {
	return "order_pay_wait"
}

// OrderExpireTime 获取订单过期时间。通过设置zSet的权重移除
func OrderExpireTime() int64 {
	return time.Now().Add(-OrderTTL()).Unix()
}

// OrderSetNXLockName 通过订单code统一管理锁，免得手写容易出问题
func OrderSetNXLockName(code string) string {
	return fmt.Sprintf("lock:order:%s", code)
}

func OrderTTL() time.Duration {
	return time.Minute * 10
}
