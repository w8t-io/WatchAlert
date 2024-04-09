package models

type Tenant struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	CreateAt         int64  `json:"createAt"`
	CreateBy         string `json:"createBy"`
	Manager          string `json:"manager"`
	Description      string `json:"description"`
	UserNumber       int64  `json:"userNumber"`
	RuleNumber       int64  `json:"ruleNumber"`
	DutyNumber       int64  `json:"dutyNumber"`
	NoticeNumber     int64  `json:"noticeNumber"`
	RemoveProtection *bool  `json:"removeProtection" gorm:"type:BOOL"`
}

type TenantQuery struct {
	ID   string `form:"id"`
	Name string `form:"name"`
}
