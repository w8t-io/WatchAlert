package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/controllers/services"
)

type ApiGroup struct {
	NoticeController
	DutyController
	DutyCalendarController
	CallbackController
	DatasourceController
	SilenceController
	RuleController
	UserController
	AlertEventController
	UserRoleController
	UserPermissionsController
	NoticeTemplateController
	RuleGroupController
	RuleTmplGroupController
	RuleTmplController
	DashboardInfoController
	TenantController
	DashboardController
	AuditLogController
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
	tenantService        = services.NewInterTenantService()
	dashboardService     = services.NewInterDashboardService()
	auditLogService      = services.NewInterAuditLogService()
)

func Service(ctx *gin.Context, fu func() (interface{}, interface{})) {
	data, err := fu()
	if err != nil {
		response.Fail(ctx, err.(error).Error(), "failed")
		ctx.Abort()
		return
	}
	response.Success(ctx, data, "success")
}

func BindJson(ctx *gin.Context, req interface{}) {
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		ctx.Abort()
		return
	}
}

func BindQuery(ctx *gin.Context, req interface{}) {
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		ctx.Abort()
		return
	}
}
