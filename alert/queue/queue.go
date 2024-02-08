package queue

import (
	"watchAlert/models"
)

var (
	AlertRuleChannel     = make(chan *models.AlertRule)
	QuitAlertRuleChannel = make(chan *string)
)
