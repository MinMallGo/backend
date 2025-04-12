package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mall_backend/dto"
	"mall_backend/util"
	"net/http"
)

type RespType int

const (
	done        RespType = 200
	fail        RespType = 400
	UserNoAuth  RespType = 401
	AdminNoAuth RespType = 403
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	Count   int         `json:"count"`
}

// Success 成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  int(done),
		Message: "success",
		Data:    data,
	})
}

func PaginateSuccess(c *gin.Context, count *dto.PaginateCount) {
	c.JSON(http.StatusOK, PaginationResponse{
		Status:  int(done),
		Message: "success",
		Data:    count.Data,
		Page:    count.Page,
		Size:    count.Size,
		Count:   count.Count,
	})
}

// Failure 失败
func Failure(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  int(fail),
		Message: message,
		Data:    struct{}{},
	})
}

func NeedLogin(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  int(UserNoAuth),
		Message: message,
		Data:    struct{}{},
	})
}

func Error(c *gin.Context, err error) {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		Failure(c, util.HandleValidationError(err))
		return
	}

	Failure(c, err.Error())
}
