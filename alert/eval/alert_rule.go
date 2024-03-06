package eval

import (
	"context"
	"fmt"
	"sync"
	"time"
	"watchAlert/alert/query"
	"watchAlert/alert/queue"
	"watchAlert/globals"
	"watchAlert/models"
)

type AlertRuleWork struct {
	sync.RWMutex
	query.RuleQuery
	rule       chan *models.AlertRule
	alertEvent models.AlertCurEvent
}

type InterAlertRuleWork interface {
	Run()
}

func NewInterAlertRuleWork() InterAlertRuleWork {

	return &AlertRuleWork{
		rule: queue.AlertRuleChannel,
	}

}

// Run 持续获取告警规则的状态
func (arw *AlertRuleWork) Run() {

	go func() {
		for {
			select {
			case rule := <-arw.rule:
				// 创建一个用于停止协程的上下文
				ctx, cancel := context.WithCancel(context.Background())
				queue.WatchCtxMap[rule.RuleId] = cancel
				go arw.watch(*rule, ctx)
			}
		}
	}()

	// 重启服务后将历史 Rule 重新推到队列中
	arw.RePushRule()

}

func (arw *AlertRuleWork) watch(rule models.AlertRule, ctx context.Context) {

	ei := time.Second * time.Duration(rule.EvalInterval)
	timer := time.NewTimer(ei)

	for {
		select {
		case <-timer.C:
			globals.Logger.Sugar().Infof("规则评估 -> %v", rule)
			arw.Query(rule)

		case <-ctx.Done():
			globals.Logger.Sugar().Infof("停止 RuleId 为 %v 的 Watch 协程", rule.RuleId)
			return
		}

		timer.Reset(time.Second * time.Duration(rule.EvalInterval))

	}

}

func (arw *AlertRuleWork) RePushRule() {

	var (
		ruleList []models.AlertRule
		// 创建一个通道用于接收处理结果
		resultCh = make(chan error)
		// 使用 WaitGroup 来等待所有规则的处理完成
		wg sync.WaitGroup
	)
	globals.DBCli.Where("enabled = ?", "true").Find(&ruleList)

	// 并发处理规则
	for _, rule := range ruleList {
		wg.Add(1)
		go func(rule models.AlertRule) {
			defer wg.Done()

			arw.rule <- rule.ParserRuleToJson()

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
