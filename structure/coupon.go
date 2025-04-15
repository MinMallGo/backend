package structure

type Context struct {
	Price    int  // 原价
	Category int  // 优惠券类型
	Rules    Rule //  规则类型
}

type Rule struct {
	Threshold     int
	DisCountPrice int
	DiscountRate  int
}
