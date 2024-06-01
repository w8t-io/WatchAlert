package process

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
)

func CalIndicatorValue(ctx *ctx.Context, m string, Threshold float64, rule models.AlertRule, v client.Vector, datasourceId string, curFiringKeys, curPendingKeys *[]string, severity string) {
	switch m {
	case ">":
		if v.Value > Threshold {
			f(ctx, datasourceId, curFiringKeys, curPendingKeys, v, rule, severity)
		}
	case ">=":
		if v.Value >= Threshold {
			f(ctx, datasourceId, curFiringKeys, curPendingKeys, v, rule, severity)
		}
	case "<":
		if v.Value < Threshold {
			f(ctx, datasourceId, curFiringKeys, curPendingKeys, v, rule, severity)
		}
	case "<=":
		if v.Value <= Threshold {
			f(ctx, datasourceId, curFiringKeys, curPendingKeys, v, rule, severity)
		}
	case "=":
		if v.Value == Threshold {
			f(ctx, datasourceId, curFiringKeys, curPendingKeys, v, rule, severity)
		}
	case "!=":
		if v.Value != Threshold {
			f(ctx, datasourceId, curFiringKeys, curPendingKeys, v, rule, severity)
		}
	}
}

func f(ctx *ctx.Context, datasourceId string, curFiringKeys, curPendingKeys *[]string, v client.Vector, rule models.AlertRule, severity string) {
	event := ParserDefaultEvent(rule)
	event.DatasourceId = datasourceId
	event.Fingerprint = v.GetFingerprint()
	event.Metric = v.GetMetric()
	event.Metric["severity"] = severity
	event.Severity = severity
	event.Annotations = cmd.ParserVariables(rule.PrometheusConfig.Annotations, event.Metric)

	firingKey := event.GetFiringAlertCacheKey()
	pendingKey := event.GetPendingAlertCacheKey()

	*curFiringKeys = append(*curFiringKeys, firingKey)
	*curPendingKeys = append(*curPendingKeys, pendingKey)

	ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
	if ok {
		SaveEventCache(ctx, event)
	}
}
