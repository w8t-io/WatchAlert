package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/models"
)

type DashboardController struct{}

func (dc DashboardController) ListDashboard(ctx *gin.Context) {
	tid, _ := ctx.Get("TenantID")
	Service(ctx, func() (interface{}, interface{}) {
		return dashboardService.ListDashboard(tid.(string))
	})
}

func (dc DashboardController) GetDashboard(ctx *gin.Context) {
	r := new(models.DashboardQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return dashboardService.GetDashboard(r)
	})
}

func (dc DashboardController) CreateDashboard(ctx *gin.Context) {
	r := new(models.Dashboard)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return dashboardService.CreateDashboard(r)
	})
}

func (dc DashboardController) UpdateDashboard(ctx *gin.Context) {
	r := new(models.Dashboard)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return dashboardService.UpdateDashboard(r)
	})
}

func (dc DashboardController) DeleteDashboard(ctx *gin.Context) {
	r := new(models.DashboardQuery)
	BindJson(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return dashboardService.DeleteDashboard(r)
	})
}

func (dc DashboardController) SearchDashboard(ctx *gin.Context) {
	r := new(models.DashboardQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return dashboardService.SearchDashboard(r)
	})
}
