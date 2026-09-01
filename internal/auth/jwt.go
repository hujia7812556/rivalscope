// Package auth 提供 JWT token 的签发与解析。
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT 载荷,携带登录身份(自包含,无需查库)。
type Claims struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

// JWT 签名与解析器。
type JWT struct {
	key         []byte
	expireHours int
}

// NewJWT 创建 JWT 实例;expireHours 为 token 有效期(小时)。
func NewJWT(key string, expireHours int) *JWT {
	return &JWT{key: []byte(key), expireHours: expireHours}
}

// GenToken 为指定用户签发 token。
func (j *JWT) GenToken(userId, username, nickname string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:   userId,
		Username: username,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.expireHours) * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.key)
}

// ParseToken 校验并解析 token,返回载荷;支持带或不带 "Bearer " 前缀。
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名算法: %v", t.Header["alg"])
		}
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token 无效")
	}
	return claims, nil
}
