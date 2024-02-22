package models

import (
	"encoding/json"
	"sort"
	"strconv"
	"watchAlert/utils/hash"
)

type Duration int64

type RuleConfig struct {
	PromQL   string `json:"promQL"`
	Severity int64  `json:"severity"`
}

type LabelsMap map[string]string

type RuleId string

type NoticeGroup []map[string]string

type AlertRule struct {
	//gorm.Model
	RuleIdStr            RuleId     `json:"ruleId" gorm:"-"`
	RuleId               string     `json:"-" gorm:"ruleId"`
	RuleGroupId          string     `json:"ruleGroupId"`
	DatasourceType       string     `json:"datasourceType"`
	DatasourceIdList     []string   `json:"datasourceId" gorm:"-"`
	DatasourceId         string     `json:"-" gorm:"datasourceId"`
	RuleName             string     `json:"ruleName"`
	EvalInterval         int64      `json:"evalInterval"`
	ForDuration          int64      `json:"forDuration"`
	RepeatNoticeInterval int64      `json:"repeatNoticeInterval"`
	Description          string     `json:"description"`
	Annotations          string     `json:"annotations"`
	RuleConfigJson       RuleConfig `json:"ruleConfig" gorm:"-"`
	RuleConfig           string     `json:"-" gorm:"ruleConfig"`
	LabelsMap            LabelsMap   `json:"labels" gorm:"-"`
	Labels               string      `json:"-" gorm:"labels"`
	NoticeId             string      `json:"noticeId"`
	NoticeGroupList      NoticeGroup `json:"noticeGroup" gorm:"-"`
	NoticeGroup          string      `json:"-" gorm:"noticeGroup"`
	EnabledBool          bool        `json:"enabled" gorm:"-"`
	Enabled              string      `json:"-" gorm:"enabled"`
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
	if len(a.LabelsMap) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	// 定义map存储所有标签
	labelNames := make([]string, 0, len(a.LabelsMap))
	for labelName := range a.LabelsMap {
		labelNames = append(labelNames, labelName)
	}
	// 标签排序。用于根据标签做hash
	sort.Strings(labelNames)

	// 在随机生成的hash的基础上，新增标签hash
	sum := hash.HashNew()
	for _, labelName := range labelNames {
		sum = hash.HashAdd(sum, labelName)
		sum = hash.HashAddByte(sum, SeparatorByte)
		sum = hash.HashAdd(sum, string(a.LabelsMap[labelName]))
		sum = hash.HashAddByte(sum, SeparatorByte)
	}
	return Fingerprint(sum)

}

func (a *AlertRule) GetRuleType() string {

	return a.DatasourceType

}

func (a *AlertRule) ParserRuleToJson() *AlertRule {

	var (
		ruleConfig       RuleConfig
		labelsMap        map[string]string
		noticeGroupList  NoticeGroup
		datasourceIdList []string
	)

	a.RuleIdStr = RuleId(a.RuleId)

	_ = json.Unmarshal([]byte(a.RuleConfig), &ruleConfig)
	a.RuleConfigJson = ruleConfig

	_ = json.Unmarshal([]byte(a.Labels), &labelsMap)
	a.LabelsMap = labelsMap

	_ = json.Unmarshal([]byte(a.NoticeGroup), &noticeGroupList)
	a.NoticeGroupList = noticeGroupList

	a.EnabledBool, _ = strconv.ParseBool(a.Enabled)

	_ = json.Unmarshal([]byte(a.DatasourceId), &datasourceIdList)
	a.DatasourceIdList = datasourceIdList

	return a

}

func (a *AlertRule) ParserRuleToGorm() *AlertRule {

	ruleConfigStr, _ := json.Marshal(a.RuleConfigJson)
	a.RuleConfig = string(ruleConfigStr)

	labelsStr, _ := json.Marshal(a.LabelsMap)
	a.Labels = string(labelsStr)

	noticeGroup, _ := json.Marshal(a.NoticeGroupList)
	a.NoticeGroup = string(noticeGroup)

	a.Enabled = strconv.FormatBool(a.EnabledBool)

	datasourceIdListStr, _ := json.Marshal(a.DatasourceIdList)
	a.DatasourceId = string(datasourceIdListStr)

	return a

}
