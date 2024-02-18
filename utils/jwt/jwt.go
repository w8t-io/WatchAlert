package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
	"watchAlert/globals"
)

// 把签发的秘钥 抛出来
var stSignKey = []byte(viper.GetString("jwt.WatchAlert"))

// JwtCustomClaims 注册声明是JWT声明集的结构化版本，仅限于注册声明名称
type JwtCustomClaims struct {
	ID             string
	Name           string
	StandardClaims jwt.StandardClaims
}

const (
	// TokenType Token 类型
	TokenType = "bearer"
	// AppGuardName 颁发者
	AppGuardName = "WatchAlert"
)

func (j JwtCustomClaims) Valid() error {
	return nil
}

// GenerateToken 生成Token
func GenerateToken(id string, name string) (string, error) {

	// 初始化
	iJwtCustomClaims := JwtCustomClaims{
		ID:   id,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + globals.Config.Jwt.Expire,
			IssuedAt:  time.Now().Unix(),
			Issuer:    AppGuardName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustomClaims)
	return token.SignedString(stSignKey)
}

// ParseToken 解析token
func parseToken(tokenStr string) (JwtCustomClaims, error) {

	iJwtCustomClaims := JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &iJwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return stSignKey, nil
	})

	if err == nil && !token.Valid {
		err = errors.New("invalid Token")
	}
	return iJwtCustomClaims, err

}

func IsTokenValid(tokenStr string) (JwtCustomClaims, bool) {

	token, err := parseToken(tokenStr)
	if err != nil {
		return token, false
	}

	// 发布者校验
	if token.StandardClaims.Issuer != AppGuardName {
		return token, false
	}

	// 校验过期时间
	ok := token.StandardClaims.VerifyExpiresAt(time.Now().Unix(), false)
	if !ok {
		return token, false
	}

	return token, true

}

func GetUser(tokenStr string) string {

	if tokenStr == "" {
		return ""
	}

	tokenStr = tokenStr[len(TokenType)+1:]
	token, err := parseToken(tokenStr)
	if err != nil {
		return ""
	}
	return token.Name

}

func GetUserID(tokenStr string) string {

	if tokenStr == "" {
		return ""
	}

	token, err := parseToken(tokenStr)
	if err != nil {
		return ""
	}

	return token.ID

}