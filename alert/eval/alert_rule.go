package eval

import (
	"context"
	"sync"
	"time"
	"watchAlert/alert/query"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	models "watchAlert/internal/models"
	"watchAlert/internal/services"
	"watchAlert/pkg/ctx"
)

type AlertRuleWork struct {
	sync.RWMutex
	query.RuleQuery
	ctx *ctx.Context
	services.InterAlertService
	rule       chan *models.AlertRule
	alertEvent models.AlertCurEvent
}

type InterAlertRuleWork interface {
	Run()
}

func NewInterAlertRuleWork(ctx *ctx.Context) InterAlertRuleWork {
	return &AlertRuleWork{
		rule: queue.AlertRuleChannel,
		ctx:  ctx,
	}
}

// Run 持续获取告警规则的状态
func (arw *AlertRuleWork) Run() {

	go func() {
		for {
			select {
			case rule := <-arw.rule:
				if *rule.Enabled {
					// 创建一个用于停止协程的上下文
					c, cancel := context.WithCancel(context.Background())
					queue.WatchCtxMap[rule.RuleId] = cancel
					go arw.worker(*rule, c)
				}
			}
		}
	}()

	// 重启服务后将历史 Rule 重新推到队列中
	services.AlertService.RePushRule(arw.ctx, arw.rule)

}

func (arw *AlertRuleWork) worker(rule models.AlertRule, ctx context.Context) {

	ei := time.Second * time.Duration(rule.EvalInterval)
	timer := time.NewTimer(ei)

	for {
		select {
		case <-timer.C:
			global.Logger.Sugar().Infof("规则评估 -> %v", rule)
			arw.Query(arw.ctx, rule)

		case <-ctx.Done():
			global.Logger.Sugar().Infof("停止 RuleId 为 %v 的 Watch 协程", rule.RuleId)
			return
		}

		timer.Reset(time.Second * time.Duration(rule.EvalInterval))

	}

}
