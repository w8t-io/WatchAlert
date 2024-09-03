package provider

import (
	"context"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
	utilsHttp "watchAlert/pkg/utils/http"
)

func CheckDatasourceHealth(datasource models.AlertDataSource) bool {
	var (
		err   error
		check bool
	)

	switch datasource.Type {
	case "Prometheus":
		prometheusClient, err := NewPrometheusClient(datasource)
		if err == nil {
			check, err = prometheusClient.Check()
		}
	case "VictoriaMetrics":
		vmClient, err := NewVictoriaMetricsClient(datasource)
		if err == nil {
			check, err = vmClient.Check()
		}
	case "Kubernetes":
		cli, err := client.NewKubernetesClient(context.Background(), datasource.KubeConfig)
		if err == nil {
			_, err = cli.GetWarningEvent("", 1)
			check = (err == nil)
		}
	case "ElasticSearch":
		res, err := utilsHttp.Get(nil, datasource.ElasticSearch.Url+"/_cat/health")
		check = (err == nil && res.StatusCode == 200)
	case "Jaeger", "Loki", "AliCloud", "CloudWatch":
		// 这几种数据源默认返回健康
		return true
	}

	// 检查数据源健康状况并返回结果
	if err != nil || !check {
		global.Logger.Sugar().Errorf("数据源不健康, Id: %s, Name: %s, Type: %s", datasource.Id, datasource.Name, datasource.Type)
		return false
	}

	return true
}
