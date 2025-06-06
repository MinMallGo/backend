package main

import (
	"github.com/gin-gonic/gin"
	"mall_backend/middleware"
	routes "mall_backend/route"
	"mall_backend/service"
	"mall_backend/service/queue"
	"mall_backend/util"
)

func main() {
	// 开启控制台颜色
	gin.ForceConsoleColor()

	g := gin.Default()

	//logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//// 记录日志
	//g.Use(gin.Logger())
	// 从错误中恢复
	g.Use(gin.Recovery())
	// 跨域中间件
	g.Use(middleware.CORS())
	// IP限流中间件
	g.Use(middleware.IPLimiter())
	//// 注册日志记录中间件
	//g.Use(sloggin.New(logger))
	// 路由注册
	routes.Register(g)
	// 启动前初始化
	initNecessary()
	service.Menu()
	// 启动
	panic(g.Run())
}

func initNecessary() {
	// 初始化配置文件解析
	util.ConfigInstance("config.ini")
	// 初始化DB,Redis连接，以提前发发现错误
	if _, err := util.DBInstance(); err != nil {
		panic(err)
	}

	if _, err := util.RedisInstance(); err != nil {
		panic(err)
	}
	// 启动支付订单过期扫描订单
	go queue.OrderExpireQueue()
	// 启动同步查询订单是否支付
	go queue.OrderPayQueue()
	// 注册自定义验证规则
	util.ValidatorRegister()
}
