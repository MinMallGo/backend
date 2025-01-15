package util

import (
	"crypto/md5"
	"encoding/hex"
	rand2 "math/rand/v2"
)

// 李少说不要用jwt，改成token和redis存储，方便服务端控制

// GenToken 就返回一个长度是 27位的字符串
func GenToken() string {
	var s []byte
	for i := 0; i <= 26; i++ {
		s = append(s, byte('a'+rand2.IntN(26)))
	}

	return string(s)
}

// EncryptPassword 密码加密
func EncryptPassword(password string) string {
	x := md5.Sum([]byte(password))
	return hex.EncodeToString(x[:])
}
