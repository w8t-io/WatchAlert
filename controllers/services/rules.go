package services

import (
	"fmt"
	"watchAlert/alert/queue"
	"watchAlert/controllers/repo"
	models "watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type RuleService struct {
	rule chan *models.AlertRule
	repo.RuleRepo
}

type InterRuleService interface {
	Create(rule models.AlertRule) error
	Update(rule models.AlertRule) error
	Delete(tid, id string) error
	List(tid, ruleGroupId string) ([]models.AlertRule, error)
	Search(tid, ruleId string) models.AlertRule
}

func NewInterRuleService() InterRuleService {
	return &RuleService{
		rule: queue.AlertRuleChannel,
	}
}

func (rs *RuleService) Create(rule models.AlertRule) error {

	ok := rs.RuleRepo.GetQuota(rule.TenantId)
	if !ok {
		return fmt.Errorf("创建失败, 配额不足")
	}

	rule.RuleId = "a-" + cmd.RandId()

	newRule := rule.ParserRuleToGorm()

	err := repo.DBCli.Create(&models.AlertRule{}, &newRule)
	if err != nil {
		return err
	}

	rs.rule <- newRule

	return nil

}

func (rs *RuleService) Update(rule models.AlertRule) error {

	newRule := rule.ParserRuleToGorm()

	/*
		重启协程
		判断当前状态是否是false 并且 历史状态是否为true
	*/
	alertInfo := models.AlertRule{}
	globals.DBCli.Model(&models.AlertRule{}).Where("tenant_id = ? AND rule_id = ?", rule.TenantId, rule.RuleId).Find(&alertInfo)

	if alertInfo.Enabled == "true" && newRule.EnabledBool == false {
		if cancel, exists := queue.WatchCtxMap[newRule.RuleId]; exists {
			cancel()
		}
	}
	if alertInfo.Enabled == "true" && newRule.EnabledBool == true {
		if cancel, exists := queue.WatchCtxMap[newRule.RuleId]; exists {
			cancel()
		}
	}

	// 删除缓存
	iter := globals.RedisCli.Scan(0, rule.TenantId+":"+models.FiringAlertCachePrefix+rule.RuleId+"*", 0).Iterator()
	keys := make([]string, 0)
	for iter.Next() {
		key := iter.Val()
		keys = append(keys, key)
	}
	globals.RedisCli.Del(keys...)

	// 启动协程
	if newRule.EnabledBool {
		rs.rule <- newRule
		globals.Logger.Sugar().Infof("重启 RuleId 为 %s 的 Watch 进程", newRule.RuleId)
	}

	// 更新数据
	db := globals.DBCli.Model(&models.AlertRule{})
	db.Where("tenant_id = ? AND rule_id = ?", rule.TenantId, rule.RuleId)
	err := db.Updates(newRule).Error
	if err != nil {
		return err
	}

	return nil
}

func (rs *RuleService) Delete(tid, id string) error {

	var alertRule models.AlertRule
	data := repo.Delete{
		Table: alertRule,
		Where: []interface{}{"tenant_id = ? and rule_id = ?", tid, id},
	}

	err := repo.DBCli.Delete(data)
	if err != nil {
		return err
	}

	// 退出该规则的协程
	if alertRule.EnabledBool {
		globals.Logger.Sugar().Infof("停止 RuleId 为 %s 的Watch 进程", id)
		if cancel, exists := queue.WatchCtxMap[id]; exists {
			cancel()
		}
		//rs.quit <- &id
	}

	iter := globals.RedisCli.Scan(0, tid+":"+models.FiringAlertCachePrefix+id+"*", 0).Iterator()
	keys := make([]string, 0)
	for iter.Next() {
		key := iter.Val()
		keys = append(keys, key)
	}

	globals.RedisCli.Del(keys...)
	globals.Logger.Sugar().Infof("删除队列数据 ->%s", keys)

	return nil

}

func (rs *RuleService) List(tid, ruleGroupId string) ([]models.AlertRule, error) {

	var alertRuleList []models.AlertRule

	db := globals.DBCli.Model(&models.AlertRule{})
	db.Where("tenant_id = ?", tid)
	db.Where("rule_group_id = ?", ruleGroupId)
	db.Find(&alertRuleList)

	for k, v := range alertRuleList {
		newRule := v.ParserRuleToJson()
		alertRuleList[k] = *newRule
	}

	return alertRuleList, nil

}

func (rs *RuleService) Search(tid, ruleId string) models.AlertRule {

	var alertRule models.AlertRule
	globals.DBCli.Where("tenant_id = ? and rule_id = ?", tid, ruleId).Find(&alertRule)

	if alertRule.RuleName == "" {
		return models.AlertRule{}
	}

	newRule := alertRule.ParserRuleToJson()

	return *newRule

}
