package mute

import (
	"time"
	"watchAlert/internal/global"
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

	return InTheEffectiveTime(alert)
}

// InTheEffectiveTime 判断生效时间
func InTheEffectiveTime(alert *models.AlertCurEvent) bool {
	if len(alert.EffectiveTime.Week) <= 0 {
		return false
	}

	var (
		p           bool
		currentTime = time.Now()
	)

	cwd := currentWeekday(currentTime)
	for _, wd := range alert.EffectiveTime.Week {
		if cwd != wd {
			continue
		}
		p = true
	}

	if !p {
		return true
	}

	cts := currentTimeSeconds(currentTime)
	if cts < alert.EffectiveTime.StartTime || cts > alert.EffectiveTime.EndTime {
		return true
	}

	return false
}

func currentWeekday(ct time.Time) string {
	// 获取当前时间
	currentDate := ct.Format("2006-01-02")

	// 解析日期字符串为时间对象
	date, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return ""
	}

	return date.Weekday().String()
}

func currentTimeSeconds(ct time.Time) int {
	cs := ct.Hour()*3600 + ct.Minute()*60
	return cs
}
