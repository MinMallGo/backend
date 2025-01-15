package util

import (
	"fmt"
	"mall_backend/structure"
	"mall_backend/util"
	"testing"
)

func TestJWTEncode(t *testing.T) {
	token, err := util.JWTEncode(structure.JWTUserInfo{
		Username: "小王",
		UserID:   1,
		Role:     "Admin",
		DeviceID: "xxx_pc",
	})
	if err != nil {
		t.Fatal("encode err:", err)
	}

	decode, err := util.JWTDecode(token)
	if err != nil {
		t.Fatal("decode err:", err)
	}

	fmt.Println(decode)
}
