package prom

import (
	"context"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"time"
	"watchAlert/controllers/services"
	"watchAlert/globals"
)

type API struct {
	api.Client
	PromV1API v1.API
}

type Warnings []string

func NewPromClient(dsId string) API {

	datasource := services.NewInterAlertDataSourceService().Get(dsId, "Prometheus")

	client, err := api.NewClient(api.Config{
		Address: datasource[0].HTTPJson.URL,
	})
	if err != nil {
		globals.Logger.Sugar().Errorf("Prometheus 初始化客户端失败: %s", err)
	}

	v1api := v1.NewAPI(client)

	return API{
		PromV1API: v1api,
	}

}

func (a API) Query(promQL string) ([]Vector, Warnings, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := a.PromV1API.Query(ctx, promQL, time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		globals.Logger.Sugar().Errorf("Prometheus 执行query失败: %s", err)
		return nil, Warnings(warnings), err
	}

	return ConvertVectors(result), Warnings(warnings), nil

}
