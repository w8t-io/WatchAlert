package api

import "prometheus-manager/controllers/services"

type ApiGroup struct {
	AlertNoticeObjectController
	AlertSilenceController
	DutyManageController
	DutyPeopleController
	DutyScheduleController
	EventController
	RuleController
	RuleGroupController
}

var ApiGroupApp = new(ApiGroup)

var (
	alertNoticeService  = services.NewInterAlertNoticeService()
	alertSilenceService = services.NewInterAlertSilenceService()
	dutyScheduleService = services.NewInterDutyScheduleService()
	dutyPeopleService   = services.NewInterDutyPeopleService()
	dutyManageService   = services.NewInterDutyManageService()
	eventService        = services.NewInterEventService()
	ruleGroupService    = services.NewInterRuleGroupService()
	ruleService         = services.NewInterRuleService()
)
