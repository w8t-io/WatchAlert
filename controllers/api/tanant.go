package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/middleware"
	"watchAlert/models"
	jwtUtils "watchAlert/public/utils/jwt"
)

type TenantController struct{}

/*
	租户 API
	/api/w8t/tenant
*/
func (tc TenantController) API(gin *gin.RouterGroup) {
	tenantA := gin.Group("tenant")
	tenantA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.AuditingLog(),
	)
	{
		tenantA.POST("createTenant", tc.CreateTenant)
		tenantA.POST("updateTenant", tc.UpdateTenant)
		tenantA.POST("deleteTenant", tc.DeleteTenant)
	}

	tenantB := gin.Group("tenant")
	tenantB.Use(
		middleware.Auth(),
		middleware.Permission(),
	)
	{
		tenantB.GET("getTenantList", tc.GetTenantList)
	}
}

func (tc TenantController) CreateTenant(ctx *gin.Context) {
	r := new(models.Tenant)
	BindJson(ctx, r)
	r.CreateBy = jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	Service(ctx, func() (interface{}, interface{}) {
		return tenantService.CreateTenant(r)
	})
}

func (tc TenantController) UpdateTenant(ctx *gin.Context) {
	r := new(models.Tenant)
	BindJson(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return tenantService.UpdateTenant(r)
	})
}

func (tc TenantController) DeleteTenant(ctx *gin.Context) {
	r := new(models.Tenant)
	BindJson(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return tenantService.DeleteTenant(r)
	})
}

func (tc TenantController) GetTenantList(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return tenantService.ListTenant()
	})
}

func (tc TenantController) SearchTenant(ctx *gin.Context) {
	// TODO
}
