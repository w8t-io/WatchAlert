package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type AlertNoticeTemplateController struct{}

func (ant *AlertNoticeTemplateController) Create(ctx *gin.Context) {

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

func (ant *AlertNoticeTemplateController) Update(ctx *gin.Context) {

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

func (ant *AlertNoticeTemplateController) Delete(ctx *gin.Context) {

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

func (ant *AlertNoticeTemplateController) List(ctx *gin.Context) {

	var templates []models.NoticeTemplateExample
	globals.DBCli.Model(&models.NoticeTemplateExample{}).Find(&templates)
	response.Success(ctx, templates, "success")

}
