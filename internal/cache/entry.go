package cache

import (
	"github.com/go-redis/redis"
	"watchAlert/pkg/client"
)

type (
	entryCache struct {
		redis    *redis.Client
		provider *ProviderPoolStore
	}

	InterEntryCache interface {
		Redis() *redis.Client
		Silence() InterSilenceCache
		Rule() InterRuleCache
		Event() InterEventCache
		ProviderPools() *ProviderPoolStore
	}
)

func NewEntryCache() InterEntryCache {
	r := client.InitRedis()
	p := NewClientPoolStore()

	return &entryCache{
		redis:    r,
		provider: p,
	}
}

func (e entryCache) Redis() *redis.Client              { return e.redis }
func (e entryCache) Silence() InterSilenceCache        { return newSilenceCacheInterface(e.redis) }
func (e entryCache) Rule() InterRuleCache              { return newRuleCacheInterface(e.redis) }
func (e entryCache) Event() InterEventCache            { return newEventCacheInterface(e.redis) }
func (e entryCache) ProviderPools() *ProviderPoolStore { return e.provider }
