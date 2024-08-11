package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type MonitorSSLController struct{}

func (m MonitorSSLController) API(gin *gin.RouterGroup) {
	mon := gin.Group("monitor")
	mon.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		mon.POST("createMon", m.create)
		mon.POST("updateMon", m.update)
		mon.POST("deleteMon", m.delete)
		mon.GET("listMon", m.list)
		mon.GET("getMon", m.get)
	}
}

func (m MonitorSSLController) create(ctx *gin.Context) {
	r := new(models.MonitorSSLRule)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.MonitorService.Create(r)
	})
}

func (m MonitorSSLController) update(ctx *gin.Context) {
	r := new(models.MonitorSSLRule)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.MonitorService.Update(r)
	})
}

func (m MonitorSSLController) delete(ctx *gin.Context) {
	r := new(models.MonitorSSLRuleQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.MonitorService.Delete(r)
	})
}

func (m MonitorSSLController) list(ctx *gin.Context) {
	r := new(models.MonitorSSLRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.MonitorService.List(r)
	})
}

func (m MonitorSSLController) get(ctx *gin.Context) {
	r := new(models.MonitorSSLRuleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.MonitorService.Get(r)
	})
}
