package services

import (
	"fmt"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	models "watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type ruleService struct {
	rule chan *models.AlertRule
	ctx  *ctx.Context
}

type InterRuleService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
}

func newInterRuleService(ctx *ctx.Context) InterRuleService {
	return &ruleService{
		rule: queue.AlertRuleChannel,
		ctx:  ctx,
	}
}

func (rs ruleService) Create(req interface{}) (interface{}, interface{}) {
	rule := req.(*models.AlertRule)
	ok := rs.ctx.DB.Rule().GetQuota(rule.TenantId)
	if !ok {
		return nil, fmt.Errorf("创建失败, 配额不足")
	}

	rs.rule <- rule

	err := rs.ctx.DB.Rule().Create(*rule)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (rs ruleService) Update(req interface{}) (interface{}, interface{}) {
	rule := req.(*models.AlertRule)
	alertInfo := models.AlertRule{}
	rs.ctx.DB.DB().Model(&models.AlertRule{}).
		Where("tenant_id = ? AND rule_id = ?", rule.TenantId, rule.RuleId).
		First(&alertInfo)

	/*
		重启协程
		判断当前状态是否是false 并且 历史状态是否为true
	*/
	if *alertInfo.Enabled == true && *rule.Enabled == false {
		if cancel, exists := queue.WatchCtxMap[rule.RuleId]; exists {
			cancel()
		}
	}
	if *alertInfo.Enabled == true && *rule.Enabled == true {
		if cancel, exists := queue.WatchCtxMap[rule.RuleId]; exists {
			cancel()
		}
	}

	// 删除缓存
	iter := rs.ctx.Redis.Redis().Scan(0, rule.TenantId+":"+models.FiringAlertCachePrefix+rule.RuleId+"*", 0).Iterator()
	keys := make([]string, 0)
	for iter.Next() {
		key := iter.Val()
		keys = append(keys, key)
	}
	rs.ctx.Redis.Redis().Del(keys...)

	// 启动协程
	if *rule.Enabled {
		rs.rule <- rule
		global.Logger.Sugar().Infof("重启 RuleId 为 %s 的 Worker 进程", rule.RuleId)
	}

	// 更新数据
	err := rs.ctx.DB.Rule().Update(*rule)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rs ruleService) Delete(req interface{}) (interface{}, interface{}) {
	rule := req.(*models.AlertRuleQuery)

	info, err := rs.ctx.DB.Rule().Search(*rule)
	if err != nil {
		return nil, err
	}

	err = rs.ctx.DB.Rule().Delete(*rule)
	if err != nil {
		return nil, err
	}

	// 退出该规则的协程
	if *info.Enabled {
		global.Logger.Sugar().Infof("停止 RuleId 为 %s 的 Worker 进程", rule.RuleId)
		if cancel, exists := queue.WatchCtxMap[rule.RuleId]; exists {
			cancel()
		}
	}

	iter := rs.ctx.Redis.Redis().Scan(0, rule.TenantId+":"+models.FiringAlertCachePrefix+rule.RuleId+"*", 0).Iterator()
	keys := make([]string, 0)
	for iter.Next() {
		key := iter.Val()
		keys = append(keys, key)
	}

	rs.ctx.Redis.Redis().Del(keys...)
	global.Logger.Sugar().Infof("删除队列数据 ->%s", keys)

	return nil, nil

}

func (rs ruleService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertRuleQuery)
	data, err := rs.ctx.DB.Rule().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rs ruleService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertRuleQuery)
	data, err := rs.ctx.DB.Rule().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}
