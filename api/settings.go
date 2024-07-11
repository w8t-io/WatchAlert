package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type SettingsController struct{}

func (a SettingsController) API(gin *gin.RouterGroup) {
	setting := gin.Group("setting")
	setting.Use(
		middleware.Auth(),
		middleware.Permission(),
	)
	{
		setting.POST("saveSystemSetting", a.Save)
		setting.GET("getSystemSetting", a.Get)
	}
}

func (a SettingsController) Save(ctx *gin.Context) {
	r := new(models.Settings)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SettingService.Save(r)
	})
}

func (a SettingsController) Get(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.SettingService.Get()
	})
}
