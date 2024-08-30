package process

import (
	"context"
	"fmt"
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
	utilsHttp "watchAlert/pkg/utils/http"
)

func CheckDatasourceHealth(datasource models.AlertDataSource) (bool, error) {
	var (
		url      = datasource.HTTP.URL
		fullPath string
	)
	switch datasource.Type {
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
		cli, err := client.NewKubernetesClient(context.Background(), datasource.KubeConfig)
		if err != nil {
			return false, err
		}

		_, err = cli.GetWarningEvent("", 1)
		if err != nil {
			return false, err
		}

		return true, nil
	case "ElasticSearch":
		url = datasource.ElasticSearch.Url
		fullPath = "/_cat/health"
	}

	res, err := utilsHttp.Get(nil, url+fullPath)
	if err != nil {
		return false, fmt.Errorf("request url: %s failed", url)
	}

	if res.StatusCode != 200 {
		return false, fmt.Errorf("request url: %s failed , StatusCode: %d", url, res.StatusCode)
	}

	return true, nil
}
