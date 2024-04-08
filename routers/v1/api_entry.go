package v1

import (
	"watchAlert/controllers/api"
)

var (
	AlertNoticeObject = api.ApiGroupApp.AlertNoticeObjectController
	AlertSilence      = api.ApiGroupApp.AlertSilenceController
	AlertDatasource   = api.ApiGroupApp.AlertDataSourceController
	DutyManage        = api.ApiGroupApp.DutyManageController
	DutySchedule      = api.ApiGroupApp.DutyScheduleController
	Event             = api.ApiGroupApp.EventController
	Rule              = api.ApiGroupApp.RuleController
	Auth              = api.ApiGroupApp.UserController
	AlertCurEvent     = api.ApiGroupApp.AlertCurEventController
	AlertHisEvent     = api.ApiGroupApp.AlertHisEventController
	Role              = api.ApiGroupApp.UserRoleController
	Permissions       = api.ApiGroupApp.UserPermissionsController
	NoticeTemplate    = api.ApiGroupApp.AlertNoticeTemplateController
	RuleGroup         = api.ApiGroupApp.RuleGroupController
	RuleTmplGroup     = api.ApiGroupApp.RuleTmplGroupController
	RuleTmpl          = api.ApiGroupApp.RuleTmplController
	DashboardInfo     = api.ApiGroupApp.DashboardInfoController
	Tenant            = api.ApiGroupApp.TenantController
	Dashboard         = api.ApiGroupApp.DashboardController
)
