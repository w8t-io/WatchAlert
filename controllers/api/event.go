package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
)

type AlertEventController struct{}

/*
	告警事件 API
	/api/w8t/event
*/
func (aec AlertEventController) API(gin *gin.RouterGroup) {
	event := gin.Group("event")
	event.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		event.GET("curEvent", aec.ListCurrentEvent)
		event.GET("hisEvent", aec.ListHistoryEvent)
	}
}

func (aec AlertEventController) ListCurrentEvent(ctx *gin.Context) {

	dsType := ctx.Query("dsType")
	tid, _ := ctx.Get("TenantID")
	data, err := alertCurEventService.List(tid.(string), dsType)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, data, "success")

}

func (aec AlertEventController) ListHistoryEvent(ctx *gin.Context) {

	data, err := alertHisEventService.List(ctx)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}
