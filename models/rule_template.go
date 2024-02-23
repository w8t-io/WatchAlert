package models

type RuleTemplateGroup struct {
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Description string `json:"description"`
}

type RuleTemplate struct {
	RuleGroupName  string     `json:"ruleGroupName"`
	RuleName       string     `json:"ruleName"`
	DatasourceType string     `json:"datasourceType"`
	RuleConfigJson RuleConfig `json:"ruleConfig" gorm:"-"`
	RuleConfig     string     `json:"-" gorm:"ruleConfig"`
	EvalInterval   int64      `json:"evalInterval"`
	ForDuration    int64      `json:"forDuration"`
	Annotations    string     `json:"annotations"`
}
