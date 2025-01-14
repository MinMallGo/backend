package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mall_backend/response"
	"mall_backend/util"
)

func Login(c *gin.Context) {
}

func Ping(c *gin.Context) {
	jwt, _ := util.JWTDecode(c.GetHeader("Authorization"))
	if err := validator.New().Struct(jwt); err != nil {
		response.Failure(c, err.Error())
	}
}
