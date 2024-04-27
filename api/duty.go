package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	jwtUtils "watchAlert/pkg/utils/jwt"
)

type DutyController struct{}

/*
	排班管理 API
	/api/w8t/dutyManage
*/
func (dc DutyController) API(gin *gin.RouterGroup) {
	dutyManageA := gin.Group("dutyManage")
	dutyManageA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		dutyManageA.POST("dutyManageCreate", dc.Create)
		dutyManageA.POST("dutyManageUpdate", dc.Update)
		dutyManageA.POST("dutyManageDelete", dc.Delete)
	}

	dutyManageB := gin.Group("dutyManage")
	dutyManageB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		dutyManageB.GET("dutyManageList", dc.List)
		dutyManageB.GET("dutyManageSearch", dc.Get)
	}
}

func (dc DutyController) List(ctx *gin.Context) {
	r := new(models.DutyManagementQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.List(r)
	})
}

func (dc DutyController) Create(ctx *gin.Context) {
	r := new(models.DutyManagement)
	BindJson(ctx, r)

	userName := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	r.CreateBy = userName

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.Create(r)
	})
}

func (dc DutyController) Update(ctx *gin.Context) {
	r := new(models.DutyManagement)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.Update(r)
	})
}

func (dc DutyController) Delete(ctx *gin.Context) {
	r := new(models.DutyManagementQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.Delete(r)
	})
}

func (dc DutyController) Get(ctx *gin.Context) {
	r := new(models.DutyManagementQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyManageService.Get(r)
	})
}
