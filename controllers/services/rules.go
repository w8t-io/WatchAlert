package services

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"watchAlert/controllers/dto"
	"watchAlert/globals"
)

type RuleService struct{}

type InterRuleService interface {
	SelectPromRules(ruleGroup string) SearchRules
	CreatePromRule(ruleGroup string, ruleBody io.ReadCloser) error
	DeletePromRule(ruleGroup, ruleName string) error
	UpdatePromRule(ruleGroup, ruleName string, ruleBody io.ReadCloser) (dto.AlertRules, error)
	GetPromRuleData(ruleGroup, ruleName string) (dto.Rules, error)
}

func NewInterRuleService() InterRuleService {
	return &RuleService{}
}

const RULEHEADER = `groups:
- name: {{ name }}
  rules:
`

const RULECONTENT = `  - alert: {{ name }}
    expr: {{ expr }}
    for: {{ for }}
    labels:
      {{ key: value }}
    annotations:
      summary: ""
      description: "{{ description }}"

`

type SearchRules struct {
	ResRule []dto.Rules
	Number  int
}

// SelectPromRules 查询告警规则
func (rs *RuleService) SelectPromRules(ruleGroup string) SearchRules {

	var (
		rules   dto.AlertRules
		resRule SearchRules
	)

	ruleConfigFile, _ := getRuleConfigInfo(ruleGroup)
	file, err := ioutil.ReadFile(ruleConfigFile)
	if err != nil {
		log.Println("文件读取失败 ->", err)
		return SearchRules{}
	}
	_ = yaml.Unmarshal(file, &rules)

	for _, v := range rules.Groups {
		if v.Name == ruleGroup {
			resRule.ResRule = v.Rules
		}
		resRule.Number = len(v.Rules)
	}

	return resRule

}

// CreatePromRule 创建告警规则
func (rs *RuleService) CreatePromRule(ruleGroup string, ruleBody io.ReadCloser) error {

	var (
		rules     dto.AlertRules
		test, aaa string
	)
	if len(ruleGroup) == 0 {
		return fmt.Errorf("RuleGroup 不能为空")
	}

	ruleByte, _ := ioutil.ReadAll(ruleBody)
	_ = json.Unmarshal(ruleByte, &rules)

	ruleFilePath, ruleConfig := getRuleConfigInfo(ruleGroup)

	if len(ruleConfig.Groups) == 0 {
		ruleHeader := RULEHEADER
		ruleHeader = strings.ReplaceAll(ruleHeader, "{{ name }}", rules.Groups[0].Name)

		err := os.WriteFile(ruleFilePath, []byte(ruleHeader), 0644)
		if err != nil {
			log.Println("内容写入文件失败 ->", err)
			return err
		}
	}

	ruleContent := RULECONTENT
	ruleContent = strings.ReplaceAll(ruleContent, "{{ name }}", rules.Groups[0].Rules[0].Alert)
	ruleContent = strings.ReplaceAll(ruleContent, "{{ expr }}", rules.Groups[0].Rules[0].Expr)
	ruleContent = strings.ReplaceAll(ruleContent, "{{ for }}", rules.Groups[0].Rules[0].For)
	for k, v := range rules.Groups[0].Rules[0].Labels {
		if len(test) != 0 {
			aaa = "      "
		}
		test += aaa + k + ": " + v + "\n"
	}
	ruleContent = strings.ReplaceAll(ruleContent, "{{ key: value }}", test)
	ruleContent = strings.ReplaceAll(ruleContent, "{{ description }}", rules.Groups[0].Rules[0].Annotations.Description)

	f, _ := os.OpenFile(ruleFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	_, err := f.WriteString(ruleContent)
	if err != nil {
		log.Println("内容写入文件失败 ->", err)
		return err
	}

	err = formatRuleConfig(ruleGroup)

	err = globals.DBCli.Exec("UPDATE rule_group_data SET rule_number = rule_number + 1 WHERE name = ?", ruleGroup).Error
	if err != nil {
		return fmt.Errorf("数据库更新失败 -> %s", err)
	}

	return nil

}

// DeletePromRule 删除告警规则
func (rs *RuleService) DeletePromRule(ruleGroup, ruleName string) error {

	var newAfter []dto.Rules

	ruleFilePath, ruleConfig := getRuleConfigInfo(ruleGroup)

	after := ruleConfig
	// 获取将被删除元素的下标
	for i, v := range after.Groups {
		for _, vv := range v.Rules {
			if vv.Alert != ruleName {
				newAfter = append(newAfter, vv)
			}
		}
		after.Groups[i].Rules = newAfter
	}

	if len(newAfter) == 0 || len(after.Groups) == 0 {
		err := os.Remove(ruleFilePath)
		if err != nil {
			log.Println("文件删除失败 ->", err)
			return err
		}
	} else {
		afterMar, _ := yaml.Marshal(after)
		err := os.WriteFile(ruleFilePath, afterMar, 0644)
		if err != nil {
			log.Println("内容写入文件失败 ->", err)
			return err
		}
	}

	err := globals.DBCli.Exec("UPDATE rule_group_data SET rule_number = rule_number - 1 WHERE name = ?", ruleGroup).Error
	if err != nil {
		return fmt.Errorf("数据库更新失败 -> %s", err)
	}

	return nil
}

// UpdatePromRule 更新告警规则
func (rs *RuleService) UpdatePromRule(ruleGroup, ruleName string, ruleBody io.ReadCloser) (dto.AlertRules, error) {

	var (
		newRuleConfig dto.AlertRules
	)

	bodyIO, _ := ioutil.ReadAll(ruleBody)
	err := json.Unmarshal(bodyIO, &newRuleConfig)

	ruleFilePath, ruleConfig := getRuleConfigInfo(ruleGroup)

	for k, v := range ruleConfig.Groups {
		for kk, vv := range v.Rules {
			if vv.Alert == ruleName {
				ruleConfig.Groups[k].Rules[kk] = newRuleConfig.Groups[0].Rules[0]
			}
		}
	}

	f, _ := yaml.Marshal(ruleConfig)
	err = os.WriteFile(ruleFilePath, f, 0644)
	if err != nil {
		log.Println("内容写入文件失败 ->", err)
		return dto.AlertRules{}, err
	}

	return newRuleConfig, nil

}

func (rs *RuleService) GetPromRuleData(ruleGroup, ruleName string) (dto.Rules, error) {

	var getRule dto.Rules

	_, allRules := getRuleConfigInfo(ruleGroup)

	for k, v := range allRules.Groups {
		for kk, vv := range v.Rules {
			if vv.Alert == ruleName {
				getRule = allRules.Groups[k].Rules[kk]
				break
			}
		}
	}

	return getRule, nil

}

func getRuleConfigInfo(ruleGroup string) (string, dto.AlertRules) {

	var ruleConfig dto.AlertRules

	ruleFilePath := globals.Config.Prometheus.RulePath + "/" + ruleGroup + ".yaml"

	_, err := os.Stat(ruleFilePath)
	if err != nil {
		return ruleFilePath, dto.AlertRules{}
	}
	fConfig, _ := ioutil.ReadFile(ruleFilePath)
	_ = yaml.Unmarshal(fConfig, &ruleConfig)

	return ruleFilePath, ruleConfig

}

func formatRuleConfig(ruleGroup string) error {

	ruleConfigPath, ruleConfig := getRuleConfigInfo(ruleGroup)
	ruleConfigYaml, _ := yaml.Marshal(ruleConfig)
	err := os.WriteFile(ruleConfigPath, ruleConfigYaml, 0644)
	if err != nil {
		return fmt.Errorf("内容写入文件失败 -> %s", err)
	}
	return nil

}
