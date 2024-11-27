package provider

import (
	"context"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"math"
	"time"
	"watchAlert/internal/models"
)

type PrometheusProvider struct {
	ExternalLabels map[string]interface{}
	apiV1          v1.API
}

func NewPrometheusClient(source models.AlertDataSource) (MetricsFactoryProvider, error) {
	client, err := api.NewClient(api.Config{
		Address: source.HTTP.URL,
	})
	if err != nil {
		return nil, err
	}

	apiV1 := v1.NewAPI(client)

	return PrometheusProvider{
		apiV1:          apiV1,
		ExternalLabels: source.Labels,
	}, nil
}

func (p PrometheusProvider) Query(promQL string) ([]Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _, err := p.apiV1.Query(ctx, promQL, time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	return ConvertVectors(result), nil
}

func ConvertVectors(value model.Value) (lst []Metrics) {
	items, ok := value.(model.Vector)
	if !ok {
		return
	}

	for _, item := range items {
		if math.IsNaN(float64(item.Value)) {
			continue
		}

		var metric = make(map[string]interface{})
		for k, v := range item.Metric {
			metric[string(k)] = string(v)
		}

		lst = append(lst, Metrics{
			Timestamp: float64(item.Timestamp),
			Value:     float64(item.Value),
			Metric:    metric,
		})
	}
	return
}

func (p PrometheusProvider) Check() (bool, error) {
	_, err := p.apiV1.Config(context.Background())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p PrometheusProvider) GetExternalLabels() map[string]interface{} {
	return p.ExternalLabels
}
