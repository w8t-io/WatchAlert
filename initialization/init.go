package initialization

import (
	"context"
	"watchAlert/alert"
	"watchAlert/config"
	"watchAlert/internal/cache"
	"watchAlert/internal/global"
	"watchAlert/internal/repo"
	"watchAlert/internal/services"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/logger"
)

func InitBasic() {

	// 初始化配置
	global.Config = config.InitConfig()

	// 初始化日志格式
	global.Logger = logger.InitLogger()

	dbRepo := repo.NewRepoEntry()
	rCache := cache.NewEntryCache()
	ctx := ctx.NewContext(context.Background(), dbRepo, rCache)

	services.NewServices(ctx)

	// 启用告警评估携程
	alert.Initialize(ctx)

	// 初始化监控分析数据
	InitResource(ctx)

	// 初始化权限数据
	InitPermissionsSQL(ctx)

	// 初始化角色数据
	InitUserRolesSQL(ctx)

	if global.Config.Ldap.Enabled {
		// 定时同步LDAP用户任务
		go services.LdapService.SyncUsersCronjob()
	}

}
