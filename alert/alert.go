package alert

import (
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/alert/consumer"
	"watchAlert/alert/eval"
	"watchAlert/alert/probing"
	"watchAlert/alert/storage"
	"watchAlert/internal/global"
	"watchAlert/pkg/ctx"
)

var (
	AlertRule eval.AlertRuleEval

	ProductProbing probing.ProductProbing
	ConsumeProbing probing.ConsumeProbing
)

func Initialize(ctx *ctx.Context) {
	// 初始化告警规则消费任务
	consumer.NewInterEvalConsumeWork(ctx).Run()
	// 初始化监控告警的基础配置
	initAlarmConfig(ctx)
	alarmRecoverWaitStore := storage.NewAlarmRecoverStore(ctx)

	// 初始化告警规则评估任务
	AlertRule = eval.NewAlertRuleEval(ctx, alarmRecoverWaitStore)
	AlertRule.RePushTask()

	// 初始化拨测任务
	ConsumeProbing = probing.NewProbingConsumerTask(ctx)
	ProductProbing = probing.NewProbingTask(ctx)
	ProductProbing.RePushRule(&ConsumeProbing)
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
