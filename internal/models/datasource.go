package models

import (
	"fmt"
	utilsHttp "watchAlert/pkg/utils/http"
)

type AlertDataSource struct {
	TenantId         string        `json:"tenantId"`
	Id               string        `json:"id"`
	Name             string        `json:"name"`
	Type             string        `json:"type"`
	HTTP             HTTP          `json:"http" gorm:"http;serializer:json"`
	AliCloudEndpoint string        `json:"alicloudEndpoint"`
	AliCloudAk       string        `json:"alicloudAk"`
	AliCloudSk       string        `json:"alicloudSk"`
	AWSCloudWatch    AWSCloudWatch `json:"awsCloudwatch" gorm:"awsCloudwatch;serializer:json"`
	Description      string        `json:"description"`
	KubeConfig       string        `json:"kubeConfig"`
	Enabled          *bool         `json:"enabled" `
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

type AWSCloudWatch struct {
	//Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

func (ds AlertDataSource) CheckHealth() (bool, error) {
	var (
		url      = ds.HTTP.URL
		fullPath string
	)
	switch ds.Type {
	case "Prometheus", "VictoriaMetrics":
		fullPath = "/-/healthy"
	case "Jaeger":
		return true, nil
	case "Loki":
		return true, nil
	case "AliCloud":
		return true, nil
	case "CloudWatch":
		return true, nil
	case "Kubernetes":
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

type PromQueryReq struct {
	DatasourceType string `json:"datasourceType"`
	Addr           string `form:"addr"`
	Query          string `form:"query"`
}

type PromQueryRes struct {
	Data data `json:"data"`
}

type data struct {
	Result     []result `json:"result"`
	ResultType string   `json:"resultType"`
}

type result struct {
	Metric map[string]interface{} `json:"metric"`
	Value  []interface{}          `json:"value"`
}
