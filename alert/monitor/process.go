package monitor

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

func SaveMonitorCacheEvent(ctx *ctx.Context, event models.AlertCurEvent) {
	firingKey := event.GetFiringAlertCacheKey()
	resFiring := ctx.Redis.Event().GetCache(firingKey)
	event.FirstTriggerTime = ctx.Redis.Event().GetFirstTime(firingKey)
	event.LastEvalTime = ctx.Redis.Event().GetLastEvalTime(firingKey)
	event.LastSendTime = resFiring.LastSendTime
	ctx.Redis.Event().SetCache("Firing", event, 0)
}
