package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/models"
	jwtUtils "watchAlert/utils/jwt"
)

type DutyManageController struct{}

func (dmc *DutyManageController) List(ctx *gin.Context) {

	data := dutyManageService.ListDutyManage()
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": data,
		"msg":  "查询成功",
	})

}

func (dmc *DutyManageController) Create(ctx *gin.Context) {

	var dutyManage models.DutyManagement
	_ = ctx.ShouldBindJSON(&dutyManage)

	userName := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	dutyManage.CreateBy = userName

	data, err := dutyManageService.CreateDutyManage(dutyManage)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3010,
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

func (dmc *DutyManageController) Update(ctx *gin.Context) {

	var dutyManage models.DutyManagement
	_ = ctx.ShouldBindJSON(&dutyManage)

	data, err := dutyManageService.UpdateDutyManage(dutyManage)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3011,
			"data": err.Error(),
			"msg":  "更新失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": data,
		"msg":  "更新成功",
	})

}

func (dmc *DutyManageController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := dutyManageService.DeleteDutyManage(id)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3012,
			"data": err.Error(),
			"msg":  "删除失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": "",
		"msg":  "删除成功",
	})

}

func (dmc *DutyManageController) Get(ctx *gin.Context) {

	id := ctx.Query("id")
	data := dutyManageService.GetDutyManage(id)
	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": data,
		"msg":  "查询成功",
	})

}
