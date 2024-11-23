package storage

import (
	"context"
	"errors"
	"sync"
	"time"
	"watchAlert/internal/models"
)

// AlertsMQErrNotFound is returned if a Store cannot find the Alert.
var (
	AlertsMQErrNotFound = errors.New("alert not found")
)

// AlertsCurEventCache provides lock-coordinated to an in-memory map of alerts, keyed by
// their fingerprint. Resolved alerts are removed from the map based on
// gcInterval. An optional callback can be set which receives a slice of all
// resolved alerts that have been removed.
type AlertsCurEventCache struct {
	sync.RWMutex
	Data map[string]models.AlertCurEvent
}

// NewCurAlertsEventMap returns a new Alerts struct.
func NewCurAlertsEventMap() *AlertsCurEventCache {
	a := &AlertsCurEventCache{
		Data: make(map[string]models.AlertCurEvent),
	}

	return a
}

// Run starts the GC loop. The interval must be greater than zero; if not, the function will panic.
func (a *AlertsCurEventCache) Run(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			a.gc()
		}
	}
}

// 清理已处理的告警信息
func (a *AlertsCurEventCache) gc() {
	a.Lock()
	defer a.Unlock()

	//var resolved []*models.AlertRule
	//for fp, alert := range a.c {
	//	if alert.Resolved() {
	//		delete(a.c, fp)
	//		resolved = append(resolved, alert)
	//	}
	//}
	//a.cb(resolved)
}

/*
获取告警指纹是否存在
*/
func (a *AlertsCurEventCache) Get(fingerprint string) (models.AlertCurEvent, error) {

	a.RLock()
	alert, prs := a.Data[fingerprint]
	a.RUnlock()
	if !prs {
		return models.AlertCurEvent{}, AlertsMQErrNotFound
	}

	return alert, nil
}

// Set unconditionally sets the alert in memory.
func (a *AlertsCurEventCache) Set(fingerprint string, alert models.AlertCurEvent) error {

	a.Lock()
	a.Data[fingerprint] = alert
	a.Unlock()

	return nil
}

// Delete removes the Alert with the matching fingerprint from the store.
func (a *AlertsCurEventCache) Delete(fingerprint string) error {
	a.Lock()
	delete(a.Data, fingerprint)
	a.Unlock()

	return nil
}

// List returns a slice of Alerts currently held in memory.

func (a *AlertsCurEventCache) List() map[string]models.AlertCurEvent {
	a.RLock()
	defer a.RUnlock()

	return a.Data
}
