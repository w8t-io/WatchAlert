package mute

import (
	models "watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

func IsMuted(ctx *ctx.Context, alert *models.AlertCurEvent) bool {
	// 判断静默
	var as models.AlertSilences
	ctx.DB.DB().Model(models.AlertSilences{}).Where("fingerprint = ?", alert.Fingerprint).First(&as)

	_, ok := ctx.Redis.Silence().GetCache(models.AlertSilenceQuery{
		TenantId:    as.TenantId,
		Fingerprint: as.Fingerprint,
	})
	if ok {
		return true
	} else {
		ttl, _ := ctx.Redis.Redis().TTL(models.SilenceCachePrefix + alert.Fingerprint).Result()
		// 如果剩余生存时间小于0，表示键已过期
		if ttl < 0 {
			ctx.DB.DB().Model(models.AlertSilences{}).
				Where("tenant_id = ? AND fingerprint = ?", alert.TenantId, alert.Fingerprint).
				Delete(models.AlertSilences{})
		}
	}

	return false
}
