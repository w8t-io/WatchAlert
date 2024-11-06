package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"watchAlert/config"
)

var (
	Layout  = "2006-01-02T15:04:05.000Z"
	Config  config.App
	Logger  *zap.Logger
	Version string
	// StSignKey 签发的秘钥
	StSignKey = []byte(viper.GetString("jwt.WatchAlert"))
)
