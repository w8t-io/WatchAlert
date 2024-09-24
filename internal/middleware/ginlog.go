package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// GinZapLogger returns a gin.HandlerFunc that logs requests using zap
func GinZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()

		// 计算请求耗时
		latency := end.Sub(start)

		// 获取请求的相关信息
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 记录到 zap
		logger.Info("",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.String("clientIP", clientIP),
			zap.Duration("latency", latency),
			zap.String("errorMessage", errorMessage),
		)
	}
}
