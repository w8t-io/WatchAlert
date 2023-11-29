package rules

import (
	"github.com/gin-gonic/gin"
	"prometheus-manager/pkg/rules"
)

type RuleGroupController struct{}

// Select /api/v1/{:ruleGroup}/select
func (rg *RuleGroupController) Select(ctx *gin.Context) {

	groupData, err := rules.SelectRuleGroup()
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": "2001",
			"data": err.Error(),
			"msg":  "查询失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "1000",
		"data": groupData,
		"msg":  "查询成功",
	})

}

// Create /api/v1/{:ruleGroup}/create
func (rg *RuleGroupController) Create(ctx *gin.Context) {

	body := ctx.Request.Body
	err := rules.CreateRuleGroup(body)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": "2002",
			"data": err.Error(),
			"msg":  "创建失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "1000",
		"data": nil,
		"msg":  "创建成功",
	})
}

// Update /api/v1/{:ruleGroup}/update?ruleGroupName=test
func (rg *RuleGroupController) Update(ctx *gin.Context) {

	ruleGroupName := ctx.Query("ruleGroupName")

	body := ctx.Request.Body

	err := rules.UpdateRuleGroup(ruleGroupName, body)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": "2003",
			"data": err.Error(),
			"msg":  "更新失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "1000",
		"data": nil,
		"msg":  "更新成功",
	})
}

// Delete /api/v1/{:ruleGroup}/delete?ruleGroupName=test
func (rg *RuleGroupController) Delete(ctx *gin.Context) {

	ruleGroupName := ctx.Query("ruleGroupName")
	err := rules.DeleteRuleGroup(ruleGroupName)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": "2004",
			"data": err.Error(),
			"msg":  "删除失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "1000",
		"data": nil,
		"msg":  "删除成功",
	})

}

// GetRuleGroup /api/v1/{:ruleGroup}/getRuleGroup?ruleGroupName=test
func (rg *RuleGroupController) GetRuleGroup(ctx *gin.Context) {

	ruleGroupName := ctx.Query("ruleGroupName")
	data, err := rules.GetRuleGroup(ruleGroupName)

	if err != nil {
		ctx.JSON(500, gin.H{
			"code": "2005",
			"data": err.Error(),
			"msg":  "查询失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "1000",
		"data": data,
		"msg":  "查询成功",
	})
}
