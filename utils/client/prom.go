package client

import (
	"context"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"math"
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
		Address: datasource[0].HTTP.URL,
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

type Vector struct {
	Key       string       `json:"key"`
	Labels    model.Metric `json:"labels"`
	Timestamp int64        `json:"timestamp"`
	Value     float64      `json:"value"`
}

func ConvertVectors(value model.Value) (lst []Vector) {
	switch value.Type() {
	case model.ValVector:
		items, ok := value.(model.Vector)
		if !ok {
			return
		}

		for _, item := range items {
			if math.IsNaN(float64(item.Value)) {
				continue
			}

			lst = append(lst, Vector{
				Key:       item.Metric.String(),
				Timestamp: item.Timestamp.Unix(),
				Value:     float64(item.Value),
				Labels:    item.Metric,
			})
		}
	case model.ValMatrix:
		items, ok := value.(model.Matrix)
		if !ok {
			return
		}

		for _, item := range items {
			if len(item.Values) == 0 {
				return
			}

			last := item.Values[len(item.Values)-1]

			if math.IsNaN(float64(last.Value)) {
				continue
			}

			lst = append(lst, Vector{
				Key:       item.Metric.String(),
				Labels:    item.Metric,
				Timestamp: last.Timestamp.Unix(),
				Value:     float64(last.Value),
			})
		}
	case model.ValScalar:
		item, ok := value.(*model.Scalar)
		if !ok {
			return
		}

		if math.IsNaN(float64(item.Value)) {
			return
		}

		lst = append(lst, Vector{
			Key:       "{}",
			Timestamp: item.Timestamp.Unix(),
			Value:     float64(item.Value),
			Labels:    model.Metric{},
		})
	default:
		return
	}

	return
}

func (v Vector) GetFingerprint() string {
	return v.Labels.FastFingerprint().String()
}

func (v Vector) GetMetric() map[string]interface{} {
	// handle series tags
	metricMap := make(map[string]interface{})
	for label, value := range v.Labels {
		metricMap[string(label)] = string(value)
	}
	metricMap["value"] = v.Value
	return metricMap
}
