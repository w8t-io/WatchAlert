package middleware

import (
	"github.com/gin-gonic/gin"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/response"
	jwtUtils "watchAlert/pkg/utils/jwt"
)

func Auth() gin.HandlerFunc {

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
		code, ok := jwtUtils.IsTokenValid(ctx.DO(), tokenStr)
		if !ok {
			if code == 401 {
				response.TokenFail(context)
				context.Abort()
				return
			}
		}

	}

}
