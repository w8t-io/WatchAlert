package globals

import (
	"github.com/go-redis/redis"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"watchAlert/config"
)

var (
	Config    config.App
	Logger    *zap.Logger
	FeiShuCli *lark.Client
	DBCli     *gorm.DB
	RedisCli  *redis.Client
)

var (
	Layout = "2006-01-02T15:04:05.000Z"
)
