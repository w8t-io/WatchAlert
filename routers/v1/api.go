package v1

import (
	"github.com/gin-gonic/gin"
)

func AlertEventMsg(gin *gin.Engine) {

	apiv1 := gin.Group("api/v1")
	{
		// 处理告警 /api/v1/prom/
		prom := apiv1.Group("prom")
		// 接收 Alert
		prom.POST("prometheusAlert", Event.AlertEventMsg)
		// 接收飞书回调
		prom.POST("feiShuEvent", Event.FeiShuEvent)

		// 静默
		/*
			/api/v1/alert/
		*/
		alert := apiv1.Group("alert")
		alert.POST("createSilence", AlertSilence.CreateSilence)

		// 告警规则
		/*
			/api/v1/ruleGroup/
		*/
		ruleGroup := apiv1.Group("ruleGroup")
		ruleGroup.GET("select", RuleGroup.Select)
		ruleGroup.POST("create", RuleGroup.Create)
		ruleGroup.POST("update", RuleGroup.Update)
		ruleGroup.POST("delete", RuleGroup.Delete)
		ruleGroup.GET("getRuleGroup", RuleGroup.GetRuleGroup)
		/*
			/api/v1/ruleGroup/:ruleGroup/rule/
		*/
		rule := ruleGroup.Group(":ruleGroup/rule")
		rule.GET("select", Rule.Select)
		rule.POST("create", Rule.Create)
		rule.POST("delete", Rule.Delete)
		rule.POST("update", Rule.Update)
		rule.GET("getRule", Rule.GetRule)

		// 告警值班
		/*
			/api/v1/dutyManage
		*/
		dutyManage := apiv1.Group("dutyManage")
		dutyManage.POST("create", DutyManage.Create)
		dutyManage.POST("update", DutyManage.Update)
		dutyManage.POST("delete", DutyManage.Delete)
		dutyManage.GET("list", DutyManage.List)
		dutyManage.GET("get", DutyManage.Get)
		/*
			/api/v1/dutyManage/schedule
		*/
		schedule := dutyManage.Group("schedule")
		schedule.POST("create", DutySchedule.Create)
		schedule.POST("update", DutySchedule.Update)
		schedule.GET("select", DutySchedule.Select)
		/*
			/api/v1/dutyManage/user
		*/
		user := dutyManage.Group("user")
		user.POST("create", DutyPeople.Create)
		user.POST("update", DutyPeople.Update)
		user.POST("delete", DutyPeople.Delete)
		user.GET("select", DutyPeople.Search)
		user.GET("getUser", DutyPeople.Get)

		// 通知对象
		/*
			/api/v1/alertNotice
		*/
		notice := apiv1.Group("alertNotice")
		notice.GET("list", AlertNoticeObject.List)
		notice.POST("create", AlertNoticeObject.Create)
		notice.POST("update", AlertNoticeObject.Update)
		notice.POST("delete", AlertNoticeObject.Delete)
		notice.GET("get", AlertNoticeObject.Get)
		notice.GET("getFeiShuChats", AlertNoticeObject.GetFeishuChats)
	}

}
