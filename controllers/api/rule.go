package api

import (
	"github.com/gin-gonic/gin"
	"prometheus-manager/utils/http"
)

type RuleController struct{}

// Select /api/v1/{:ruleGroup}/rule/select
func (rc *RuleController) Select(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")

	data := ruleService.SelectPromRules(ruleGroup)
	ctx.JSON(200, gin.H{
		"code": 1001,
		"data": data,
		"msg":  "查询成功",
	})

}

// Create /api/v1/{:ruleGroup}/rule/create
func (rc *RuleController) Create(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleBody := ctx.Request.Body
	err := ruleService.CreatePromRule(ruleGroup, ruleBody)
	err = http.PostReloadPrometheus()
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
func (rc *RuleController) Delete(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleName := ctx.Query("ruleName")

	err := ruleService.DeletePromRule(ruleGroup, ruleName)
	err = http.PostReloadPrometheus()
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
func (rc *RuleController) Update(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleName := ctx.Query("ruleName")
	body := ctx.Request.Body

	data, err := ruleService.UpdatePromRule(ruleGroup, ruleName, body)
	err = http.PostReloadPrometheus()
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
func (rc *RuleController) GetRule(ctx *gin.Context) {

	ruleGroup := ctx.Param("ruleGroup")
	ruleName := ctx.Query("ruleName")
	data, err := ruleService.GetPromRuleData(ruleGroup, ruleName)
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
