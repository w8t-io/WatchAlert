package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type AlertNoticeObjectController struct{}

func (ano *AlertNoticeObjectController) List(ctx *gin.Context) {

	object := alertNoticeService.SearchNoticeObject(ctx)
	response.Success(ctx, object, "success")

}

func (ano *AlertNoticeObjectController) Create(ctx *gin.Context) {

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

func (ano *AlertNoticeObjectController) Update(ctx *gin.Context) {

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

func (ano *AlertNoticeObjectController) Delete(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	tid, _ := ctx.Get("TenantID")
	err := alertNoticeService.DeleteNoticeObject(tid.(string), uuid)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (ano *AlertNoticeObjectController) Get(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	tid, _ := ctx.Get("TenantID")
	object := alertNoticeService.GetNoticeObject(tid.(string), uuid)
	response.Success(ctx, object, "success")

}

func (ano *AlertNoticeObjectController) CheckNoticeStatus(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	tid, _ := ctx.Get("TenantID")
	status := alertNoticeService.CheckNoticeObjectStatus(tid.(string), uuid)
	response.Success(ctx, status, "success")

}
