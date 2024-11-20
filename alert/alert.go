package alert

import (
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/alert/consumer"
	"watchAlert/alert/eval"
	"watchAlert/alert/monitor"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	"watchAlert/pkg/ctx"
)

var (
	MonEvalTask     monitor.MonitorSSLEval
	MonConsumerTask consumer.MonitorSslConsumer
	AlertRule       eval.AlertRule
)

func Initialize(ctx *ctx.Context) {
	// 初始化告警规则消费任务
	consumer.NewInterEvalConsumeWork(ctx).Run()
	// 初始化监控告警的基础配置
	initAlarmConfig(ctx)
	alarmRecoverWaitStore := queue.NewAlarmRecoverStore(ctx)
	// 初始化证书监控的消费任务
	MonConsumerTask = consumer.NewMonitorSslConsumer(ctx)
	// 初始化证书监控任务
	MonEvalTask = monitor.NewMonitorSSLEval()
	MonEvalTask.RePushTask(ctx, &MonConsumerTask)
	// 初始化告警规则评估任务
	AlertRule = eval.NewAlertRuleEval(ctx, alarmRecoverWaitStore)
	AlertRule.RePushTask()
}

func initAlarmConfig(ctx *ctx.Context) {
	get, err := ctx.DB.Setting().Get()
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
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
