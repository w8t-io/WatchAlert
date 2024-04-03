package models

import (
	"time"
	"watchAlert/globals"
	"watchAlert/utils/cmd"
)

const SilenceCachePrefix = "mute-"

type AlertSilences struct {
	TenantId       string `json:"tenantId"`
	Id             string `json:"id"`
	Fingerprint    string `json:"fingerprint"`
	Datasource     string `json:"datasource"`
	DatasourceType string `json:"datasource_type"`
	StartsAt       int64  `json:"starts_at"`
	EndsAt         int64  `json:"ends_at"`
	CreateBy       string `json:"create_by"`
	UpdateBy       string `json:"update_by"`
	CreateAt       int64  `json:"create_at"`
	UpdateAt       int64  `json:"update_at"`
	Comment        string `json:"comment"`
}

func (as *AlertSilences) SetCache(expiration time.Duration) {

	globals.RedisCli.Set(as.TenantId+":"+SilenceCachePrefix+as.Fingerprint, cmd.JsonMarshal(as), expiration)

}

func (as *AlertSilences) GetCache(fingerprint string) (string, bool) {

	event, err := globals.RedisCli.Get(as.TenantId + ":" + SilenceCachePrefix + fingerprint).Result()
	if err != nil {
		return "", false
	}
	return event, true

}
