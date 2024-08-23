package monitor

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
	"watchAlert/alert/consumer"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/http"
)

type MonitorSSLEval struct {
	l           sync.RWMutex
	WatchCtxMap map[string]context.CancelFunc
}

func NewMonitorSSLEval() MonitorSSLEval {
	return MonitorSSLEval{
		WatchCtxMap: make(map[string]context.CancelFunc),
	}
}

func (t *MonitorSSLEval) Submit(ctx *ctx.Context, rule models.MonitorSSLRule) {
	t.l.Lock()
	defer t.l.Unlock()

	c, cancel := context.WithCancel(context.Background())
	t.WatchCtxMap[rule.ID] = cancel
	go t.Eval(c, ctx, rule)
}

func (t *MonitorSSLEval) Stop(id string) {
	t.l.Lock()
	defer t.l.Unlock()

	if cancel, exists := t.WatchCtxMap[id]; exists {
		cancel()
		delete(t.WatchCtxMap, id)
	}
}

func (t *MonitorSSLEval) Eval(ctx context.Context, w8tCtx *ctx.Context, rule models.MonitorSSLRule) {
	timer := time.NewTicker(time.Hour * time.Duration(rule.EvalInterval))
	defer timer.Stop()
	t.worker(w8tCtx, rule)

	for {
		select {
		case <-timer.C:
			t.worker(w8tCtx, rule)
		case <-ctx.Done():
			return
		}
	}
}

func (t *MonitorSSLEval) worker(w8tCtx *ctx.Context, rule models.MonitorSSLRule) {
	// 记录开始时间
	startTime := time.Now()

	resp, err := http.Get("https://" + rule.Domain)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 证书为空, 跳过检测
	if resp.TLS == nil {
		t.Stop(rule.ID)
		return
	}

	// 获取证书信息
	certs := resp.TLS.PeerCertificates[0]
	certTime := certs.NotAfter.Unix()
	currentTime := time.Now().Unix()

	// 计算剩余有效期时间
	TimeRemaining := (certTime - currentTime) / 86400

	// 创建事件
	event := t.processDefaultEvent(rule)
	event.TimeRemaining = TimeRemaining
	event.ResponseTime = fmt.Sprintf("%dms", time.Since(startTime).Milliseconds())
	event.Metric = rule.GetMetrics()
	event.Annotations = fmt.Sprintf("域名: %s, SSL证书即将到期, 剩余: %d天", rule.Domain, TimeRemaining)

	// 更新规则信息
	rule.TimeRemaining = event.TimeRemaining
	rule.ResponseTime = event.ResponseTime
	if err := w8tCtx.DB.MonitorSSL().Update(rule); err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	// 根据剩余时间与期望时间判断是否触发告警或恢复
	if TimeRemaining <= rule.ExpectTime {
		SaveMonitorCacheEvent(w8tCtx, event)
	} else {
		t.processRecover(w8tCtx, event)
	}
}

func (t *MonitorSSLEval) processDefaultEvent(rule models.MonitorSSLRule) models.AlertCurEvent {
	return models.AlertCurEvent{
		TenantId:             rule.TenantId,
		RuleId:               rule.ID,
		RuleName:             rule.Name,
		EvalInterval:         rule.EvalInterval,
		NoticeId:             rule.NoticeId,
		IsRecovered:          false,
		RepeatNoticeInterval: rule.RepeatNoticeInterval,
		DutyUser:             "暂无", // 默认暂无值班人员, 渲染模版时会实际判断 Notice 是否存在值班人员
		RecoverNotify:        rule.RecoverNotify,
	}
}

func (t *MonitorSSLEval) processRecover(ctx *ctx.Context, event models.AlertCurEvent) {
	key := fmt.Sprintf("%s:%s%s--", event.TenantId, models.FiringAlertCachePrefix, event.RuleId)
	cache := ctx.Redis.Event().GetCache(key)

	// 提前返回，如果缓存中没有相关数据或已经恢复
	if cache.RuleId == "" || cache.IsRecovered {
		return
	}

	// 检查事件的剩余时间是否大于或等于缓存中的剩余时间
	if event.TimeRemaining >= cache.TimeRemaining {
		event.FirstTriggerTime = ctx.Redis.Event().GetFirstTime(key)
		event.IsRecovered = true
		event.RecoverTime = time.Now().Unix()
		event.LastSendTime = 0
		ctx.Redis.Event().SetCache("Firing", event, 0)
	}
}

func (t *MonitorSSLEval) RePushTask(ctx *ctx.Context, consumer *consumer.MonitorSslConsumer) {
	var ruleList []models.MonitorSSLRule
	if err := ctx.DB.DB().Where("enabled = ?", "1").Find(&ruleList).Error; err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	g := new(errgroup.Group)
	for _, rule := range ruleList {
		rule := rule
		g.Go(func() error {
			t.Submit(ctx, rule)
			consumer.Add(rule)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		global.Logger.Sugar().Error(err.Error())
	}
}
