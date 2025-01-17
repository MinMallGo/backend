package util

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
)

var fieldMaP = map[string]string{
	"Account":        "账号",
	"username":       "用户名",
	"RepeatPassword": "重复密码",
	"Password":       "密码",
	"Sex":            "性别",
	"Birthday":       "出生日期",
	"Phone":          "手机号",
}

// ValidatorRegister 注册自定义验证规则
func ValidatorRegister() {
	// 获取 Gin 的默认校验器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义校验规则
		err := v.RegisterValidation("regexpPhone", RegexpPhone)
		if err != nil {
			panic("failed to register regexpPhone validator")
		}
		log.Println("register regexpPhone validator SUCCESS")
	}
}

// HandleValidationError 接管验证器错误
func HandleValidationError(err error) string {
	message := ""
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for index, fe := range validationErrors {
			message += translateError(fe)
			if index < len(validationErrors)-1 {
				message += ", "
			}
		}
	} else {
		message += "unknown " + err.Error()
	}

	return message
}

// translateError 错误翻译器（自定义消息）
func translateError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s是必填项", getField(fe.Field()))
	case "min":
		return fmt.Sprintf("%s最小值或者最小长度是%s", getField(fe.Field()), fe.Param())
	case "max":
		return fmt.Sprintf("%s最大值或者最大长度是%s", getField(fe.Field()), fe.Param())
	case "len":
		return fmt.Sprintf("%s长度必须是%s", getField(fe.Field()), fe.Param())
	case "regexpPhone":
		return fmt.Sprintf("%s格式不正确", getField(fe.Field()))
	case "datetime":
		return fmt.Sprintf("%s必须是有效日期，格式为%s", getField(fe.Field()), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s必须是数字且区间为%s", getField(fe.Field()), fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s与%s必须相同", getField(fe.Field()), getField(fe.Param()))
	default:
		return fmt.Sprintf("%s格式不符合要求", getField(fe.Field()))
	}
}

func getField(field string) string {
	if str, ok := fieldMaP[field]; ok {
		return str
	}
	return field
}

// RegexpPhone 注册一些自定义的验证器规则
func RegexpPhone(fl validator.FieldLevel) bool {
	regex := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(regex)
	return reg.MatchString(fl.Field().String())
}
