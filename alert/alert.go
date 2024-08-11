package alert

import (
	"watchAlert/alert/consumer"
	"watchAlert/alert/eval"
	"watchAlert/alert/task"
	"watchAlert/internal/global"
	"watchAlert/pkg/ctx"
)

var (
	MonEvalTask     task.MonitorSSLEval
	MonConsumerTask consumer.MonitorSslConsumer
)

func Initialize(ctx *ctx.Context) {
	consumer.NewInterEvalConsumeWork(ctx).Run()
	eval.NewInterAlertRuleWork(ctx).Run()
	initAlarmConfig(ctx)
	MonConsumerTask = consumer.NewMonitorSslConsumer(ctx)
	MonEvalTask = task.NewMonitorSSLEval()
	MonEvalTask.RePushTask(ctx, &MonConsumerTask)
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
