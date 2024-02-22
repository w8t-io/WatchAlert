package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type RuleGroupController struct {
}

func (rc *RuleGroupController) Create(ctx *gin.Context) {

	var ruleGroup models.RuleGroups
	_ = ctx.ShouldBindJSON(&ruleGroup)

	err := ruleGroupService.Create(ruleGroup)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc *RuleGroupController) Update(ctx *gin.Context) {

	var ruleGroup models.RuleGroups
	_ = ctx.ShouldBindJSON(&ruleGroup)

	err := ruleGroupService.Update(ruleGroup)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc *RuleGroupController) List(ctx *gin.Context) {

	data := ruleGroupService.List()
	response.Success(ctx, data, "success")

}

func (rc *RuleGroupController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := ruleGroupService.Delete(id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}
