package alert

import (
	"watchAlert/alert/consumer"
	"watchAlert/alert/eval"
	"watchAlert/pkg/ctx"
)

func Initialize(ctx *ctx.Context) {

	consumer.NewInterEvalConsumeWork(ctx).Run()
	eval.NewInterAlertRuleWork(ctx).Run()

}
