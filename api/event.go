package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type AlertEventController struct{}

/*
	告警事件 API
	/api/w8t/event
*/
func (e AlertEventController) API(gin *gin.RouterGroup) {
	event := gin.Group("event")
	event.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		event.GET("curEvent", e.ListCurrentEvent)
		event.GET("hisEvent", e.ListHistoryEvent)
	}
}

func (e AlertEventController) ListCurrentEvent(ctx *gin.Context) {
	r := new(models.AlertCurEventQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListCurrentEvent(r)
	})
}

func (e AlertEventController) ListHistoryEvent(ctx *gin.Context) {
	r := new(models.AlertHisEventQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.EventService.ListHistoryEvent(r)
	})
}
