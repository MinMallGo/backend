package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
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
	_, err = cc.Set(ctx, token, marshal, time.Duration(time.Second*constants.UserLoginTTL)).Result()
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

func Logout(c *gin.Context) {
	token := c.GetHeader("token")
	ctx := context.Background()
	util.CacheClient().Set(ctx, token, "", time.Duration(time.Second))
	response.Success(c, []string{})
	return
}

func Register(c *gin.Context, register *dto.UserRegister) {
	// 下面要做的验证唯一性，比如手机是否已经注册，然后账号是否已经存在
	res := &model.MmUser{}
	x := util.DBClient().Model(model.MmUser{}).Where("account = ? OR phone = ?", register.Account, register.Phone).First(res)
	if x.RowsAffected != 0 {
		msg := ""
		if res.Account == register.Account {
			msg += "用户名已经被使用"
		}

		if res.Phone == register.Phone {
			msg += "，手机号已经被注册"
		}

		response.Failure(c, msg)
		return
	}
	salt := util.GenSalt()
	password := util.EncryptPassword(register.Password + salt)
	birthday, err := time.Parse("2006-01-02", register.Birthday)
	thirdParty, _ := json.Marshal("")
	if err != nil {
		response.Failure(c, err.Error())
		return
	}
	user := &model.MmUser{
		Account:  register.Account,
		Password: password,
		Salt:     salt,
		Sex:      int32(register.Sex),
		Birthday: birthday,
		Type:     int32(constants.NormalUser),
		Phone:    register.Phone,
		// TODO 注册时的角色管理
		// TODO 这要等到写完角色之后再来联动了
		Role:       "1",
		UniqueCode: util.UserCode(),
		// TODO 用户的第三方信息，如果时通过微信小程序注册过来或者其他三方平台，这里需要保存
		ThirdParty:    string(thirdParty),
		LastLoginTime: util.MinDateTime(),
		CreateTime:    time.Now(),
		UpdateTime:    util.MinDateTime(),
		DeleteTime:    util.MinDateTime(),
	}

	x = util.DBClient().Create(user)
	if x.Error != nil {
		response.Failure(c, x.Error.Error())
		return
	}

	resp := &dto.UserRegisterResponse{
		Account:  register.Account,
		Username: register.Username,
	}

	response.Success(c, resp)
	return
}
