package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/models"
)

type RuleController struct {
}

func (rc *RuleController) Create(ctx *gin.Context) {

	var rule models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)

	err := ruleService.Create(rule)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": "401",
			"data": err.Error(),
			"msg":  "failed",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "",
		"msg":  "success",
	})

}

func (rc *RuleController) Update(ctx *gin.Context) {

	var rule models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)

	err := ruleService.Update(rule)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": "401",
			"data": err.Error(),
			"msg":  "failed",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "",
		"msg":  "success",
	})

}

func (rc *RuleController) List(ctx *gin.Context) {

	var rule []models.AlertRule
	_ = ctx.ShouldBindJSON(&rule)

	ruleGroupId := ctx.Query("ruleGroupId")

	data, err := ruleService.List(ruleGroupId)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": "401",
			"data": err.Error(),
			"msg":  "failed",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": data,
		"msg":  "success",
	})

}

func (rc *RuleController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := ruleService.Delete(id)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": "401",
			"data": err.Error(),
			"msg":  "failed",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "",
		"msg":  "success",
	})

}

func (rc *RuleController) Search(ctx *gin.Context) {

	ruleId := ctx.Query("ruleId")
	data := ruleService.Search(ruleId)
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": data,
		"msg":  "success",
	})

}
