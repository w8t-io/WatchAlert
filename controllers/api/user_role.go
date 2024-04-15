package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
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

	var userRole models.UserRole
	_ = ctx.ShouldBindJSON(&userRole)

	pString, _ := json.Marshal(userRole.PermissionsJson)
	userRole.Permissions = string(pString)

	userRole.ID = "ur-" + cmd.RandId()
	userRole.CreateAt = time.Now().Unix()
	err := repo.DBCli.Create(models.UserRole{}, &userRole)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (urc UserRoleController) Update(ctx *gin.Context) {

	var userRole models.UserRole
	_ = ctx.ShouldBindJSON(&userRole)

	pString, _ := json.Marshal(userRole.PermissionsJson)
	userRole.Permissions = string(pString)

	err := repo.DBCli.Updates(repo.Updates{
		Table:   &models.UserRole{},
		Where:   []interface{}{"id = ?", userRole.ID},
		Updates: userRole,
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (urc UserRoleController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := repo.DBCli.Delete(repo.Delete{
		Table: models.UserRole{},
		Where: []interface{}{"id = ?", id},
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (urc UserRoleController) List(ctx *gin.Context) {

	var data []models.UserRole

	err := globals.DBCli.Model(&models.UserRole{}).Find(&data).Error

	for k, v := range data {
		var pJson []models.UserPermissions
		_ = json.Unmarshal([]byte(v.Permissions), &pJson)
		data[k].PermissionsJson = pJson
	}

	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, data, "success")

}
