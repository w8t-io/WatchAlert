package services

import (
	"watchAlert/alert/queue"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	models "watchAlert/models"
	"watchAlert/utils/cmd"
)

type RuleService struct {
	rule chan *models.AlertRule
	quit chan *string
}

type InterRuleService interface {
	Create(rule models.AlertRule) error
	Update(rule models.AlertRule) error
	Delete(id string) error
	List() ([]models.AlertRule, error)
	Search(ruleId string) models.AlertRule
}

func NewInterRuleService() InterRuleService {
	return &RuleService{
		rule: queue.AlertRuleChannel,
		quit: queue.QuitAlertRuleChannel,
	}
}

func (rs *RuleService) Create(rule models.AlertRule) error {

	rule.RuleIdStr = models.RuleId("a-" + cmd.RandId())
	rule.RuleId = string(rule.RuleIdStr)

	newRule := rule.ParserRuleToGorm()

	err := repo.DBCli.Create(&models.AlertRule{}, &newRule)
	if err != nil {
		return err
	}

	rs.rule <- newRule

	return nil

}

func (rs *RuleService) Update(rule models.AlertRule) error {

	rule.RuleId = string(rule.RuleIdStr)

	newRule := rule.ParserRuleToGorm()

	data := repo.Updates{
		Table:   models.AlertRule{},
		Where:   []string{"rule_id = ?", newRule.RuleId},
		Updates: &newRule,
	}
	err := repo.DBCli.Updates(data)
	if err != nil {
		return err
	}

	// 更新后重启协程
	rs.quit <- &newRule.RuleId

	rs.rule <- newRule

	return nil

}

func (rs *RuleService) Delete(id string) error {

	var alertRule models.AlertRule
	data := repo.Delete{
		Table: alertRule,
		Where: []string{"rule_id = ?", id},
	}

	err := repo.DBCli.Delete(data)
	if err != nil {
		return err
	}

	// 退出该规则的协程
	rs.quit <- &id

	return nil

}

func (rs *RuleService) List() ([]models.AlertRule, error) {

	var alertRuleList []models.AlertRule

	globals.DBCli.Find(&alertRuleList)

	for k, v := range alertRuleList {
		newRule := v.ParserRuleToJson()
		alertRuleList[k] = *newRule
	}

	return alertRuleList, nil

}

func (rs *RuleService) Search(ruleId string) models.AlertRule {

	var alertRule models.AlertRule
	globals.DBCli.Where("rule_id", ruleId).Find(&alertRule)

	if alertRule.RuleName == "" {
		return models.AlertRule{}
	}

	newRule := alertRule.ParserRuleToJson()

	return *newRule

}
