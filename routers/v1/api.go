package v1

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/middleware/jwt"
	"watchAlert/middleware/permission"
)

func AlertEventMsg(gin *gin.Engine) {

	apiV1 := gin.Group("api")
	{

		/*
			不需要鉴权
			/api/system
		*/
		system := apiV1.Group("system")
		{
			system.POST("register", Auth.Register)
			system.POST("login", Auth.Login)
			system.GET("checkUser", Auth.CheckUser)
			// 接收飞书回调
			system.POST("feiShuEvent", Event.FeiShuEvent)
			system.GET("checkNoticeStatus", AlertNoticeObject.CheckNoticeStatus)
		}

		/*
			需要鉴权
			/api/w8t
		*/
		w8t := apiV1.Group("w8t")
		w8t.Use(
			middleware.JwtAuth(),
			permission.Permission(),
		)
		{
			/*
				用户
				/api/w8t/user
			*/
			user := w8t.Group("user")
			{
				user.POST("userUpdate", Auth.Update)
				user.GET("userList", Auth.List)
				user.POST("userDelete", Auth.Delete)
				user.POST("userChangePass", Auth.ChangePass)
				user.GET("searchDutyUser", Auth.SearchDutyUser)
			}

			/*
				角色
				/api/w8t/role
			*/
			role := w8t.Group("role")
			{
				role.POST("roleCreate", Role.Create)
				role.POST("roleDelete", Role.Delete)
				role.GET("roleList", Role.List)
			}

			perms := w8t.Group("permissions")
			{
				perms.GET("permsList", Permissions.List)
			}

			/*
				告警静默
				/api/w8t/silence
			*/
			silence := w8t.Group("silence")
			{
				silence.POST("silenceCreate", AlertSilence.Create)
				silence.POST("silenceUpdate", AlertSilence.Update)
				silence.POST("silenceDelete", AlertSilence.Delete)
				silence.GET("silenceList", AlertSilence.List)
			}

			/*
				告警规则
				/api/w8t/rule
			*/
			rule := w8t.Group("rule")
			{
				rule.POST("ruleCreate", Rule.Create)
				rule.POST("ruleUpdate", Rule.Update)
				rule.POST("ruleDelete", Rule.Delete)
				rule.GET("ruleList", Rule.List)
				rule.GET("ruleSearch", Rule.Search)
			}

			/*
				排班管理
				/api/w8t/dutyManage
			*/
			dutyManage := w8t.Group("dutyManage")
			{
				dutyManage.POST("dutyManageCreate", DutyManage.Create)
				dutyManage.POST("dutyManageUpdate", DutyManage.Update)
				dutyManage.POST("dutyManageDelete", DutyManage.Delete)
				dutyManage.GET("dutyManageList", DutyManage.List)
				dutyManage.GET("dutyManageSearch", DutyManage.Get)
			}

			/*
				值班表
				/api/w8t/calendar
			*/
			schedule := w8t.Group("calendar")
			{
				schedule.POST("calendarCreate", DutySchedule.Create)
				schedule.POST("calendarUpdate", DutySchedule.Update)
				schedule.GET("calendarSearch", DutySchedule.Select)
			}

			/*
				通知对象
				/api/w8t/notice
			*/
			notice := w8t.Group("notice")
			{
				notice.GET("noticeList", AlertNoticeObject.List)
				notice.POST("noticeCreate", AlertNoticeObject.Create)
				notice.POST("noticeUpdate", AlertNoticeObject.Update)
				notice.POST("noticeDelete", AlertNoticeObject.Delete)
				notice.GET("noticeSearch", AlertNoticeObject.Get)
			}
			
			/*
				通知模版
				/api/w8t/noticeTemplate
			*/
			noticeTemplate := w8t.Group("noticeTemplate")
			{
				noticeTemplate.GET("noticeTemplateList", NoticeTemplate.List)
				noticeTemplate.POST("noticeTemplateCreate", NoticeTemplate.Create)
				noticeTemplate.POST("noticeTemplateUpdate", NoticeTemplate.Update)
				noticeTemplate.POST("noticeTemplateDelete", NoticeTemplate.Delete)
			}

			/*
				通知对象
				/api/w8t/datasource
			*/
			alert := w8t.Group("datasource")
			{
				alert.POST("dataSourceCreate", AlertDatasource.Create)
				alert.POST("dataSourceUpdate", AlertDatasource.Update)
				alert.POST("dataSourceDelete", AlertDatasource.Delete)
				alert.GET("dataSourceList", AlertDatasource.List)
				alert.GET("dataSourceSearch", AlertDatasource.Search)
			}

			/*
				告警事件
				/api/w8t/event
			*/
			event := w8t.Group("event")
			{
				event.GET("curEvent", AlertCurEvent.List)
				event.GET("hisEvent", AlertHisEvent.List)
			}

		}

	}

}
