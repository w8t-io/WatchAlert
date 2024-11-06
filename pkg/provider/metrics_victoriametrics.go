package provider

import (
	"fmt"
	"io"
	"strconv"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	utilsHttp "watchAlert/pkg/tools"
)

type VictoriaMetricsProvider struct {
	address string
}

func NewVictoriaMetricsClient(ds models.AlertDataSource) (MetricsFactoryProvider, error) {
	return VictoriaMetricsProvider{address: ds.HTTP.URL}, nil
}

type QueryResponse struct {
	VMData VMData `json:"data"`
}

type VMData struct {
	VMResult   []VMResult `json:"result"`
	ResultType string     `json:"resultType"`
}

type VMResult struct {
	Metric map[string]interface{} `json:"metric"`
	Value  []interface{}          `json:"value"`
}

func (v VictoriaMetricsProvider) Query(promQL string) ([]Metrics, error) {
	apiEndpoint := fmt.Sprintf("%s%s?query=%s&time=%d", v.address, "/api/v1/query", promQL, time.Now().Unix())

	resp, err := utilsHttp.Get(nil, apiEndpoint)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
		}
	}(resp.Body)

	var vmRespBody QueryResponse
	if err := utilsHttp.ParseReaderBody(resp.Body, &vmRespBody); err != nil {
		global.Logger.Sugar().Error(err.Error())
		return nil, err
	}

	return vmVectors(vmRespBody.VMData.VMResult), nil
}

func vmVectors(res []VMResult) []Metrics {
	var vectors []Metrics
	for _, item := range res {
		valueFloat, err := strconv.ParseFloat(item.Value[1].(string), 64)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return nil
		}
		vectors = append(vectors, Metrics{
			Metric:    item.Metric,
			Value:     valueFloat,
			Timestamp: item.Value[0].(float64),
		})
	}

	return vectors
}

func (v VictoriaMetricsProvider) Check() (bool, error) {
	res, err := utilsHttp.Get(nil, v.address+"/api/v1/labels")
	if err != nil {
		return false, err
	}

	if res.StatusCode != 200 {
		return false, err
	}
	return true, nil
}
