package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type NoticeTemplateController struct{}

/*
	通知模版 API
	/api/w8t/noticeTemplate
*/
func (ntc NoticeTemplateController) API(gin *gin.RouterGroup) {
	noticeTemplateA := gin.Group("noticeTemplate")
	noticeTemplateA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		noticeTemplateA.POST("noticeTemplateCreate", ntc.Create)
		noticeTemplateA.POST("noticeTemplateUpdate", ntc.Update)
		noticeTemplateA.POST("noticeTemplateDelete", ntc.Delete)
	}
	noticeTemplateB := gin.Group("noticeTemplate")
	noticeTemplateB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		noticeTemplateB.GET("noticeTemplateList", ntc.List)
	}
}

func (ntc NoticeTemplateController) Create(ctx *gin.Context) {

	var tmpl models.NoticeTemplateExample
	_ = ctx.ShouldBindJSON(&tmpl)

	tmpl.Id = "nt-" + cmd.RandId()
	err := repo.DBCli.Create(&models.NoticeTemplateExample{}, &tmpl)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (ntc NoticeTemplateController) Update(ctx *gin.Context) {

	var tmpl models.NoticeTemplateExample
	_ = ctx.ShouldBindJSON(&tmpl)

	err := repo.DBCli.Updates(repo.Updates{
		Table:   &models.NoticeTemplateExample{},
		Where:   []interface{}{"id = ?", tmpl.Id},
		Updates: &tmpl,
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (ntc NoticeTemplateController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")

	err := repo.DBCli.Delete(repo.Delete{
		Table: &models.NoticeTemplateExample{},
		Where: []interface{}{"id = ?", id},
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (ntc NoticeTemplateController) List(ctx *gin.Context) {

	var templates []models.NoticeTemplateExample
	globals.DBCli.Model(&models.NoticeTemplateExample{}).Find(&templates)
	response.Success(ctx, templates, "success")

}
