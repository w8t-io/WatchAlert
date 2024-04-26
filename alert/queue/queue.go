package queue

import (
	"context"
	"watchAlert/internal/models"
)

var (
	// WatchCtxMap 用于存储每个协程的上下文
	WatchCtxMap = make(map[string]context.CancelFunc)

	// AlertRuleChannel 用于消费用户创建的 Rule
	AlertRuleChannel = make(chan *models.AlertRule)

	// RecoverWaitMap 存储等待被恢复的告警的 Key
	RecoverWaitMap = make(map[string]int64)
)
