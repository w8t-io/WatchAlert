package api

import (
	"watchAlert/controllers/services"
)

type ApiGroup struct {
	AlertNoticeObjectController
	DutyManageController
	DutyPeopleController
	DutyScheduleController
	EventController
	AlertDataSourceController
	AlertSilenceController
	RuleController
	UserController
	AlertCurEventController
	AlertHisEventController
	UserRoleController
	UserPermissionsController
	AlertNoticeTemplateController
	RuleGroupController
	RuleTmplGroupController
	RuleTmplController
	DashboardInfoController
}

var ApiGroupApp = new(ApiGroup)

var (
	alertNoticeService   = services.NewInterAlertNoticeService()
	dutyScheduleService  = services.NewInterDutyScheduleService()
	dutyPeopleService    = services.NewInterDutyPeopleService()
	dutyManageService    = services.NewInterDutyManageService()
	ruleService          = services.NewInterRuleService()
	dataSourceService    = services.NewInterAlertDataSourceService()
	alertSilenceService  = services.NewInterAlertSilenceService()
	alertCurEventService = services.NewInterAlertCurEventService()
	alertHisEventService = services.NewInterAlertHisEventService()
	ruleGroupService     = services.NewInterRuleGroupService()
)
