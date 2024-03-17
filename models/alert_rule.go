package models

import (
	"encoding/json"
	"sort"
	"strconv"
	"watchAlert/utils/hash"
)

type Duration int64

type LabelsMap map[string]string

type NoticeGroup []map[string]string

type AlertRule struct {
	//gorm.Model
	RuleId               string    `json:"ruleId" gorm:"ruleId"`
	RuleGroupId          string    `json:"ruleGroupId"`
	DatasourceType       string    `json:"datasourceType"`
	DatasourceIdList     []string  `json:"datasourceId" gorm:"-"`
	DatasourceId         string    `json:"-" gorm:"datasourceId"`
	RuleName             string    `json:"ruleName"`
	EvalInterval         int64     `json:"evalInterval"`
	ForDuration          int64     `json:"forDuration"`
	RepeatNoticeInterval int64     `json:"repeatNoticeInterval"`
	Description          string    `json:"description"`
	Annotations          string    `json:"annotations"`
	Labels               LabelsMap `json:"labels" gorm:"labels;serializer:json"`
	Severity             string    `json:"severity"`

	// Prometheus
	PrometheusConfig PrometheusConfig `json:"prometheusConfig" gorm:"prometheusConfig;serializer:json"`

	// 阿里云SLS
	AliCloudSLSConfig AliCloudSLSConfig `json:"alicloudSLSConfig" gorm:"alicloudSLSConfig;serializer:json"`

	// Loki
	LokiConfig LokiConfig `json:"lokiConfig" gorm:"lokiConfig;serializer:json"`

	NoticeId        string      `json:"noticeId"`
	NoticeGroupList NoticeGroup `json:"noticeGroup" gorm:"-"`
	NoticeGroup     string      `json:"-" gorm:"noticeGroup"`
	EnabledBool     bool        `json:"enabled" gorm:"-"`
	Enabled         string      `json:"-" gorm:"enabled"`
}

type PrometheusConfig struct {
	PromQL string `json:"promQL"`
}

type AliCloudSLSConfig struct {
	Project       string        `json:"project"`
	Logstore      string        `json:"logstore"`
	LogQL         string        `json:"logQL"`    // 查询语句
	LogScope      int           `json:"logScope"` // 相对查询的日志范围（单位分钟）,1(min) 5(min)...
	EvalCondition EvalCondition `json:"evalCondition" gorm:"evalCondition;serializer:json"`
}

type LokiConfig struct {
	LogQL         string        `json:"logQL"`
	LogScope      int           `json:"logScope"`
	EvalCondition EvalCondition `json:"evalCondition" gorm:"evalCondition;serializer:json"`
}

// EvalCondition 日志评估条件
type EvalCondition struct {
	Type     string `json:"type"`
	Operator string `json:"operator"`
	Value    int    `json:"value"`
}

type Fingerprint uint64

var (
	// cache the signature of an empty label set.
	emptyLabelSignature = hash.HashNew()
)

const SeparatorByte byte = 255

// Fingerprint returns a unique hash for the alert. It is equivalent to
// the fingerprint of the alert's label set.
func (a *AlertRule) Fingerprint() Fingerprint {

	// 没有配置标签，则用随机生成
	if len(a.Labels) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	// 定义map存储所有标签
	labelNames := make([]string, 0, len(a.Labels))
	for labelName := range a.Labels {
		labelNames = append(labelNames, labelName)
	}
	// 标签排序。用于根据标签做hash
	sort.Strings(labelNames)

	// 在随机生成的hash的基础上，新增标签hash
	sum := hash.HashNew()
	for _, labelName := range labelNames {
		sum = hash.HashAdd(sum, labelName)
		sum = hash.HashAddByte(sum, SeparatorByte)
		sum = hash.HashAdd(sum, a.Labels[labelName])
		sum = hash.HashAddByte(sum, SeparatorByte)
	}
	return Fingerprint(sum)

}

func (a *AlertRule) GetRuleType() string {

	return a.DatasourceType

}

func (a *AlertRule) ParserRuleToJson() *AlertRule {

	var (
		noticeGroupList  NoticeGroup
		datasourceIdList []string
	)

	_ = json.Unmarshal([]byte(a.NoticeGroup), &noticeGroupList)
	a.NoticeGroupList = noticeGroupList

	a.EnabledBool, _ = strconv.ParseBool(a.Enabled)

	_ = json.Unmarshal([]byte(a.DatasourceId), &datasourceIdList)
	a.DatasourceIdList = datasourceIdList

	return a

}

func (a *AlertRule) ParserRuleToGorm() *AlertRule {

	noticeGroup, _ := json.Marshal(a.NoticeGroupList)
	a.NoticeGroup = string(noticeGroup)

	a.Enabled = strconv.FormatBool(a.EnabledBool)

	datasourceIdListStr, _ := json.Marshal(a.DatasourceIdList)
	a.DatasourceId = string(datasourceIdListStr)

	return a

}
