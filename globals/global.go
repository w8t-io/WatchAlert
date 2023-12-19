package globals

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"prometheus-manager/controllers/dto"
	"prometheus-manager/utils/cache"
)

var (
	Config    dto.App
	Logger    *zap.Logger
	FeiShuCli *lark.Client
	CacheCli  *cache.InMemoryCache
	DBCli     *gorm.DB
)

var (
	Layout = "2006-01-02T15:04:05.000Z"
)
