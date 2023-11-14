package cache

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}

type InMemoryCache struct {
	m         sync.Mutex
	StoreSync sync.Map
	Store     map[string]interface{}
}

type CacheItem struct {
	Values interface{}
	Expire time.Time
}

func NewMemoryCache() *InMemoryCache {

	var cacheData InMemoryCache
	cacheData = InMemoryCache{
		m:         sync.Mutex{},
		StoreSync: sync.Map{},
		Store:     make(map[string]interface{}),
	}

	go func() {
		for {

			cacheData.StoreSync.Range(func(key, value any) bool {
				v, _ := cacheData.StoreSync.Load(key)
				if expireTime := v.(CacheItem).Expire; time.Now().Sub(expireTime) > 0 {
					cacheData.StoreSync.Delete(key)
				}

				return true
			})

		}
	}()

	return &cacheData

}

func (i *InMemoryCache) Get(key string) interface{} {

	v, _ := i.StoreSync.Load(key)

	return v

}

func (i *InMemoryCache) Set(key string, value interface{}) {

	i.StoreSync.Store(key, CacheItem{
		Values: value,
		Expire: time.Now().Add(6 * time.Hour),
	})

}
