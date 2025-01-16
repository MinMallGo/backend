package constants

const (
	UserLoginTTL = 60 * 60 * 24 // 用户登录状态保存时间
)

// MinDateTime 定义一些该死的常量 比如 datetime 没有默认值的时候时 0000-00-00 00:00:00 会报错
const (
	MinDateTime string = "1970-01-01"
)
