package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
)

type NoticeController struct{}

/*
	通知对象 API
	/api/w8t/notice
*/
func (nc NoticeController) API(gin *gin.RouterGroup) {
	noticeA := gin.Group("notice")
	noticeA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		noticeA.POST("noticeCreate", nc.Create)
		noticeA.POST("noticeUpdate", nc.Update)
		noticeA.POST("noticeDelete", nc.Delete)
	}

	noticeB := gin.Group("notice")
	noticeB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		noticeB.GET("noticeList", nc.List)
		noticeB.GET("noticeSearch", nc.Get)
		noticeB.GET("searchNotice", nc.Search)
	}
}

func (nc NoticeController) List(ctx *gin.Context) {

	object := alertNoticeService.ListNoticeObject(ctx)
	response.Success(ctx, object, "success")

}

func (nc NoticeController) Create(ctx *gin.Context) {

	var alertNotice models.AlertNotice
	_ = ctx.ShouldBindJSON(&alertNotice)
	tid, _ := ctx.Get("TenantID")
	alertNotice.TenantId = tid.(string)
	object, err := alertNoticeService.CreateNoticeObject(alertNotice)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, object, "success")

}

func (nc NoticeController) Update(ctx *gin.Context) {

	var alertNotice models.AlertNotice
	_ = ctx.ShouldBindJSON(&alertNotice)
	tid, _ := ctx.Get("TenantID")
	alertNotice.TenantId = tid.(string)
	object, err := alertNoticeService.UpdateNoticeObject(alertNotice)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, object, "success")

}

func (nc NoticeController) Delete(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	tid, _ := ctx.Get("TenantID")
	err := alertNoticeService.DeleteNoticeObject(tid.(string), uuid)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (nc NoticeController) Get(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	tid, _ := ctx.Get("TenantID")
	object := alertNoticeService.GetNoticeObject(tid.(string), uuid)
	response.Success(ctx, object, "success")

}

func (nc NoticeController) CheckNoticeStatus(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	tid, _ := ctx.Get("TenantID")
	status := alertNoticeService.CheckNoticeObjectStatus(tid.(string), uuid)
	response.Success(ctx, status, "success")

}

func (nc NoticeController) Search(ctx *gin.Context) {
	r := new(models.NoticeQuery)
	BindQuery(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return noticeService.Search(r)
	})
}
