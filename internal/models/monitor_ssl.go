package models

type MonitorSSLRule struct {
	TenantId             string `json:"tenantId"`
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Domain               string `json:"domain"`
	Description          string `json:"description"`
	ExpectTime           int64  `json:"expectTime"`   // 预期剩余时间
	EvalInterval         int64  `json:"evalInterval"` // Second
	NoticeId             string `json:"noticeId"`
	RepeatNoticeInterval int64  `json:"repeatNoticeInterval"` // 重复通知间隔时间
	TimeRemaining        int64  `json:"timeRemaining"`
	ResponseTime         string `json:"responseTime"`
	Enabled              *bool  `json:"enabled"`
	RecoverNotify        *bool  `json:"recoverNotify"`
}

type MonitorSSLRuleQuery struct {
	TenantId string `json:"tenantId"`
	ID       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Domain   string `json:"domain" form:"domain"`
	Query    string `json:"query" form:"query"`
}

func (m MonitorSSLRule) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"DomainName": m.Domain,
	}
}
