package middleware

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/response"
	"mall_backend/util"
	"strconv"
)

func AuthVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		login, err := util.CacheClient().Get(context.Background(), token).Result()
		if len(login) < 2 || err != nil {
			// 返回请登录的提示
			response.NeedLogin(c, "请登录后操作")
			c.Abort() // 终止请求
			return
		} else {
			// 调用该请求的剩余处理程序
			c.Next()
		}
	}
}
func MenuVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		login, err := util.CacheClient().Get(context.Background(), token).Result()
		if len(login) < 2 || err != nil {
			// 返回请登录的提示
			response.NeedLogin(c, "请登录后操作")
			c.Abort() // 终止请求
			return
		} else {
			request := c.Request.URL.Path
			var userData model.MmUser
			err = json.Unmarshal([]byte(login), &userData)
			if err != nil {
				response.NeedLogin(c, "无权限访问")
				c.Abort() // 终止请求
				return
			}
			menu, e := util.CacheClient().Get(context.Background(), strconv.Itoa(int(userData.Type))).Result()
			if e != nil {
				response.NeedLogin(c, "请登录后操作")
				c.Abort() // 终止请求
				return
			}
			var menuData []model.MmMenu
			err = json.Unmarshal([]byte(menu), &menuData)
			if err != nil {
				response.NeedLogin(c, "无权限访问")
				c.Abort() // 终止请求
				return
			}
			for _, menu := range menuData {
				if "/v1"+menu.Component == request {
					c.Next()
				}
			}
			response.NeedLogin(c, "无权限访问")
			c.Abort() // 终止请求
			return
			// 调用该请求的剩余处理程序
		}
	}
}

func AdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 根据路由来判断是否是商户/管理员
		// 往header里面加入用户
		c.Set("user", constants.AdminUser)
	}
}
