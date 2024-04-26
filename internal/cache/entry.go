package cache

import (
	"github.com/go-redis/redis"
	"watchAlert/pkg/client"
)

type (
	entryCache struct {
		redis *redis.Client
	}

	InterEntryCache interface {
		Redis() *redis.Client
		Silence() InterSilenceCache
		Rule() InterRuleCache
		Event() InterEventCache
	}
)

func NewEntryCache() InterEntryCache {
	r := client.InitRedis()
	return &entryCache{
		redis: r,
	}
}

func (e entryCache) Redis() *redis.Client       { return e.redis }
func (e entryCache) Silence() InterSilenceCache { return newSilenceCacheInterface(e.redis) }
func (e entryCache) Rule() InterRuleCache       { return newRuleCacheInterface(e.redis) }
func (e entryCache) Event() InterEventCache     { return newEventCacheInterface(e.redis) }
