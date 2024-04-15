package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
	"watchAlert/public/globals"
)

type RuleTmplGroupController struct{}

/*
	规则模版组 API
	/api/w8t/ruleTmplGroup
*/
func (rtgc RuleTmplGroupController) API(gin *gin.RouterGroup) {
	ruleTmplGroupA := gin.Group("ruleTmplGroup")
	ruleTmplGroupA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		ruleTmplGroupA.POST("ruleTmplGroupCreate", rtgc.Create)
		ruleTmplGroupA.POST("ruleTmplGroupDelete", rtgc.Delete)
	}

	ruleTmplGroupB := gin.Group("ruleTmplGroup")
	ruleTmplGroupB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		ruleTmplGroupB.GET("ruleTmplGroupList", rtgc.List)
	}
}

func (rtgc RuleTmplGroupController) Create(ctx *gin.Context) {

	var resRT models.RuleTemplateGroup
	_ = ctx.ShouldBindJSON(&resRT)
	err := repo.DBCli.Create(&models.RuleTemplateGroup{}, &resRT)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rtgc RuleTmplGroupController) Delete(ctx *gin.Context) {

	tmplGroupName := ctx.Query("tmplGroupName")
	err := repo.DBCli.Delete(repo.Delete{
		Table: &models.RuleTemplateGroup{},
		Where: []interface{}{"name = ?", tmplGroupName},
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rtgc RuleTmplGroupController) List(ctx *gin.Context) {

	var resRTG []models.RuleTemplateGroup
	globals.DBCli.Model(&models.RuleTemplateGroup{}).Find(&resRTG)
	for k, v := range resRTG {
		var resRT []models.RuleTemplate
		globals.DBCli.Model(&models.RuleTemplate{}).Where("rule_group_name = ?", v.Name).Find(&resRT)
		resRTG[k].Number = len(resRT)
	}
	response.Success(ctx, resRTG, "success")

}

type RuleTmplController struct{}

/*
	规则模版 API
	/api/w8t/ruleTmpl
*/
func (rtc RuleTmplController) API(gin *gin.RouterGroup) {
	ruleTmplA := gin.Group("ruleTmpl")
	ruleTmplA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		ruleTmplA.POST("ruleTmplCreate", rtc.Create)
		ruleTmplA.POST("ruleTmplDelete", rtc.Delete)
	}

	ruleTmplB := gin.Group("ruleTmpl")
	ruleTmplB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		ruleTmplB.GET("ruleTmplList", rtc.List)
	}
}

func (rtc RuleTmplController) Create(ctx *gin.Context) {

	var resRT models.RuleTemplate
	_ = ctx.ShouldBindJSON(&resRT)

	err := repo.DBCli.Create(&models.RuleTemplate{}, &resRT)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rtc RuleTmplController) Delete(ctx *gin.Context) {

	ruleName := ctx.Query("ruleName")
	err := repo.DBCli.Delete(repo.Delete{
		Table: &models.RuleTemplate{},
		Where: []interface{}{"rule_name = ?", ruleName},
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rtc RuleTmplController) List(ctx *gin.Context) {

	ruleGroupName := ctx.Query("ruleGroupName")

	var resRT []models.RuleTemplate
	globals.DBCli.Model(&models.RuleTemplate{}).Where("rule_group_name = ?", ruleGroupName).Find(&resRT)

	response.Success(ctx, resRT, "success")

}
