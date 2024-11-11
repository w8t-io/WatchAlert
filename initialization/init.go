package initialization

import (
	"context"
	"golang.org/x/sync/errgroup"
	"watchAlert/alert"
	"watchAlert/config"
	"watchAlert/internal/cache"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
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

	// 导入数据源 Client 到存储池
	importClientPools(ctx)

	if global.Config.Ldap.Enabled {
		// 定时同步LDAP用户任务
		go services.LdapService.SyncUsersCronjob()
	}

}

func importClientPools(ctx *ctx.Context) {
	list, err := ctx.DB.Datasource().List(models.DatasourceQuery{})
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	g := new(errgroup.Group)
	for _, datasource := range list {
		datasource := datasource
		if !*datasource.Enabled {
			continue
		}
		g.Go(func() error {
			err := services.DatasourceService.WithAddClientToProviderPools(datasource)
			if err != nil {
				global.Logger.Sugar().Error("添加到 Client 存储池失败, err: %s", err.Error())
				return err
			}
			return nil
		})
	}
}
