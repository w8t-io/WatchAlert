package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/dao"
)

type DutyPeopleController struct{}

func (dpc *DutyPeopleController) Search(ctx *gin.Context) {

	data := dutyPeopleService.SelectDutyUser()

	ctx.JSON(200, gin.H{
		"code": 3000,
		"data": data,
		"msg":  "查询成功",
	})
}

func (dpc *DutyPeopleController) Create(ctx *gin.Context) {

	var userInfo dao.People

	_ = ctx.ShouldBindJSON(&userInfo)

	data, err := dutyPeopleService.CreateDutyUser(userInfo)
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
		"data": data,
		"msg":  "创建成功",
	})

}

func (dpc *DutyPeopleController) Delete(ctx *gin.Context) {

	userId := ctx.Query("userId")
	err := dutyPeopleService.DeleteDutyUser(userId)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3004,
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

func (dpc *DutyPeopleController) Update(ctx *gin.Context) {

	var userInfo dao.People
	_ = ctx.ShouldBindJSON(&userInfo)

	data, err := dutyPeopleService.UpdateDutyUser(userInfo)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": 3003,
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

func (dpc *DutyPeopleController) Get(ctx *gin.Context) {

	search := ctx.Query("search")
	data, err := dutyPeopleService.GetDutyUser(search)
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
