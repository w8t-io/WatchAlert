package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/globals"
	"watchAlert/models"
)

type RuleTmplGroupController struct {
}

func (rtg *RuleTmplGroupController) Create(ctx *gin.Context) {

	var resRT models.RuleTemplateGroup
	_ = ctx.ShouldBindJSON(&resRT)
	err := repo.DBCli.Create(&models.RuleTemplateGroup{}, &resRT)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rtg *RuleTmplGroupController) Delete(ctx *gin.Context) {

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

func (rtg *RuleTmplGroupController) List(ctx *gin.Context) {

	var resRTG []models.RuleTemplateGroup
	globals.DBCli.Model(&models.RuleTemplateGroup{}).Find(&resRTG)
	for k, v := range resRTG {
		var resRT []models.RuleTemplate
		globals.DBCli.Model(&models.RuleTemplate{}).Where("rule_group_name = ?", v.Name).Find(&resRT)
		resRTG[k].Number = len(resRT)
	}
	response.Success(ctx, resRTG, "success")

}

type RuleTmplController struct {
}

func (rtg *RuleTmplController) Create(ctx *gin.Context) {

	var resRT models.RuleTemplate
	_ = ctx.ShouldBindJSON(&resRT)

	err := repo.DBCli.Create(&models.RuleTemplate{}, &resRT)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rtg *RuleTmplController) Delete(ctx *gin.Context) {

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

func (rtg *RuleTmplController) List(ctx *gin.Context) {

	ruleGroupName := ctx.Query("ruleGroupName")

	var resRT []models.RuleTemplate
	globals.DBCli.Model(&models.RuleTemplate{}).Where("rule_group_name = ?", ruleGroupName).Find(&resRT)

	response.Success(ctx, resRT, "success")

}
