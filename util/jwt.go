package util

import (
	"github.com/golang-jwt/jwt/v5"
	"mall_backend/structure"
	"time"
)

var mySigningKey = []byte("123123")

type CustomClaims struct {
	structure.JWTUserInfo
	jwt.RegisteredClaims
}

func createJwtRegisteredClaims() jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "test",
		Subject:   "somebody",
		ID:        "1",
		Audience:  []string{"somebody_else"},
	}
}

func JWTEncode(userinfo structure.JWTUserInfo) (string, error) {
	claims := CustomClaims{
		userinfo,
		createJwtRegisteredClaims(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(mySigningKey)
}

func JWTDecode(tokenString string) (claims CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return
	}

	if custom, ok := token.Claims.(*CustomClaims); ok {
		return *custom, nil
	}

	err = jwt.ErrInvalidKey

	return
}
