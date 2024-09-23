package initialization

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/middleware"
	"watchAlert/internal/routers"
	"watchAlert/internal/routers/v1"
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

func InitRoute() {
	global.Logger.Sugar().Info("服务启动")

	mode := global.Config.Server.Mode
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	ginEngine := gin.New()

	ginEngine.Use(
		// 启用CORS中间件
		middleware.Cors(),
		// 自定义请求日志格式
		GinZapLogger(global.Logger),
		gin.Recovery(),
		//gin.LoggerWithFormatter(middleware.RequestLoggerFormatter),
	)
	allRouter(ginEngine)

	err := ginEngine.Run(":" + global.Config.Server.Port)
	if err != nil {
		global.Logger.Sugar().Error("服务启动失败:", err)
		return
	}
}

func allRouter(engine *gin.Engine) {

	routers.HealthCheck(engine)
	v1.Router(engine)

}
