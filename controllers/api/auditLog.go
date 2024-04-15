package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/middleware"
	"watchAlert/models"
)

type AuditLogController struct{}

func (ac AuditLogController) API(gin *gin.RouterGroup) {
	auditLog := gin.Group("auditLog")
	auditLog.Use(
		middleware.Cors(),
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		auditLog.GET("listAuditLog", ac.List)
		auditLog.GET("searchAuditLog", ac.Search)
	}
}

func (ac AuditLogController) List(ctx *gin.Context) {
	r := new(models.AuditLogQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return auditLogService.ListAuditLog(r)
	})
}

func (ac AuditLogController) Search(ctx *gin.Context) {
	r := new(models.AuditLogQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return auditLogService.SearchAuditLog(r)
	})
}
