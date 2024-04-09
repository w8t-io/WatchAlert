package initialize

import (
	"github.com/gin-gonic/gin"
	"watchAlert/globals"
	"watchAlert/middleware/cors"
	"watchAlert/middleware/requestLoggerFormatter"
	"watchAlert/routers"
	"watchAlert/routers/v1"
)

func InitRoute() {
	globals.Logger.Sugar().Info("服务启动")
	ginEngine := gin.New()

	var mode string
	if globals.Config.Server.Mode != "" {
		mode = globals.Config.Server.Mode
	} else {
		mode = "debug"
	}

	gin.SetMode(mode)
	ginEngine.Use(
		// 启用CORS中间件
		cors.Cors(),
		// 自定义请求日志格式
		gin.LoggerWithFormatter(requestLoggerFormatter.CustomLogFormatter),
	)
	allRouter(ginEngine)

	err := ginEngine.Run(":" + globals.Config.Server.Port)
	if err != nil {
		globals.Logger.Sugar().Error("服务启动失败:", err)
		return
	}
}

func allRouter(engine *gin.Engine) {

	routers.HealthCheck(engine)
	v1.AlertEventMsg(engine)

}
