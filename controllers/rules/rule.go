package rules

import (
	"github.com/gin-gonic/gin"
	promRule "prometheus-manager/pkg/rules"
	"prometheus-manager/utils"
)

type RuleController struct{}

// Select /api/v1/{:ruleGroup}/rule/select
func (r *RuleController) Select(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")

	data := promRule.SelectPromRules(ruleGroup)
	ctx.JSON(200, gin.H{
		"code": 1001,
		"data": data,
		"msg":  "查询成功",
	})

}

// Create /api/v1/{:ruleGroup}/rule/create
func (r *RuleController) Create(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleBody := ctx.Request.Body
	err := promRule.CreatePromRule(ruleGroup, ruleBody)
	err = utils.PostReloadPrometheus()
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 1002,
			"data": err.Error(),
			"msg":  "创建失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 1000,
		"data": nil,
		"msg":  "创建成功",
	})
}

// Delete /api/v1/{:ruleGroup}/rule/delete?ruleName=test
func (r *RuleController) Delete(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleName := ctx.Query("ruleName")

	err := promRule.DeletePromRule(ruleGroup, ruleName)
	err = utils.PostReloadPrometheus()
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 1003,
			"data": err.Error(),
			"msg":  "删除失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 1000,
		"data": nil,
		"msg":  "删除成功",
	})

}

// Update /api/v1/{:ruleGroup}/rule/update?ruleName=test
func (r *RuleController) Update(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleName := ctx.Query("ruleName")
	body := ctx.Request.Body

	data, err := promRule.UpdatePromRule(ruleGroup, ruleName, body)
	err = utils.PostReloadPrometheus()
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 1004,
			"data": err.Error(),
			"msg":  "更新失败",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 1000,
		"data": data,
		"msg":  "更新成功",
	})
}

// GetRule /api/v1/{:ruleGroup}/rule/getRuleInfo?ruleName=test
func (r *RuleController) GetRule(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleName := ctx.Query("ruleName")
	data, err := promRule.GetPromRuleData(ruleGroup, ruleName)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 1005,
			"data": err.Error(),
			"msg":  "获取失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 1000,
		"data": data,
		"msg":  "获取成功",
	})
	return

}
