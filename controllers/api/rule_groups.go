package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
)

type RuleGroupController struct{}

/*
	规则组 API
	/api/w8t/ruleGroup
*/
func (rc RuleGroupController) API(gin *gin.RouterGroup) {
	ruleGroupA := gin.Group("ruleGroup")
	ruleGroupA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		ruleGroupA.POST("ruleGroupCreate", rc.Create)
		ruleGroupA.POST("ruleGroupUpdate", rc.Update)
		ruleGroupA.POST("ruleGroupDelete", rc.Delete)
	}
	ruleGroupB := gin.Group("ruleGroup")
	ruleGroupB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		ruleGroupB.GET("ruleGroupList", rc.List)
	}
}

func (rc RuleGroupController) Create(ctx *gin.Context) {

	var ruleGroup models.RuleGroups
	_ = ctx.ShouldBindJSON(&ruleGroup)

	tid, _ := ctx.Get("TenantID")
	ruleGroup.TenantId = tid.(string)
	err := ruleGroupService.Create(ruleGroup)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc RuleGroupController) Update(ctx *gin.Context) {

	var ruleGroup models.RuleGroups
	_ = ctx.ShouldBindJSON(&ruleGroup)

	tid, _ := ctx.Get("TenantID")
	ruleGroup.TenantId = tid.(string)
	err := ruleGroupService.Update(ruleGroup)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc RuleGroupController) List(ctx *gin.Context) {

	data := ruleGroupService.List(ctx)
	response.Success(ctx, data, "success")

}

func (rc RuleGroupController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	tid, _ := ctx.Get("TenantID")
	err := ruleGroupService.Delete(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}
