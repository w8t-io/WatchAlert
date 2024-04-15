package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
	jwtUtils "watchAlert/public/utils/jwt"
)

type SilenceController struct{}

/*
	告警静默 API
	/api/w8t/silence
*/
func (sc SilenceController) API(gin *gin.RouterGroup) {
	silenceA := gin.Group("silence")
	silenceA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		silenceA.POST("silenceCreate", sc.Create)
		silenceA.POST("silenceUpdate", sc.Update)
		silenceA.POST("silenceDelete", sc.Delete)
	}

	silenceB := gin.Group("silence")
	silenceB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		silenceB.GET("silenceList", sc.List)

	}
}

func (sc SilenceController) Create(ctx *gin.Context) {

	var silence models.AlertSilences
	_ = ctx.ShouldBindJSON(&silence)

	tid, _ := ctx.Get("TenantID")
	silence.TenantId = tid.(string)
	user := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	silence.CreateBy = user

	err := alertSilenceService.CreateAlertSilence(silence)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (sc SilenceController) Update(ctx *gin.Context) {

	var silence models.AlertSilences
	_ = ctx.ShouldBindJSON(&silence)

	tid, _ := ctx.Get("TenantID")
	silence.TenantId = tid.(string)
	user := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	silence.UpdateBy = user

	data, err := alertSilenceService.UpdateAlertSilence(silence)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}

func (sc SilenceController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	tid, _ := ctx.Get("TenantID")
	err := alertSilenceService.DeleteAlertSilence(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (sc SilenceController) List(ctx *gin.Context) {

	data, err := alertSilenceService.ListAlertSilence(ctx)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}
