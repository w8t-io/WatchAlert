package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type RuleController struct {
}

func (rc *RuleController) Create(ctx *gin.Context) {

	var rule models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)

	err := ruleService.Create(rule)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc *RuleController) Update(ctx *gin.Context) {

	var rule models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)

	err := ruleService.Update(rule)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc *RuleController) List(ctx *gin.Context) {

	var rule []models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)

	ruleGroupId := ctx.Query("ruleGroupId")

	data, err := ruleService.List(ruleGroupId)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, data, "success")

}

func (rc *RuleController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := ruleService.Delete(id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc *RuleController) Search(ctx *gin.Context) {

	ruleId := ctx.Query("ruleId")
	data := ruleService.Search(ruleId)
	response.Success(ctx, data, "success")

}
