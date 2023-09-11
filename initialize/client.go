package initialize

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"prometheus-manager/globals"
)

func InitClient() {

	feiShuClient()

}

func feiShuClient() {

	globals.FeiShuCli = lark.NewClient(globals.Config.FeiShu.AppID, globals.Config.FeiShu.AppSecret, lark.WithEnableTokenCache(true))

}