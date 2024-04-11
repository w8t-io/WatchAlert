package initialize

import (
	"watchAlert/alert/eval"
	"watchAlert/config"
	"watchAlert/public/client"
	"watchAlert/public/globals"
	"watchAlert/public/utils/logger"
)

func InitBasic() {

	// 初始化配置
	globals.Config = config.InitConfig()

	// 初始化日志格式
	globals.Logger = logger.InitLogger()

	// 初始化数据库
	globals.DBCli = client.InitDB()

	// 初始化缓存
	globals.RedisCli = client.InitRedis()

	// 启用告警评估携程
	eval.Initialize()

	// 初始化监控分析数据
	InitResource()

	// 初始化权限数据
	InitPermissionsSQL()

	// 初始化角色数据
	InitUserRolesSQL()

}
