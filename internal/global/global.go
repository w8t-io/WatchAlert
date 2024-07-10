package global

import (
	"go.uber.org/zap"
	"watchAlert/config"
)

var (
	Layout  = "2006-01-02T15:04:05.000Z"
	Config  config.App
	Logger  *zap.Logger
	Version string
)
