package models

type AlertHisEvent struct {
	DatasourceId     string            `json:"datasource_id" gorm:"datasource_id"`
	DatasourceType   string            `json:"datasource_type"`
	Fingerprint      string            `json:"fingerprint"`
	RuleId           string            `json:"rule_id"`
	RuleName         string            `json:"rule_name"`
	Severity         int64             `json:"severity"`
	Metric           string            `json:"-" gorm:"metric"`
	MetricMap        map[string]string `json:"metric" gorm:"-"`
	EvalInterval     int64             `json:"eval_interval"`
	Annotations      string            `json:"annotations"`
	IsRecovered      bool              `json:"is_recovered" gorm:"-"`
	FirstTriggerTime int64             `json:"first_trigger_time"` // 第一次触发时间
	LastEvalTime     int64             `json:"last_eval_time"`     // 最近评估时间
	LastSendTime     int64             `json:"last_send_time"`     // 最近发送时间
	RecoverTime      int64             `json:"recover_time"`       // 恢复时间
}
