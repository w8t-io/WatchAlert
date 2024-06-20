package models

type RuleGroups struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Description string `json:"description"`
}

type RuleGroupQuery struct {
	TenantId    string `json:"tenantId" form:"tenantId"`
	ID          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Query       string `json:"query" form:"query"`
	Page
}

type RuleGroupResponse struct {
	List []RuleGroups `json:"list"`
	Page
}
