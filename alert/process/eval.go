package process

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

func EvalConditionMetric(ctx *ctx.Context, m string, Threshold float64, v float64, f func() models.AlertCurEvent) {
	switch m {
	case ">":
		if v > Threshold {
			event := f()
			ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				SaveEventCache(ctx, event)
			}
		}
	case ">=":
		if v >= Threshold {
			event := f()
			ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				SaveEventCache(ctx, event)
			}
		}
	case "<":
		if v < Threshold {
			event := f()
			ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				SaveEventCache(ctx, event)
			}
		}
	case "<=":
		if v <= Threshold {
			event := f()
			ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				SaveEventCache(ctx, event)
			}
		}
	case "=":
		if v == Threshold {
			event := f()
			ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				SaveEventCache(ctx, event)
			}
		}
	case "!=":
		if v != Threshold {
			event := f()
			ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				SaveEventCache(ctx, event)
			}
		}
	}
}
