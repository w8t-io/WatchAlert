package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
)

type DutyCalendarController struct{}

/*
	值班表 API
	/api/w8t/calendar
*/
func (dcc DutyCalendarController) API(gin *gin.RouterGroup) {
	calendarA := gin.Group("calendar")
	calendarA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		calendarA.POST("calendarCreate", dcc.Create)
		calendarA.POST("calendarUpdate", dcc.Update)
	}

	calendarB := gin.Group("calendar")
	calendarB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		calendarB.GET("calendarSearch", dcc.Select)
	}
}

func (dcc DutyCalendarController) Create(ctx *gin.Context) {

	var dutySchedule models.DutyScheduleCreate
	_ = ctx.ShouldBindJSON(&dutySchedule)

	tid, _ := ctx.Get("TenantID")
	dutySchedule.TenantId = tid.(string)
	data, err := dutyScheduleService.CreateAndUpdateDutySystem(dutySchedule)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}

func (dcc DutyCalendarController) Update(ctx *gin.Context) {

	var dutySchedule models.DutySchedule
	_ = ctx.ShouldBindJSON(&dutySchedule)

	tid, _ := ctx.Get("TenantID")
	dutySchedule.TenantId = tid.(string)
	err := dutyScheduleService.UpdateDutySystem(dutySchedule)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (dcc DutyCalendarController) Select(ctx *gin.Context) {

	tid, _ := ctx.Get("TenantID")
	dutyId := ctx.Query("dutyId")
	date := ctx.Query("time")

	data, err := dutyScheduleService.SelectDutySystem(tid.(string), dutyId, date)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}
