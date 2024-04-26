package models

import (
	"fmt"
	utilsHttp "watchAlert/pkg/utils/http"
)

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
	Enabled          *bool  `json:"enabled" `
}

type HTTP struct {
	URL     string `json:"url"`
	Timeout int64  `json:"timeout"`
}

type DatasourceQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	Id       string `json:"id" form:"id"`
	Type     string `json:"type" form:"type"`
	Query    string `json:"query" form:"query"`
}

func (ds AlertDataSource) CheckHealth() (bool, error) {
	var (
		url      = ds.HTTP.URL
		fullPath string
	)
	switch ds.Type {
	case "Prometheus":
		fullPath = "/-/healthy"
	case "Jaeger":
		fullPath = "/"
	case "Loki":
		fullPath = "/"
	case "AliCloud":
		return true, nil
	}

	res, err := utilsHttp.Get(url + fullPath)
	if err != nil {
		return false, fmt.Errorf("request url: %s failed", url)
	}

	if res.StatusCode != 200 {
		return false, fmt.Errorf("request url: %s failed , StatusCode: %d", url, res.StatusCode)
	}

	return true, nil
}
