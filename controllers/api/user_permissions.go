package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
	"watchAlert/public/globals"
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

	var data []models.UserPermissions
	globals.DBCli.Model(&models.UserPermissions{}).Find(&data)
	response.Success(ctx, data, "success")

}
