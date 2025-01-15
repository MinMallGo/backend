package util

import (
	"log"
	"mall_backend/util"
	"testing"
)

// 开始之前要初始化config连接
func TestDBConn(t *testing.T) {
	t.Helper()
	util.ConfigInstance("F:\\cxxzoom\\mall\\backend\\config.ini")
	if db, err := util.DBInstance(); err != nil {
		log.Println(db, err)
	} else {
		log.Println(db, err)
	}
}
