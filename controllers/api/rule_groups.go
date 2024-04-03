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

	tid, _ := ctx.Get("TenantID")
	ruleGroup.TenantId = tid.(string)
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

	tid, _ := ctx.Get("TenantID")
	ruleGroup.TenantId = tid.(string)
	err := ruleGroupService.Update(ruleGroup)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (rc *RuleGroupController) List(ctx *gin.Context) {

	data := ruleGroupService.List(ctx)
	response.Success(ctx, data, "success")

}

func (rc *RuleGroupController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	tid, _ := ctx.Get("TenantID")
	err := ruleGroupService.Delete(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}
