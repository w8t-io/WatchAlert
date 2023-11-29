package v1

import (
	"github.com/gin-gonic/gin"
	"prometheus-manager/controllers"
	"prometheus-manager/controllers/rules"
)

var (
	ec  controllers.EventController
	ac  controllers.AlertController
	rc  rules.RuleController
	rgc rules.RuleGroupController
)

func AlertEventMsg(gin *gin.Engine) {

	apiv1 := gin.Group("api/v1")
	{
		// 处理告警 /api/v1/prom/
		prom := apiv1.Group("prom")
		// 接收 Alert
		prom.POST("prometheusAlert", ec.AlertEventMsg)
		// 接收飞书回调
		prom.POST("feiShuEvent", ec.FeiShuEvent)

		// 静默 /api/v1/alert/
		alert := apiv1.Group("alert")
		alert.POST("createSilence", ac.CreateSilence)

		// 告警规则组 /api/v1/ruleGroup/
		ruleGroup := apiv1.Group("ruleGroup")
		ruleGroup.GET("select", rgc.Select)
		ruleGroup.POST("create", rgc.Create)
		ruleGroup.POST("update", rgc.Update)
		ruleGroup.POST("delete", rgc.Delete)
		ruleGroup.GET("getRuleGroup", rgc.GetRuleGroup)

		// 告警规则 /api/v1/ruleGroup/:ruleGroup/rule/
		rule := ruleGroup.Group(":ruleGroup/rule")
		rule.GET("select", rc.Select)
		rule.POST("create", rc.Create)
		rule.POST("delete", rc.Delete)
		rule.POST("update", rc.Update)
		rule.GET("getRule", rc.GetRule)
	}

}
