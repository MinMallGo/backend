package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/enumerate"
	"mall_backend/response"
	"mall_backend/util"
	"time"
)

// Login 用户登录只要账号密码
func Login(c *gin.Context, login *dto.UserLogin) {
	userRes := &model.MmUser{}
	util.DBClient().First(&model.MmUser{}).Where("account = ?", login.Username).First(&userRes)
	// TODO 1分钟输错几次让他锁定几分钟
	if util.EncryptPassword(login.Password+userRes.Salt) != userRes.Password {
		response.Failure(c, "账号或密码错误")
		return
	}

	cc := util.CacheClient()

	token := util.GenToken()

	ctx := context.Background()

	// 保存用户登录状态。需要对要存入的数据进行json序列化。取的时候再反序列化。暂时没想到更好的办法
	marshal, err := json.Marshal(userRes)
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	_, err = cc.Set(ctx, token, marshal, time.Duration(time.Second*enumerate.UserLoginTTL)).Result()
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	// 保存用户登录状态

	resp := &dto.UserLoginResponse{
		Token: token,
	}
	response.Success(c, resp)
	return
}
