package eval

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
	"watchAlert/alert/process"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

type AlertRule struct {
	ctx                   *ctx.Context
	watchCtxMap           map[string]context.CancelFunc
	alarmRecoverWaitStore queue.AlarmRecoverWaitStore
}

func NewAlertRuleEval(ctx *ctx.Context, alarmRecoverWaitStore queue.AlarmRecoverWaitStore) AlertRule {
	return AlertRule{
		ctx:                   ctx,
		watchCtxMap:           make(map[string]context.CancelFunc),
		alarmRecoverWaitStore: alarmRecoverWaitStore,
	}
}

func (t *AlertRule) Submit(rule models.AlertRule) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	c, cancel := context.WithCancel(context.Background())
	t.watchCtxMap[rule.RuleId] = cancel
	go t.Eval(c, rule)
}

func (t *AlertRule) Stop(ruleId string) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	if cancel, exists := t.watchCtxMap[ruleId]; exists {
		cancel()
		delete(t.watchCtxMap, ruleId)
	}
}

func (t *AlertRule) Eval(ctx context.Context, rule models.AlertRule) {
	timer := time.NewTicker(time.Second * time.Duration(rule.EvalInterval))
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// 在规则评估前检查是否仍然启用，避免不必要的操作
			if !t.isRuleEnabled(rule.RuleId) {
				return
			}

			var curFiringKeys, curPendingKeys []string
			for _, dsId := range rule.DatasourceIdList {
				instance, err := t.ctx.DB.Datasource().GetInstance(dsId)
				if err != nil {
					global.Logger.Sugar().Error(err.Error())
				}

				switch rule.DatasourceType {
				case "Prometheus", "VictoriaMetrics":
					curFiringKeys, curPendingKeys = metrics(t.ctx, dsId, instance.Type, rule)
				case "AliCloudSLS", "Loki", "ElasticSearch":
					curFiringKeys = logs(t.ctx, dsId, instance.Type, rule)
				case "Jaeger":
					curFiringKeys = traces(t.ctx, dsId, instance.Type, rule)
				case "CloudWatch":
					curFiringKeys = cloudWatch(t.ctx, dsId, rule)
				case "KubernetesEvent":
					curFiringKeys = kubernetesEvent(t.ctx, dsId, rule)
				}
			}
			global.Logger.Sugar().Infof("规则评估 -> %v", tools.JsonMarshal(rule))

			t.Recover(rule, curFiringKeys)
			t.GC(rule, curFiringKeys, curPendingKeys)
		case <-ctx.Done():
			global.Logger.Sugar().Infof("停止 RuleId: %v, RuleName: %s 的 Watch 协程", rule.RuleId, rule.RuleName)
			return
		}
		timer.Reset(time.Second * time.Duration(rule.EvalInterval))
	}
}

func (t *AlertRule) Recover(rule models.AlertRule, curKeys []string) {
	firingKeys, err := ctx.Redis.Rule().GetAlertFiringCacheKeys(models.AlertRuleQuery{
		TenantId:         rule.TenantId,
		RuleId:           rule.RuleId,
		DatasourceIdList: rule.DatasourceIdList,
	})
	if err != nil {
		return
	}
	// 获取已恢复告警的keys
	recoverKeys := tools.GetSliceDifference(firingKeys, curKeys)
	if recoverKeys == nil {
		return
	}

	curTime := time.Now().Unix()
	for _, key := range recoverKeys {
		event := ctx.Redis.Event().GetCache(key)
		if event.IsRecovered == true {
			return
		}

		if _, exists := t.alarmRecoverWaitStore.Get(key); !exists {
			// 如果没有，则记录当前时间
			t.alarmRecoverWaitStore.Set(key, curTime)
			continue
		}

		// 判断是否在等待时间范围内
		wTime, _ := t.alarmRecoverWaitStore.Get(key)
		rt := time.Unix(wTime, 0).Add(time.Minute * time.Duration(global.Config.Server.AlarmConfig.RecoverWait)).Unix()
		if rt > curTime {
			continue
		}

		event.IsRecovered = true
		event.RecoverTime = curTime
		event.LastSendTime = 0

		ctx.Redis.Event().SetCache("Firing", event, 0)

		// 触发恢复删除带恢复中的 key
		t.alarmRecoverWaitStore.Remove(key)
	}
}

func (t *AlertRule) GC(rule models.AlertRule, curFiringKeys, curPendingKeys []string) {
	go process.GcPendingCache(t.ctx, rule, curPendingKeys)
	go process.GcRecoverWaitCache(t.ctx, t.alarmRecoverWaitStore, rule, curFiringKeys)
}

func (t *AlertRule) RePushTask() {
	ruleList, err := t.getRuleList()
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	g := new(errgroup.Group)
	for _, rule := range ruleList {
		rule := rule
		g.Go(func() error {
			t.Submit(rule)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		global.Logger.Sugar().Error(err.Error())
	}
}

// isRuleEnabled 检查规则是否启用
func (t *AlertRule) isRuleEnabled(ruleId string) bool {
	// 直接检查数据库或缓存中的当前启用状态
	return *t.ctx.DB.Rule().GetRuleObject(ruleId).Enabled
}

func (t *AlertRule) getRuleList() ([]models.AlertRule, error) {
	var ruleList []models.AlertRule
	if err := t.ctx.DB.DB().Where("enabled = ?", "1").Find(&ruleList).Error; err != nil {
		return ruleList, fmt.Errorf("获取 Rule List 失败, err: %s", err.Error())
	}
	return ruleList, nil
}
