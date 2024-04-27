package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type UserRoleController struct{}

/*
	用户角色 API
	/api/w8t/role
*/
func (urc UserRoleController) API(gin *gin.RouterGroup) {
	roleA := gin.Group("role")
	roleA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		roleA.POST("roleCreate", urc.Create)
		roleA.POST("roleUpdate", urc.Update)
		roleA.POST("roleDelete", urc.Delete)
	}

	roleB := gin.Group("role")
	roleB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		roleB.GET("roleList", urc.List)
	}
}

func (urc UserRoleController) Create(ctx *gin.Context) {
	r := new(models.UserRole)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.Create(r)
	})
}

func (urc UserRoleController) Update(ctx *gin.Context) {
	r := new(models.UserRole)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.Update(r)
	})
}

func (urc UserRoleController) Delete(ctx *gin.Context) {
	r := new(models.UserRoleQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.Delete(r)
	})
}

func (urc UserRoleController) List(ctx *gin.Context) {
	r := new(models.UserRoleQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserRoleService.List(r)
	})
}
