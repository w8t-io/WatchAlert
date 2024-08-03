package process

import (
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

// EvalCondition 评估告警条件
func EvalCondition(ctx *ctx.Context, f func() models.AlertCurEvent, value float64, ec models.EvalCondition) {

	switch ec.Type {
	case "count", "metric":
		switch ec.Operator {
		case ">":
			if value > ec.Value {
				SaveAlertEvent(ctx, f())
			}
		case ">=":
			if value >= ec.Value {
				SaveAlertEvent(ctx, f())
			}
		case "<":
			if value < ec.Value {
				SaveAlertEvent(ctx, f())
			}
		case "<=":
			if value <= ec.Value {
				SaveAlertEvent(ctx, f())
			}
		case "==":
			if value == ec.Value {
				SaveAlertEvent(ctx, f())
			}
		case "!=":
			if value != ec.Value {
				SaveAlertEvent(ctx, f())
			}
		default:
			global.Logger.Sugar().Error("无效的评估条件", ec.Type, ec.Operator, ec.Value)
		}
	default:
		global.Logger.Sugar().Error("无效的评估类型", ec.Type)
	}

}

func SaveAlertEvent(ctx *ctx.Context, event models.AlertCurEvent) {
	ok := ctx.DB.Rule().GetRuleIsExist(event.RuleId)
	if ok {
		SaveEventCache(ctx, event)
	}
}
