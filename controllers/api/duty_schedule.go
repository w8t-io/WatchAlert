package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/models"
)

type DutyScheduleController struct{}

func (sc *DutyScheduleController) Create(ctx *gin.Context) {

	var dutySchedule models.DutyScheduleCreate
	_ = ctx.ShouldBindJSON(&dutySchedule)

	data, err := dutyScheduleService.CreateAndUpdateDutySystem(dutySchedule)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3003,
			"data": err.Error(),
			"msg":  "创建失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": data,
		"msg":  "创建成功",
	})
}

func (sc *DutyScheduleController) Update(ctx *gin.Context) {

	var dutySchedule models.DutySchedule
	_ = ctx.ShouldBindJSON(&dutySchedule)

	err := dutyScheduleService.UpdateDutySystem(dutySchedule)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3004,
			"data": err.Error(),
			"msg":  "更新失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": nil,
		"msg":  "更新成功",
	})
}

func (sc *DutyScheduleController) Select(ctx *gin.Context) {

	dutyId := ctx.Query("dutyId")
	date := ctx.Query("time")

	data, err := dutyScheduleService.SelectDutySystem(dutyId, date)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3003,
			"data": err.Error(),
			"msg":  "查询失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": data,
		"msg":  "查询成功",
	})

}
