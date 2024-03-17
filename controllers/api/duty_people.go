package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type DutyPeopleController struct{}

func (dpc *DutyPeopleController) Search(ctx *gin.Context) {

	data := dutyPeopleService.SelectDutyUser()

	response.Success(ctx, data, "success")

}

func (dpc *DutyPeopleController) Create(ctx *gin.Context) {

	var userInfo models.People

	_ = ctx.ShouldBindJSON(&userInfo)

	data, err := dutyPeopleService.CreateDutyUser(userInfo)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}

func (dpc *DutyPeopleController) Delete(ctx *gin.Context) {

	userId := ctx.Query("userId")
	err := dutyPeopleService.DeleteDutyUser(userId)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (dpc *DutyPeopleController) Update(ctx *gin.Context) {

	var userInfo models.People
	_ = ctx.ShouldBindJSON(&userInfo)

	data, err := dutyPeopleService.UpdateDutyUser(userInfo)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}

func (dpc *DutyPeopleController) Get(ctx *gin.Context) {

	search := ctx.Query("search")
	data, err := dutyPeopleService.GetDutyUser(search)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")
	
}
