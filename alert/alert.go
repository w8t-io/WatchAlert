package alert

import (
	"watchAlert/alert/consumer"
	"watchAlert/alert/eval"
	"watchAlert/internal/global"
	"watchAlert/pkg/ctx"
)

func Initialize(ctx *ctx.Context) {

	consumer.NewInterEvalConsumeWork(ctx).Run()
	eval.NewInterAlertRuleWork(ctx).Run()
	initAlarmConfig(ctx)

}

func initAlarmConfig(ctx *ctx.Context) {
	get, err := ctx.DB.Setting().Get()
	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		return
	}

	global.Config.Server.AlarmConfig = get.AlarmConfig

	if global.Config.Server.AlarmConfig.RecoverWait == 0 {
		global.Config.Server.AlarmConfig.RecoverWait = 1
	}

	if global.Config.Server.AlarmConfig.GroupInterval == 0 {
		global.Config.Server.AlarmConfig.GroupInterval = 120
	}

	if global.Config.Server.AlarmConfig.GroupWait == 0 {
		global.Config.Server.AlarmConfig.GroupWait = 10
	}
}
