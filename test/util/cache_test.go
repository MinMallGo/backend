package util

import (
	"context"
	"fmt"
	"log"
	"mall_backend/util"
	"testing"
	"time"
)

func TestCacheConn(t *testing.T) {
	t.Helper()
	util.ConfigInstance("F:\\cxxzoom\\mall\\backend\\config.ini")
	instance, err := util.RedisInstance()
	if err != nil {
		log.Println(err)
		return
	}

	c := instance.Client

	ctx := context.Background()
	key := "xxx"
	c.Set(ctx, key, 123, time.Duration(600*time.Second))
	fmt.Println(c.Get(ctx, key))
}
