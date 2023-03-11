package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type myClaims struct {
	UserID               string `json:"user_id"`
	jwt.RegisteredClaims        // 内嵌申明字段
}

var mySecret = []byte("秘密")

func ParseToken(tokenString string) (*myClaims, error) {
	// 解析token
	// 如果自定义的Claims结构体需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claims进行类型断言
	if claims, ok := token.Claims.(*myClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
