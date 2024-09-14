package initialization

import (
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
)

var perms []models.UserPermissions

func InitPermissionsSQL(ctx *ctx.Context) {
	var psData []models.UserPermissions

	for _, v := range models.PermissionsInfo() {
		psData = append(psData, v)
	}
	perms = psData

	ctx.DB.DB().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.UserPermissions{})
	ctx.DB.DB().Model(&models.UserPermissions{}).Create(&psData)
}

func InitUserRolesSQL(ctx *ctx.Context) {
	var adminRole models.UserRole
	var db = ctx.DB.DB().Model(&models.UserRole{})

	roles := models.UserRole{
		ID:          "ur-" + cmd.RandId(),
		Name:        "admin",
		Description: "system",
		Permissions: perms,
		CreateAt:    time.Now().Unix(),
	}

	err := db.Where("name = ?", "admin").First(&adminRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ctx.DB.DB().Create(&roles).Error
		}
	} else {
		err = db.Where("name = ?", "admin").Updates(models.UserRole{Permissions: perms}).Error
	}

	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		panic(err)
	}
}
