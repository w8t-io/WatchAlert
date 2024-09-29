package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type DashboardController struct{}

/*
	仪表盘 API
	/api/w8t/dashboard
*/
func (dc DashboardController) API(gin *gin.RouterGroup) {
	dashboardA := gin.Group("dashboard")
	dashboardA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		dashboardA.POST("createDashboard", dc.Create)
		dashboardA.POST("updateDashboard", dc.Update)
		dashboardA.POST("deleteDashboard", dc.Delete)
		dashboardA.POST("createFolder", dc.CreateFolder)
		dashboardA.POST("updateFolder", dc.UpdateFolder)
		dashboardA.POST("deleteFolder", dc.DeleteFolder)
	}
	dashboardB := gin.Group("dashboard")
	dashboardB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		dashboardB.GET("getDashboard", dc.Get)
		dashboardB.GET("listDashboard", dc.List)
		dashboardB.GET("searchDashboard", dc.Search)
		dashboardB.GET("listFolder", dc.ListFolder)
		dashboardB.GET("getFolder", dc.GetFolder)
		dashboardB.GET("listGrafanaDashboards", dc.ListGrafanaDashboards)
		dashboardB.GET("getDashboardFullUrl", dc.GetDashboardFullUrl)
	}
}

func (dc DashboardController) List(ctx *gin.Context) {
	r := new(models.DashboardQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.List(r)
	})
}

func (dc DashboardController) Get(ctx *gin.Context) {
	r := new(models.DashboardQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.Get(r)
	})
}

func (dc DashboardController) Create(ctx *gin.Context) {
	r := new(models.Dashboard)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.Create(r)
	})
}

func (dc DashboardController) Update(ctx *gin.Context) {
	r := new(models.Dashboard)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.Update(r)
	})
}

func (dc DashboardController) Delete(ctx *gin.Context) {
	r := new(models.DashboardQuery)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.Delete(r)
	})
}

func (dc DashboardController) Search(ctx *gin.Context) {
	r := new(models.DashboardQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.Search(r)
	})
}

func (dc DashboardController) ListFolder(ctx *gin.Context) {
	r := new(models.DashboardFolders)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.ListFolder(r)
	})
}

func (dc DashboardController) GetFolder(ctx *gin.Context) {
	r := new(models.DashboardFolders)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetFolder(r)
	})
}

func (dc DashboardController) CreateFolder(ctx *gin.Context) {
	r := new(models.DashboardFolders)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.CreateFolder(r)
	})
}

func (dc DashboardController) UpdateFolder(ctx *gin.Context) {
	r := new(models.DashboardFolders)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.UpdateFolder(r)
	})
}

func (dc DashboardController) DeleteFolder(ctx *gin.Context) {
	r := new(models.DashboardFolders)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.DeleteFolder(r)
	})
}

func (dc DashboardController) ListGrafanaDashboards(ctx *gin.Context) {
	r := new(models.DashboardFolders)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.ListGrafanaDashboards(r)
	})
}

func (dc DashboardController) GetDashboardFullUrl(ctx *gin.Context) {
	r := new(models.DashboardFolders)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DashboardService.GetDashboardFullUrl(r)
	})
}
