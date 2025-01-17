package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mall_backend/util"
	"net/http"
)

type RespType int

const (
	done RespType = 400200
	fail RespType = 400400
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  int(done),
		Message: "success",
		Data:    data,
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

func Error(c *gin.Context, err error) {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		Failure(c, util.HandleValidationError(err))
		return
	}

	Failure(c, err.Error())
}
