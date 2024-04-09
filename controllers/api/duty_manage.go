package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
	jwtUtils "watchAlert/utils/jwt"
)

type DutyManageController struct{}

func (dmc *DutyManageController) List(ctx *gin.Context) {

	data := dutyManageService.ListDutyManage(ctx)
	response.Success(ctx, data, "success")

}

func (dmc *DutyManageController) Create(ctx *gin.Context) {

	var dutyManage models.DutyManagement
	_ = ctx.ShouldBindJSON(&dutyManage)

	userName := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	dutyManage.CreateBy = userName

	tid, _ := ctx.Get("TenantID")
	dutyManage.TenantId = tid.(string)
	data, err := dutyManageService.CreateDutyManage(dutyManage)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")
}

func (dmc *DutyManageController) Update(ctx *gin.Context) {

	var dutyManage models.DutyManagement
	_ = ctx.ShouldBindJSON(&dutyManage)

	tid, _ := ctx.Get("TenantID")
	dutyManage.TenantId = tid.(string)
	data, err := dutyManageService.UpdateDutyManage(dutyManage)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}

func (dmc *DutyManageController) Delete(ctx *gin.Context) {

	tid, _ := ctx.Get("TenantID")
	id := ctx.Query("id")
	err := dutyManageService.DeleteDutyManage(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (dmc *DutyManageController) Get(ctx *gin.Context) {

	tid, _ := ctx.Get("TenantID")
	id := ctx.Query("id")
	data := dutyManageService.GetDutyManage(tid.(string), id)
	response.Success(ctx, data, "success")

}
