package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
	jwtUtils "watchAlert/public/utils/jwt"
)

type DutyController struct{}

/*
	排班管理 API
	/api/w8t/dutyManage
*/
func (dc DutyController) API(gin *gin.RouterGroup) {
	dutyManageA := gin.Group("dutyManage")
	dutyManageA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		dutyManageA.POST("dutyManageCreate", dc.Create)
		dutyManageA.POST("dutyManageUpdate", dc.Update)
		dutyManageA.POST("dutyManageDelete", dc.Delete)
	}

	dutyManageB := gin.Group("dutyManage")
	dutyManageB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		dutyManageB.GET("dutyManageList", dc.List)
		dutyManageB.GET("dutyManageSearch", dc.Get)
	}
}

func (dc DutyController) List(ctx *gin.Context) {

	data := dutyManageService.ListDutyManage(ctx)
	response.Success(ctx, data, "success")

}

func (dc DutyController) Create(ctx *gin.Context) {

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

func (dc DutyController) Update(ctx *gin.Context) {

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

func (dc DutyController) Delete(ctx *gin.Context) {

	tid, _ := ctx.Get("TenantID")
	id := ctx.Query("id")
	err := dutyManageService.DeleteDutyManage(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (dc DutyController) Get(ctx *gin.Context) {

	tid, _ := ctx.Get("TenantID")
	id := ctx.Query("id")
	data := dutyManageService.GetDutyManage(tid.(string), id)
	response.Success(ctx, data, "success")

}
