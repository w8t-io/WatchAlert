package client

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/hash"
	utilsHttp "watchAlert/pkg/utils/http"
)

type VM struct {
	address string
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

type VMVector struct {
	Metric    map[string]interface{}
	Value     float64
	Timestamp float64
}

func NewVictoriaMetricsClient(ds models.AlertDataSource) VM {
	_, err := ds.CheckHealth()
	if err != nil {
		global.Logger.Sugar().Errorf(fmt.Sprintf("数据源不健康, Type: %s, Name: %s, Address: %s, Msg: %s", ds.Type, ds.Name, ds.HTTP.URL, err.Error()))
		return VM{}
	}

	return VM{address: ds.HTTP.URL}
}

func (a VM) Query(promQL string) ([]VMVector, error) {
	apiEndpoint := fmt.Sprintf("%s%s?query=%s&time=%d", a.address, "/prometheus/api/v1/query", promQL, time.Now().Unix())

	resp, err := utilsHttp.Get(apiEndpoint)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return nil, err
	}

	var vmRespBody QueryResponse
	err = json.Unmarshal(body, &vmRespBody)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return nil, err
	}

	return vmVectors(vmRespBody.VMData.VMResult), nil
}

func vmVectors(res []VMResult) []VMVector {
	var vectors []VMVector
	for _, item := range res {
		valueFloat, err := strconv.ParseFloat(item.Value[1].(string), 64)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return nil
		}
		vectors = append(vectors, VMVector{
			Metric:    item.Metric,
			Value:     valueFloat,
			Timestamp: item.Value[0].(float64),
		})
	}

	return vectors
}

func (a VMVector) GetFingerprint() string {
	if len(a.Metric) == 0 {
		return strconv.FormatUint(hash.HashNew(), 10)
	}

	var result uint64
	for labelName, labelValue := range a.Metric {
		sum := hash.HashNew()
		sum = hash.HashAdd(sum, labelName)
		sum = hash.HashAdd(sum, labelValue.(string))
		result ^= sum
	}

	return strconv.FormatUint(result, 10)
}

func (a VMVector) GetMetric() map[string]interface{} {
	a.Metric["value"] = a.Value
	return a.Metric
}
