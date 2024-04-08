package initialize

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

var perms []models.UserPermissions

func InitPermissionsSQL() {

	permissions := []models.UserPermissions{
		{
			Key: "UserRegister",
			API: "/api/system/register",
		},
		{
			Key: "UserList",
			API: "/api/w8t/user/userList",
		},
		{
			Key: "UserUpdate",
			API: "/api/w8t/user/userUpdate",
		},
		{
			Key: "UserDelete",
			API: "/api/w8t/user/userDelete",
		},
		{
			Key: "UserChangePass",
			API: "/api/w8t/user/userChangePass",
		},
		{
			Key: "SearchDutyUser",
			API: "/api/w8t/user/searchDutyUser",
		},
		{
			Key: "RoleCreate",
			API: "/api/w8t/role/roleCreate",
		},
		{
			Key: "RoleUpdate",
			API: "/api/w8t/role/roleUpdate",
		},
		{
			Key: "RoleDelete",
			API: "/api/w8t/role/roleDelete",
		},
		{
			Key: "RoleList",
			API: "/api/w8t/role/roleList",
		},
		{
			Key: "SilenceCreate",
			API: "/api/w8t/silence/silenceCreate",
		},
		{
			Key: "SilenceUpdate",
			API: "/api/w8t/silence/silenceUpdate",
		},
		{
			Key: "SilenceDelete",
			API: "/api/w8t/silence/silenceDelete",
		},
		{
			Key: "SilenceList",
			API: "/api/w8t/silence/silenceList",
		},
		{
			Key: "RuleCreate",
			API: "/api/w8t/rule/ruleCreate",
		},
		{
			Key: "RuleUpdate",
			API: "/api/w8t/rule/ruleUpdate",
		},
		{
			Key: "RuleDelete",
			API: "/api/w8t/rule/ruleDelete",
		},
		{
			Key: "RuleList",
			API: "/api/w8t/rule/ruleList",
		},
		{
			Key: "RuleSearch",
			API: "/api/w8t/rule/RuleSearch",
		},
		{
			Key: "DutyManageCreate",
			API: "/api/w8t/dutyManage/dutyManageCreate",
		},
		{
			Key: "DutyManageUpdate",
			API: "/api/w8t/dutyManage/dutyManageUpdate",
		},
		{
			Key: "DutyManageDelete",
			API: "/api/w8t/dutyManage/dutyManageDelete",
		},
		{
			Key: "DutyManageList",
			API: "/api/w8t/dutyManage/dutyManageList",
		},
		{
			Key: "DutyManageSearch",
			API: "/api/w8t/dutyManage/dutyManageSearch",
		},
		{
			Key: "DutyScheduleCreate",
			API: "/api/w8t/calendar/calendarCreate",
		},
		{
			Key: "DutyScheduleUpdate",
			API: "/api/w8t/calendar/calendarUpdate",
		},
		{
			Key: "DutyScheduleSearch",
			API: "/api/w8t/calendar/calendarSearch",
		},
		{
			Key: "NoticeCreate",
			API: "/api/w8t/notice/noticeCreate",
		},
		{
			Key: "NoticeUpdate",
			API: "/api/w8t/notice/noticeUpdate",
		},
		{
			Key: "NoticeDelete",
			API: "/api/w8t/notice/noticeDelete",
		},
		{
			Key: "NoticeList",
			API: "/api/w8t/notice/noticeList",
		},
		{
			Key: "NoticeSearch",
			API: "/api/w8t/notice/noticeSearch",
		},
		{
			Key: "DataSourceCreate",
			API: "/api/w8t/datasource/dataSourceCreate",
		},
		{
			Key: "DataSourceUpdate",
			API: "/api/w8t/datasource/dataSourceUpdate",
		},
		{
			Key: "DataSourceDelete",
			API: "/api/w8t/datasource/dataSourceDelete",
		},
		{
			Key: "DataSourceList",
			API: "/api/w8t/datasource/dataSourceList",
		},
		{
			Key: "DataSourceGet",
			API: "/api/w8t/datasource/dataSourceGet",
		},
		{
			Key: "DataSourceSearch",
			API: "/api/w8t/datasource/dataSourceSearch",
		},
		{
			Key: "CurrentEventList",
			API: "/api/w8t/event/curEvent",
		},
		{
			Key: "HistoryEventList",
			API: "/api/w8t/event/hisEvent",
		},
		{
			Key: "PermissionsList",
			API: "/api/w8t/permissions/permsList",
		},
		{
			Key: "NoticeTemplateList",
			API: "/api/w8t/noticeTemplate/noticeTemplateList",
		},
		{
			Key: "NoticeTemplateCreate",
			API: "/api/w8t/noticeTemplate/noticeTemplateCreate",
		},
		{
			Key: "NoticeTemplateUpdate",
			API: "/api/w8t/noticeTemplate/noticeTemplateUpdate",
		},
		{
			Key: "NoticeTemplateDelete",
			API: "/api/w8t/noticeTemplate/noticeTemplateDelete",
		},
		{
			Key: "RuleGroupCreate",
			API: "/api/w8t/ruleGroup/ruleGroupCreate",
		},
		{
			Key: "RuleGroupUpdate",
			API: "/api/w8t/ruleGroup/ruleGroupUpdate",
		},
		{
			Key: "RuleGroupDelete",
			API: "/api/w8t/ruleGroup/ruleGroupDelete",
		},
		{
			Key: "RuleGroupList",
			API: "/api/w8t/ruleGroup/ruleGroupList",
		},
		{
			Key: "RuleTmplGroupCreate",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupCreate",
		},
		{
			Key: "RuleTmplGroupDelete",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupDelete",
		},
		{
			Key: "RuleTmplGroupList",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupList",
		},
		{
			Key: "RuleTmplCreate",
			API: "/api/w8t/ruleTmpl/ruleTmplCreate",
		},
		{
			Key: "RuleTmplDelete",
			API: "/api/w8t/ruleTmpl/ruleTmplDelete",
		},
		{
			Key: "RuleTmplList",
			API: "/api/w8t/ruleTmpl/ruleTmplList",
		},
		{
			Key: "CreateTenant",
			API: "/api/w8t/tenant/createTenant",
		},
		{
			Key: "UpdateTenant",
			API: "/api/w8t/tenant/updateTenant",
		},
		{
			Key: "DeleteTenant",
			API: "/api/w8t/tenant/deleteTenant",
		},
		{
			Key: "GetTenantList",
			API: "/api/w8t/tenant/getTenantList",
		},
		{
			Key: "CreateDashboard",
			API: "/api/w8t/dashboard/createDashboard",
		},
		{
			Key: "UpdateDashboard",
			API: "/api/w8t/dashboard/updateDashboard",
		},
		{
			Key: "DeleteDashboard",
			API: "/api/w8t/dashboard/deleteDashboard",
		},
		{
			Key: "ListDashboard",
			API: "/api/w8t/dashboard/listDashboard",
		},
		{
			Key: "GetDashboard",
			API: "/api/w8t/dashboard/getDashboard",
		},
		{
			Key: "SearchDashboard",
			API: "/api/w8t/dashboard/searchDashboard",
		},
	}

	perms = permissions

	globals.DBCli.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.UserPermissions{})
	repo.DBCli.Create(&models.UserPermissions{}, &permissions)

}

func InitUserRolesSQL() {

	permsString, _ := json.Marshal(perms)

	roles := models.UserRole{
		ID:          "ur-" + cmd.RandId(),
		Name:        "admin",
		Description: "system",
		Permissions: string(permsString),
		CreateAt:    time.Now().Unix(),
	}

	var adminRole models.UserRole
	globals.DBCli.Model(&models.UserRole{}).Where("name = ?", "admin").First(&adminRole)

	if adminRole.Name != "" {
		return
	}
	repo.DBCli.Create(&models.UserRole{}, &roles)

}
