package controllers

import (
	"github.com/gin-gonic/gin"
	"prometheus-manager/models/dao"
	"prometheus-manager/pkg/schedule"
	"strconv"
)

type ScheduleController struct {
}

func (sc *ScheduleController) CreateSchedule(ctx *gin.Context) {

	var dutySystem []string
	_ = ctx.ShouldBindJSON(&dutySystem)

	dutyPeriod := ctx.Query("dutyPeriod")

	dutyPeriodInt, _ := strconv.Atoi(dutyPeriod)
	data, err := schedule.CreateAndUpdateDutySystem(dutySystem, dutyPeriodInt)
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

func (sc *ScheduleController) UpdateSchedule(ctx *gin.Context) {

	var dutySystem dao.DutySystem
	_ = ctx.ShouldBindJSON(&dutySystem)

	err := schedule.UpdateDutySystem(dutySystem)
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

func (sc *ScheduleController) SelectDutySystem(ctx *gin.Context) {

	date := ctx.Query("time")

	data, err := schedule.SelectDutySystem(date)
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

func (sc *ScheduleController) SearchUser(ctx *gin.Context) {

	data := schedule.SelectDutyUser()

	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": data,
		"msg":  "查询成功",
	})
}

func (sc *ScheduleController) CreateUser(ctx *gin.Context) {

	var userInfo dao.People

	_ = ctx.ShouldBindJSON(&userInfo)

	err := schedule.CreateDutyUser(userInfo)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3001,
			"data": err.Error(),
			"msg":  "创建失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": nil,
		"msg":  "创建成功",
	})

}

func (sc *ScheduleController) DeleteUser(ctx *gin.Context) {

}

func (sc *ScheduleController) UpdateUser(ctx *gin.Context) {

}

func (sc *ScheduleController) GetUser(ctx *gin.Context) {

	search := ctx.Query("search")
	data, err := schedule.GetDutyUser(search)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3002,
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
