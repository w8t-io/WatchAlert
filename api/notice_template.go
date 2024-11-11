package api

import (
	"github.com/gin-gonic/gin"
	middleware2 "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type NoticeTemplateController struct{}

/*
	通知模版 API
	/api/w8t/noticeTemplate
*/
func (ntc NoticeTemplateController) API(gin *gin.RouterGroup) {
	noticeTemplateA := gin.Group("noticeTemplate")
	noticeTemplateA.Use(
		middleware2.Auth(),
		middleware2.Permission(),
		middleware2.ParseTenant(),
		middleware2.AuditingLog(),
	)
	{
		noticeTemplateA.POST("noticeTemplateCreate", ntc.Create)
		noticeTemplateA.POST("noticeTemplateUpdate", ntc.Update)
		noticeTemplateA.POST("noticeTemplateDelete", ntc.Delete)
	}
	noticeTemplateB := gin.Group("noticeTemplate")
	noticeTemplateB.Use(
		middleware2.Auth(),
		middleware2.Permission(),
		middleware2.ParseTenant(),
	)
	{
		noticeTemplateB.GET("noticeTemplateList", ntc.List)
		noticeTemplateB.GET("searchNoticeTmpl", ntc.Search)
	}
}

func (ntc NoticeTemplateController) Create(ctx *gin.Context) {
	r := new(models.NoticeTemplateExample)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeTmplService.Create(r)
	})
}

func (ntc NoticeTemplateController) Update(ctx *gin.Context) {
	r := new(models.NoticeTemplateExample)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeTmplService.Update(r)
	})
}

func (ntc NoticeTemplateController) Delete(ctx *gin.Context) {
	r := new(models.NoticeTemplateExampleQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeTmplService.Delete(r)
	})
}

func (ntc NoticeTemplateController) List(ctx *gin.Context) {
	r := new(models.NoticeTemplateExampleQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeTmplService.List(r)
	})
}

func (ntc NoticeTemplateController) Search(ctx *gin.Context) {
	r := new(models.NoticeTemplateExampleQuery)
	BindQuery(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeTmplService.Search(r)
	})
}
