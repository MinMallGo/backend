package coupon

type Context struct {
	OriginalPrice float64                // 原价
	Rules         map[string]interface{} // 打折规则
	Category      string                 // 优惠券类型
}

type Strategy interface {
	ApplyStrategy(ctx *Context)
}
