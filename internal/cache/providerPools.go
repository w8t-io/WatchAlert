package cache

import (
	"fmt"
	"sync"
)

// ProviderPoolStore 提供商客户端存储池
type ProviderPoolStore struct {
	clients map[string]interface{}
	mux     sync.RWMutex
}

// NewClientPoolStore 创建一个新的 ProviderPoolStore 实例
func NewClientPoolStore() *ProviderPoolStore {
	return &ProviderPoolStore{
		clients: make(map[string]interface{}),
	}
}

// SetClient 设置通用客户端
func (p *ProviderPoolStore) SetClient(key string, client interface{}) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.clients[key] = client
}

// GetClient 获取通用客户端
func (p *ProviderPoolStore) GetClient(key string) (interface{}, error) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if client, exists := p.clients[key]; exists {
		return client, nil
	}

	return nil, fmt.Errorf("获取客户端错误, 客户端在缓存中不存在, datasourceId: %s", key)
}

// RemoveClient 移除通用客户端
func (p *ProviderPoolStore) RemoveClient(key string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	delete(p.clients, key)
}
