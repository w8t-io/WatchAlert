package tenant

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
)

func ParseTenantInfo() gin.HandlerFunc {
	// 从HTTP头部获取TenantID并存储到上下文中，可以提高代码的可维护性、可重用性、安全性和性能，同时也使得错误处理和业务逻辑的实现更加高效和灵活。
	return func(context *gin.Context) {
		tid := context.Request.Header.Get("TenantID")
		if tid == "" {
			response.Fail(context, "租户ID不能为空", "failed")
			context.Abort()
			return
		}
		// TODO 判断租户ID是否是有效的
		context.Set("TenantID", tid)
	}
}
