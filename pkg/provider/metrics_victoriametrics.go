package provider

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"net/url"
	"strconv"
	"time"
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
	Status string `json:"status"`
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
	params := url.Values{}
	params.Add("query", promQL)
	params.Add("time", strconv.FormatInt(time.Now().Unix(), 10))
	fullURL := fmt.Sprintf("%s%s?%s", v.address, "/api/v1/query", params.Encode())
	resp, err := utilsHttp.Get(nil, fullURL, 10)
	if err != nil {
		logc.Error(context.Background(), err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var vmRespBody QueryResponse
	if err := utilsHttp.ParseReaderBody(resp.Body, &vmRespBody); err != nil {
		logc.Error(context.Background(), err.Error())
		return nil, err
	}

	return vmVectors(vmRespBody.VMData.VMResult), nil
}

func vmVectors(res []VMResult) []Metrics {
	var vectors []Metrics
	for _, item := range res {
		valueFloat, err := strconv.ParseFloat(item.Value[1].(string), 64)
		if err != nil {
			logc.Error(context.Background(), err.Error())
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
	res, err := utilsHttp.Get(nil, v.address+"/api/v1/labels", 10)
	if err != nil {
		return false, err
	}

	if res.StatusCode != 200 {
		return false, err
	}
	return true, nil
}
