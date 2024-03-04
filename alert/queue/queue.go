package queue

import (
	"context"
	"watchAlert/models"
)

var (
	// WatchCtxMap 用于存储每个协程的上下文
	WatchCtxMap = make(map[string]context.CancelFunc)

	// AlertRuleChannel 用于消费用户创建的 Rule
	AlertRuleChannel = make(chan *models.AlertRule)
)
