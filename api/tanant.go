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
		tenantA.POST("addUsersToTenant", tc.AddUsersToTenant)
		tenantA.POST("delUsersOfTenant", tc.DelUsersOfTenant)
		tenantA.POST("changeTenantUserRole", tc.ChangeTenantUserRole)
	}

	tenantB := gin.Group("tenant")
	tenantB.Use(
		middleware.Auth(),
		middleware.Permission(),
	)
	{
		tenantB.GET("getTenantList", tc.List)
		tenantB.GET("getTenant", tc.Get)
		tenantB.GET("getUsersForTenant", tc.GetUsersForTenant)
	}
}

func (tc TenantController) Create(ctx *gin.Context) {
	r := new(models.Tenant)
	BindJson(ctx, r)

	token := ctx.Request.Header.Get("Authorization")
	r.CreateBy = jwtUtils.GetUser(token)
	r.UserId = jwtUtils.GetUserID(token)
	if r.UserId == "" {
		r.UserId = "admin"
	}

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
	r := new(models.TenantQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.List(r)
	})
}

func (tc TenantController) Get(ctx *gin.Context) {
	r := new(models.TenantQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.Get(r)
	})
}

func (tc TenantController) Search(ctx *gin.Context) {
	// TODO
}

func (tc TenantController) AddUsersToTenant(ctx *gin.Context) {
	r := new(models.TenantLinkedUsers)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.AddUsersToTenant(r)
	})
}

func (tc TenantController) DelUsersOfTenant(ctx *gin.Context) {
	r := new(models.TenantQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.DelUsersOfTenant(r)
	})
}

func (tc TenantController) GetUsersForTenant(ctx *gin.Context) {
	r := new(models.TenantQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.GetUsersForTenant(r)
	})
}

func (tc TenantController) ChangeTenantUserRole(ctx *gin.Context) {
	r := new(models.ChangeTenantUserRole)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TenantService.ChangeTenantUserRole(r)
	})
}
