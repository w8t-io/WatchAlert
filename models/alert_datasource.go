package models

type AlertDataSource struct {
	TenantId         string `json:"tenantId"`
	Id               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	HTTP             HTTP   `json:"http" gorm:"http;serializer:json"`
	AliCloudEndpoint string `json:"alicloudEndpoint"`
	AliCloudAk       string `json:"alicloudAk"`
	AliCloudSk       string `json:"alicloudSk"`
	Description      string `json:"description"`
	EnabledBool      bool   `json:"enabled" gorm:"-"`
	Enabled          string `json:"-" gorm:"enabled"`
}

type HTTP struct {
	URL     string `json:"url"`
	Timeout int64  `json:"timeout"`
}
