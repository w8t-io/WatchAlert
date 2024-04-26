package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	jwtUtils "watchAlert/pkg/utils/jwt"
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
		tenantA.POST("createTenant", tc.Create)
		tenantA.POST("updateTenant", tc.Update)
		tenantA.POST("deleteTenant", tc.Delete)
	}

	tenantB := gin.Group("tenant")
	tenantB.Use(
		middleware.Auth(),
		middleware.Permission(),
	)
	{
		tenantB.GET("getTenantList", tc.List)
	}
}

func (tc TenantController) Create(ctx *gin.Context) {
	r := new(models.Tenant)
	BindJson(ctx, r)
	r.CreateBy = jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Create(r)
	})
}

func (tc TenantController) Update(ctx *gin.Context) {
	r := new(models.Tenant)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Update(r)
	})
}

func (tc TenantController) Delete(ctx *gin.Context) {
	r := new(models.TenantQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Delete(r)
	})
}

func (tc TenantController) List(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.List()
	})
}

func (tc TenantController) Search(ctx *gin.Context) {
	// TODO
}
