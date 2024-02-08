package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

type UserRoleController struct {
}

func (urc *UserRoleController) Create(ctx *gin.Context) {

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

func (urc *UserRoleController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := repo.DBCli.Delete(repo.Delete{
		Table: models.UserRole{},
		Where: []string{"id = ?", id},
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (urc *UserRoleController) List(ctx *gin.Context) {

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
