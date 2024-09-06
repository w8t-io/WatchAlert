package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/http"
)

type JaegerDsProvider struct {
	url string
}

func NewJaegerClient(datasource models.AlertDataSource) (TracesFactoryProvider, error) {
	_, err := http.Get(nil, datasource.HTTP.URL)
	if err != nil {
		return JaegerDsProvider{}, err
	}

	return JaegerDsProvider{url: datasource.HTTP.URL}, nil
}

type JaegerResult struct {
	Data []JaegerData `json:"data"`
}

type JaegerData struct {
	TraceId string `json:"traceID"`
}

func (j JaegerDsProvider) Query(options TraceQueryOptions) ([]Traces, error) {
	curTime := time.Now()

	if options.Limit == 0 {
		options.Limit = 100
	}

	if options.StartAt == 0 {
		duration, _ := time.ParseDuration(strconv.Itoa(1) + "h")
		options.StartAt = curTime.Add(-duration).UnixNano()
	}

	if options.EndAt == 0 {
		options.EndAt = curTime.UnixNano()
	}

	args := fmt.Sprintf("/api/traces?service=%s&start=%d&end=%d&limit=%d&tags=%s", options.Service, options.StartAt, options.EndAt, options.Limit, options.Tags)
	requestURL := j.url + args
	res, err := http.Get(nil, requestURL)
	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(res.Body)
	var jaegerResult JaegerResult
	_ = json.Unmarshal(body, &jaegerResult)

	var data []Traces
	for _, t := range jaegerResult.Data {
		data = append(data, Traces{
			Service: options.Service,
			TraceId: t.TraceId,
		})
	}

	return data, nil
}

func (j JaegerDsProvider) Check() (bool, error) {
	res, err := http.Get(nil, j.url)
	if err != nil {
		return false, err
	}

	if res.StatusCode != 200 {
		return false, err
	}
	return true, nil
}

type JaegerServiceData struct {
	Data []string `json:"data"`
}

func (j JaegerDsProvider) GetJaegerService() (JaegerServiceData, error) {
	url := j.url + "/api/services"
	res, err := http.Get(nil, url)
	if err != nil {
		return JaegerServiceData{}, err
	}

	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return JaegerServiceData{}, fmt.Errorf("后端服务请求异常, Status: %d, Msg: %s", res.StatusCode, string(b))
	}

	body, _ := io.ReadAll(res.Body)
	var resData JaegerServiceData
	err = json.Unmarshal(body, &resData)
	if err != nil {
		return JaegerServiceData{}, fmt.Errorf("json.Unmarshal failed, %s", err.Error())
	}

	return resData, nil
}
