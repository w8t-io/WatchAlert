package queue

import (
	"strings"
	"watchAlert/pkg/ctx"
)

// AlarmRecoverWaitStore 存储等待被恢复的告警的 Key
type AlarmRecoverWaitStore struct {
	Ctx  *ctx.Context
	Data map[string]int64
}

func NewAlarmRecoverStore(ctx *ctx.Context) AlarmRecoverWaitStore {
	return AlarmRecoverWaitStore{
		Ctx:  ctx,
		Data: make(map[string]int64),
	}
}

func (a AlarmRecoverWaitStore) Set(key string, t int64) {
	a.Ctx.Mux.Lock()
	defer a.Ctx.Mux.Unlock()
	a.Data[key] = t
}

func (a AlarmRecoverWaitStore) Get(key string) (int64, bool) {
	t, ok := a.Data[key]
	return t, ok
}

func (a AlarmRecoverWaitStore) Remove(key string) {
	a.Ctx.Mux.Lock()
	defer a.Ctx.Mux.Unlock()
	delete(a.Data, key)
}

func (a AlarmRecoverWaitStore) Search(keyPrefix string) []string {
	var keys []string
	for k := range a.Data {
		// 只获取当前规则组的告警。
		if strings.HasPrefix(k, keyPrefix) {
			keys = append(keys, k)
		}
	}

	return keys
}
