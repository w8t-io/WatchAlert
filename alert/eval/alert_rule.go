package eval

import (
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
	rule       chan *models.AlertRule
	alertEvent models.AlertCurEvent
	quit       <-chan *string
	prom       query.Prometheus
}

type InterAlertRuleWork interface {
	Run()
}

func NewInterAlertRuleWork() InterAlertRuleWork {

	return &AlertRuleWork{
		rule: queue.AlertRuleChannel,
		quit: queue.QuitAlertRuleChannel,
	}

}

// Run 持续获取告警规则的状态
func (arw *AlertRuleWork) Run() {

	go func() {
		for {
			select {
			case rule := <-arw.rule:
				go arw.watch(*rule)
			}
		}
	}()

	// 重启服务后将历史 Rule 重新推到队列中
	arw.RePushRule()

}

func (arw *AlertRuleWork) watch(rule models.AlertRule) {

	ei := time.Second * time.Duration(rule.EvalInterval)
	timer := time.NewTimer(ei)

	for {
		select {
		case <-timer.C:
			switch rule.GetRuleType() {
			case "Prometheus":
				globals.Logger.Sugar().Infof("规则评估 -> %v", rule)
				arw.prom.Query(rule)
			}

		case <-arw.quit:
			timer.Stop()
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
