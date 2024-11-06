package provider

import (
	"strconv"
	"watchAlert/pkg/tools"
)

const (
	PrometheusDsProvider      string = "Prometheus"
	VictoriaMetricsDsProvider string = "VictoriaMetrics"
)

type MetricsFactoryProvider interface {
	Query(promQL string) ([]Metrics, error)
	Check() (bool, error)
}

type Metrics struct {
	Metric    map[string]interface{}
	Value     float64
	Timestamp float64
}

func (m Metrics) GetFingerprint() string {
	if len(m.Metric) == 0 {
		return strconv.FormatUint(tools.HashNew(), 10)
	}

	var result uint64
	for labelName, labelValue := range m.Metric {
		sum := tools.HashNew()
		sum = tools.HashAdd(sum, labelName)
		sum = tools.HashAdd(sum, labelValue.(string))
		result ^= sum
	}

	return strconv.FormatUint(result, 10)
}

func (m Metrics) GetMetric() map[string]interface{} {
	m.Metric["value"] = m.Value
	return m.Metric
}
