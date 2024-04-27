package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
)

type (
	SilenceCache struct {
		rc *redis.Client
	}

	InterSilenceCache interface {
		SetCache(r models.AlertSilences, expiration time.Duration)
		DelCache(r models.AlertSilenceQuery) error
		GetCache(r models.AlertSilenceQuery) (string, bool)
	}
)

func newSilenceCacheInterface(r *redis.Client) InterSilenceCache {
	return &SilenceCache{
		r,
	}
}

func (sc SilenceCache) SetCache(r models.AlertSilences, expiration time.Duration) {
	sc.rc.Set(r.TenantId+":"+models.SilenceCachePrefix+r.Fingerprint, cmd.JsonMarshal(r), expiration)
}

func (sc SilenceCache) DelCache(r models.AlertSilenceQuery) error {
	key := fmt.Sprintf("%s:%s%s", r.TenantId, models.SilenceCachePrefix, r.Fingerprint)
	_, err := sc.rc.Del(key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (sc SilenceCache) GetCache(r models.AlertSilenceQuery) (string, bool) {
	event, err := sc.rc.Get(r.TenantId + ":" + models.SilenceCachePrefix + r.Fingerprint).Result()
	if err != nil {
		return "", false
	}

	return event, true
}
