package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
)

type UserPermissionsController struct{}

/*
	用户权限 API
	/api/w8t/permissions
*/
func (urc UserPermissionsController) API(gin *gin.RouterGroup) {
	perms := gin.Group("permissions")
	perms.Use(
		middleware.Auth(),
	)
	{
		perms.GET("permsList", urc.List)
	}
}

func (urc UserPermissionsController) List(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.UserPermissionService.List()
	})
}
