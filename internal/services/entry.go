package services

import "watchAlert/pkg/ctx"

var (
	DatasourceService     InterDatasourceService
	AuditLogService       InterAuditLogService
	DashboardService      InterDashboardService
	DutyManageService     InterDutyManageService
	DutyCalendarService   InterDutyCalendarService
	EventService          InterEventService
	NoticeService         InterNoticeService
	NoticeTmplService     InterNoticeTmplService
	RuleService           InterRuleService
	RuleGroupService      InterRuleGroupService
	RuleTmplService       InterRuleTmplService
	SilenceService        InterSilenceService
	TenantService         InterTenantService
	UserService           InterUserService
	UserRoleService       InterUserRoleService
	AlertService          InterAlertService
	RuleTmplGroupService  InterRuleTmplGroupService
	UserPermissionService InterUserPermissionService
)

func NewServices(ctx *ctx.Context) {
	DatasourceService = newInterDatasourceService(ctx)
	AuditLogService = newInterAuditLogService(ctx)
	DashboardService = newInterDashboardService(ctx)
	DutyManageService = newInterDutyManageService(ctx)
	DutyCalendarService = newInterDutyCalendarService(ctx)
	EventService = newInterEventService(ctx)
	NoticeService = newInterAlertNoticeService(ctx)
	NoticeTmplService = newInterNoticeTmplService(ctx)
	RuleService = newInterRuleService(ctx)
	RuleGroupService = newInterRuleGroupService(ctx)
	RuleTmplService = newInterRuleTmplService(ctx)
	RuleTmplGroupService = newInterRuleTmplGroupService(ctx)
	SilenceService = newInterSilenceService(ctx)
	TenantService = newInterTenantService(ctx)
	UserService = newInterUserService(ctx)
	UserRoleService = newInterUserRoleService(ctx)
	AlertService = newInterAlertService(ctx)
	UserPermissionService = newInterUserPermissionService(ctx)
}
