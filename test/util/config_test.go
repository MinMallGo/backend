package util

import (
	"fmt"
	"mall_backend/util"
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Helper()
	util.ConfigInstance("F:\\cxxzoom\\mall\\backend\\config.ini")
	fmt.Println(util.GetRedis())
	fmt.Println(util.GetDatabase())
}
