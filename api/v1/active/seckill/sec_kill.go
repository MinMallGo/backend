package seckill

import (
	"github.com/gin-gonic/gin"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/service"
)

// CreateSecKill 秒杀活动
func CreateSecKill(c *gin.Context) {
	param := &dto.SecKillCreate{}
	if err := c.ShouldBindJSON(param); err != nil {
		response.Failure(c, err.Error())
	}

	service.SecKillCreate(c, param)
}
