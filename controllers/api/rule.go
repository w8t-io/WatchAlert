package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
)

type RuleController struct{}

/*
	告警规则 API
	/api/w8t/rule
*/
func (rc RuleController) API(gin *gin.RouterGroup) {
	ruleA := gin.Group("rule")
	ruleA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		ruleA.POST("ruleCreate", rc.Create)
		ruleA.POST("ruleUpdate", rc.Update)
		ruleA.POST("ruleDelete", rc.Delete)
	}
	ruleB := gin.Group("rule")
	ruleB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		ruleB.GET("ruleList", rc.List)
		ruleB.GET("ruleSearch", rc.Search)
	}
}

func (rc RuleController) Create(ctx *gin.Context) {

	var rule models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)
	tid, _ := ctx.Get("TenantID")
	rule.TenantId = tid.(string)
	err := ruleService.Create(rule)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc RuleController) Update(ctx *gin.Context) {

	var rule models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)
	tid, _ := ctx.Get("TenantID")
	rule.TenantId = tid.(string)
	err := ruleService.Update(rule)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc RuleController) List(ctx *gin.Context) {

	var rule []models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)

	ruleGroupId := ctx.Query("ruleGroupId")
	tid, _ := ctx.Get("TenantID")

	data, err := ruleService.List(tid.(string), ruleGroupId)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, data, "success")

}

func (rc RuleController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	tid, _ := ctx.Get("TenantID")
	err := ruleService.Delete(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc RuleController) Search(ctx *gin.Context) {

	ruleId := ctx.Query("ruleId")
	tid, _ := ctx.Get("TenantID")
	data := ruleService.Search(tid.(string), ruleId)
	response.Success(ctx, data, "success")

}
