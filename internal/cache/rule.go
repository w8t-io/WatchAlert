package cache

import (
	"github.com/go-redis/redis"
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
)

type (
	RuleCache struct {
		rc *redis.Client
	}

	InterRuleCache interface {
		GetAlertFiringCacheKeys(s models.AlertRuleQuery) ([]string, error)
		GetAlertPendingCacheKeys(s models.AlertRuleQuery) ([]string, error)
	}
)

func newRuleCacheInterface(r *redis.Client) InterRuleCache {
	return &RuleCache{
		r,
	}
}

// GetAlertFiringCacheKeys 获取当前规则中所有的 Firing 告警数据
func (rc RuleCache) GetAlertFiringCacheKeys(s models.AlertRuleQuery) ([]string, error) {
	var keys []string
	for _, v := range s.DatasourceIdList {
		keyPrefix := s.TenantId + ":" + models.FiringAlertCachePrefix + alertCacheTailKeys(s.RuleId, v)
		k, _ := client.Redis.Keys(keyPrefix).Result()
		keys = append(keys, k...)
	}

	return keys, nil
}

// GetAlertPendingCacheKeys 获取当前规则中所有的 Pending 告警数据
func (rc RuleCache) GetAlertPendingCacheKeys(s models.AlertRuleQuery) ([]string, error) {
	var keys []string
	for _, v := range s.DatasourceIdList {
		keyPrefix := s.TenantId + ":" + models.PendingAlertCachePrefix + alertCacheTailKeys(s.RuleId, v)
		k, _ := client.Redis.Keys(keyPrefix).Result()
		keys = append(keys, k...)
	}

	return keys, nil
}

func alertCacheTailKeys(ruleId, dsId string) string {
	return ruleId + "-" + dsId + "-" + "*"
}
