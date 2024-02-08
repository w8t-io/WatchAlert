package middleware

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	jwtUtils "watchAlert/utils/jwt"
)

func JwtAuth() gin.HandlerFunc {

	return func(context *gin.Context) {
		// 获取 Token
		tokenStr := context.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(context)
			context.Abort()
			return
		}

		// Bearer Token, 获取 Token 值
		tokenStr = tokenStr[len(jwtUtils.TokenType)+1:]

		// 校验 Token
		token, ok := jwtUtils.IsTokenValid(tokenStr)
		if !ok {
			response.TokenFail(context)
			context.Abort()
			return
		}

		// 发布者校验
		if token.StandardClaims.Issuer != jwtUtils.AppGuardName {
			response.TokenFail(context)
			context.Abort()
			return
		}

		context.Set("token", token)
		context.Set("id", token.ID)

	}

}
