package util

import (
	constants "mall_backend/constant"
	"time"
)

func MinDateTime() time.Time {
	x, _ := time.Parse("2006-01-02", constants.MinDateTime)
	return x
}
