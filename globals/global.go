package globals

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"go.uber.org/zap"
	"prometheus-manager/models"
	"prometheus-manager/services/cache"
)

var (
	Config    models.App
	Logger    *zap.Logger
	FeiShuCli *lark.Client
	CacheCli  *cache.InMemoryCache
)
