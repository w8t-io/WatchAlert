package eval

import (
	"context"
	"fmt"
	"sync"
	"time"
	"watchAlert/alert/query"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	models "watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type AlertRuleWork struct {
	sync.RWMutex
	query.RuleQuery
	ctx        *ctx.Context
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

	rePushRule(arw.ctx, arw.rule)
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

func rePushRule(ctx *ctx.Context, alertRule chan *models.AlertRule) {

	var (
		ruleList []models.AlertRule
		// 创建一个通道用于接收处理结果
		resultCh = make(chan error)
		// 使用 WaitGroup 来等待所有规则的处理完成
		wg sync.WaitGroup
	)
	ctx.DB.DB().Where("enabled = ?", "1").Find(&ruleList)

	// 并发处理规则
	for _, rule := range ruleList {
		wg.Add(1)
		go func(rule models.AlertRule) {
			defer wg.Done()

			alertRule <- &rule

			resultCh <- nil
		}(rule)
	}

	// 等待所有规则的处理完成
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// 处理结果
	for result := range resultCh {
		if result != nil {
			fmt.Println("Error:", result)
		}
	}

}
