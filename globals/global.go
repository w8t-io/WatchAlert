package globals

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"go.uber.org/zap"
	"prometheus-manager/models"
)

var (
	Config    models.App
	Logger    *zap.Logger
	FeiShuCli *lark.Client
)
