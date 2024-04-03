package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/models"
	jwtUtils "watchAlert/utils/jwt"
)

type TenantController struct{}

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
