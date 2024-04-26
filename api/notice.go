package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type NoticeController struct{}

/*
	通知对象 API
	/api/w8t/sender
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
		noticeB.GET("noticeSearch", nc.Search)
	}
}

func (nc NoticeController) List(ctx *gin.Context) {
	r := new(models.NoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.List(r)
	})
}

func (nc NoticeController) Create(ctx *gin.Context) {
	r := new(models.AlertNotice)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Create(r)
	})
}

func (nc NoticeController) Update(ctx *gin.Context) {
	r := new(models.AlertNotice)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Update(r)
	})
}

func (nc NoticeController) Delete(ctx *gin.Context) {
	r := new(models.NoticeQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Delete(r)
	})
}

func (nc NoticeController) Get(ctx *gin.Context) {
	r := new(models.NoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Get(r)
	})

}

func (nc NoticeController) Check(ctx *gin.Context) {
	r := new(models.NoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Check(r)
	})
}

func (nc NoticeController) Search(ctx *gin.Context) {
	r := new(models.NoticeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.NoticeService.Search(r)
	})
}
