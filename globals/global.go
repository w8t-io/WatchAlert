package globals

import (
	"go.uber.org/zap"
	"prometheus-manager/models"
)

var (
	Config models.App
	Logger *zap.Logger
)
