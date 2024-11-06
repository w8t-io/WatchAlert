package mute

import (
	"time"
	models "watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

func IsMuted(ctx *ctx.Context, alert *models.AlertCurEvent) bool {
	//if IsSilence(ctx, alert) {
	//	return true
	//}

	if InTheEffectiveTime(alert) {
		return true
	}

	if RecoverNotify(alert) {
		return true
	}

	return false
}

// IsSilence 判断是否创建静默规则
//func IsSilence(ctx *ctx.Context, alert *models.AlertCurEvent) bool {
//	var as models.AlertSilences
//	ctx.DB.DB().Model(models.AlertSilences{}).Where("fingerprint = ?", alert.Fingerprint).First(&as)
//
//	_, ok := ctx.Redis.Silence().GetCache(models.AlertSilenceQuery{
//		TenantId:    as.TenantId,
//		Fingerprint: as.Fingerprint,
//	})
//
//	if ok {
//		return true
//	} else {
//		ttl, _ := ctx.Redis.Redis().TTL(alert.TenantId + ":" + models.SilenceCachePrefix + alert.Fingerprint).Result()
//		// 如果剩余生存时间小于0，表示键已过期
//		if ttl < 0 {
//			// 过期后标记为1
//			ctx.DB.DB().Model(models.AlertSilences{}).
//				Where("fingerprint = ? and status = ?", alert.Fingerprint, 0).
//				Update("status", 1)
//		}
//	}
//
//	return false
//}

// InTheEffectiveTime 判断生效时间
func InTheEffectiveTime(alert *models.AlertCurEvent) bool {
	if len(alert.EffectiveTime.Week) <= 0 {
		return false
	}

	var (
		p           bool
		currentTime = time.Now()
	)

	cwd := tools.TimeTransformToWeek(currentTime)
	for _, wd := range alert.EffectiveTime.Week {
		if cwd != wd {
			continue
		}
		p = true
	}

	if !p {
		return true
	}

	cts := tools.TimeTransformToSeconds(currentTime)
	if cts < alert.EffectiveTime.StartTime || cts > alert.EffectiveTime.EndTime {
		return true
	}

	return false
}

// RecoverNotify 判断是否推送恢复通知
func RecoverNotify(alert *models.AlertCurEvent) bool {
	// 如果是恢复告警，并且 恢复通知 == 1，即关闭恢复通知
	if alert.IsRecovered && !*alert.RecoverNotify {
		return true
	}

	return false
}
