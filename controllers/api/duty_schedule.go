package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type DutyScheduleController struct{}

func (sc *DutyScheduleController) Create(ctx *gin.Context) {

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

func (sc *DutyScheduleController) Update(ctx *gin.Context) {

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

func (sc *DutyScheduleController) Select(ctx *gin.Context) {

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
