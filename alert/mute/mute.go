package mute

import (
	"watchAlert/controllers/repo"
	"watchAlert/models"
	"watchAlert/public/globals"
)

func IsMuted(alert *models.AlertCurEvent) bool {

	// 判断静默
	var am models.AlertSilences
	_, ok := am.GetCache(alert.Fingerprint)
	if ok {
		return true
	} else {
		ttl, _ := globals.RedisCli.TTL(models.SilenceCachePrefix + alert.Fingerprint).Result()
		// 如果剩余生存时间小于0，表示键已过期
		if ttl < 0 {
			repo.DBCli.Delete(repo.Delete{
				Table: models.AlertSilences{},
				Where: []interface{}{"fingerprint = ?", alert.Fingerprint},
			})
		}
	}

	return false

}
