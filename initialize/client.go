package initialize

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"prometheus-manager/globals"
	"prometheus-manager/pkg/cache"
)

func InitClient() {

	feiShuClient()
	cacheClient()

}

func feiShuClient() {

	globals.FeiShuCli = lark.NewClient(globals.Config.FeiShu.AppID, globals.Config.FeiShu.AppSecret, lark.WithEnableTokenCache(true))

}

func cacheClient() {

	globals.CacheCli = cache.NewMemoryCache()

}
