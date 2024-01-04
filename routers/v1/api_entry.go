package v1

import "watchAlert/controllers/api"

var (
	AlertNoticeObject = api.ApiGroupApp.AlertNoticeObjectController
	AlertSilence      = api.ApiGroupApp.AlertSilenceController
	DutyManage        = api.ApiGroupApp.DutyManageController
	DutyPeople        = api.ApiGroupApp.DutyPeopleController
	DutySchedule      = api.ApiGroupApp.DutyScheduleController
	Event             = api.ApiGroupApp.EventController
	Rule              = api.ApiGroupApp.RuleController
	RuleGroup         = api.ApiGroupApp.RuleGroupController
)
