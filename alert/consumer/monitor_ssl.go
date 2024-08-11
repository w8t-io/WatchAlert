package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"
	"watchAlert/alert/process"
	"watchAlert/alert/sender"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type MonitorSslConsumer struct {
	l            sync.RWMutex
	consumerPool map[string]context.CancelFunc
	ctx          *ctx.Context
}

func NewMonitorSslConsumer(ctx *ctx.Context) MonitorSslConsumer {
	return MonitorSslConsumer{
		ctx:          ctx,
		consumerPool: make(map[string]context.CancelFunc),
	}
}

func (m *MonitorSslConsumer) Add(r models.MonitorSSLRule) {
	m.l.Lock()
	m.l.Unlock()

	c, cancel := context.WithCancel(context.Background())
	m.consumerPool[r.ID] = cancel

	ticker := time.Tick(time.Second)
	go func(ctx context.Context, r models.MonitorSSLRule) {
		for {
			select {
			case <-ticker:
				key := fmt.Sprintf("%s:%s%s--", r.TenantId, models.FiringAlertCachePrefix, r.ID)
				result := m.ctx.Redis.Event().GetCache(key)
				handleAlert(m.ctx, result)
			case <-ctx.Done():
				return
			}
		}
	}(c, r)
}

func (m *MonitorSslConsumer) Stop(id string) {
	m.l.Lock()
	m.l.Unlock()

	if cancel, exists := m.consumerPool[id]; exists {
		cancel()
	}
}

func filterEvent(ctx *ctx.Context, alert models.AlertCurEvent) bool {
	var pass bool
	if !alert.IsRecovered {
		if alert.LastSendTime == 0 || alert.LastEvalTime >= alert.LastSendTime+alert.RepeatNoticeInterval*60 {
			alert.LastSendTime = time.Now().Unix()
			ctx.Redis.Event().SetCache("Firing", alert, 0)
			return true
		}
	} else {
		removeAlertFromCache(ctx, alert)
		err := process.RecordAlertHisEvent(ctx, alert)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
		}
		return true
	}
	return pass
}

// 删除缓存
func removeAlertFromCache(ctx *ctx.Context, alert models.AlertCurEvent) {
	key := fmt.Sprintf("%s:%s%s--", alert.TenantId, models.FiringAlertCachePrefix, alert.RuleId)
	ctx.Redis.Event().DelCache(key)
}

// 推送告警
func handleAlert(ctx *ctx.Context, alert models.AlertCurEvent) {
	if alert.RuleId == "" {
		return
	}

	if filterEvent(ctx, alert) {
		r := models.NoticeQuery{
			TenantId: alert.TenantId,
			Uuid:     alert.NoticeId,
		}
		noticeData, _ := ctx.DB.Notice().Get(r)
		alert.DutyUser = process.GetDutyUser(ctx, noticeData)
		err := sender.Sender(ctx, alert, noticeData)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}
	}
}
